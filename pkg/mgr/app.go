package mgr

import (
  "context"
  "fmt"
  "github.com/sirupsen/logrus"
  "os"
  "os/signal"
  "sync"
  "syscall"
  "time"
)

type mgr struct {
  Timeout           time.Duration // Таймаут ожидания 0 - без таймаута
  log               *logrus.Entry
  ctx               context.Context
  IsCli             bool          // Вывод в консоль
  ShutdownTimeoutMs time.Duration // Время ожидания закрытия приложения
  ctxCancel         context.CancelFunc
  tasksCompleted    bool
  failFn            func()         // Функция при ошибке в одной из задач
  wait              bool           // Флаг начала ожидания работы приложения
  isErr             bool           // Одна из задач завершилась с ошибкой
  started           bool           // Флаг запуска всех задач
  wg                sync.WaitGroup // Ожидание задач запуска
  mx                sync.Mutex
  timer             *time.Timer
  fnsStarted        []func()       // Выполнить после всех задач. Не влияет на успешность старта
  fnsAfterTasks     []func() error // Выполнить после всех задач на запуск
  fnsFailed         []func()       // Подписка на событие ошибки запуска
  sig               chan os.Signal
  chansExit         []chan struct{} // Каналы ожидания завершения работы
}

func New(ctx context.Context, cancel context.CancelFunc, l *logrus.Entry) *mgr {
  o := &mgr{
    ctx:               ctx,
    ctxCancel:         cancel,
    log:               l,
    ShutdownTimeoutMs: defShutdownTimeout,
  }

  go func() {
    o.sig = make(chan os.Signal, 1)
    signal.Notify(o.sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-o.sig
    cancel()
  }()

  return o
}

func (a *mgr) Ready() bool {
  return a.started
}

// Запускает таймер ожидания завершения запуска
func (a *mgr) startTimeout() {
  if a.Timeout > 0 && a.timer == nil {
    fn := func() {
      if !a.tasksCompleted {
        a.log.Warn("Application launch was stopped due to timeout")
        a.ctxCancel()
      }
    }
    a.timer = time.AfterFunc(a.Timeout, fn)
  }
}

// Task Добавляет задачу на запуск
func (a *mgr) Task(fn func() error) {
  a.mx.Lock()
  if a.wait {
    a.log.Panic(ErrAfterWait)
  }
  a.mx.Unlock()

  a.startTimeout()
  a.wg.Add(1)
  go func() {
    if err := fn(); err != nil {
      a.isErr = true
    }
    a.wg.Done()
  }()
}

// Run Добавляет задачу на запуск с каналом ожидания при завершении работы
func (a *mgr) Run(fn func() (chan struct{}, error)) {
  a.mx.Lock()
  if a.wait {
    a.log.Panic(ErrAfterWait)
  }

  l := len(a.chansExit)
  a.chansExit = append(a.chansExit, nil)

  a.mx.Unlock()

  a.startTimeout()
  a.wg.Add(1)

  go func() {
    ch, err := fn()
    if err == nil {
      a.chansExit[l] = ch
    } else {
      a.isErr = true
    }
    a.wg.Done()
  }()
}

// Wait ожидание запуска задач
func (a *mgr) Wait() {
  a.mx.Lock()
  if a.wait {
    a.log.Panic(ErrWaitOnce)
  }
  a.wait = true
  a.mx.Unlock()

  if a.IsCli {
    go a.cli()
  }

  go a.performWaitTasks()

  <-a.ctx.Done()
  a.log.Info("Application stopping")

  if a.ShutdownTimeoutMs > 0 || len(a.chansExit) > 0 {
    ch := make(chan struct{}, 1)
    do := true

    if a.IsCli {
      go func() {
        for do {
          fmt.Print(".")
          time.Sleep(time.Millisecond * delayRepeatDot)
        }
        close(ch)
      }()
    } else {
      close(ch)
    }

    select {
    case <-a.waitingCompletion():
    case <-a.sig:
    }

    if a.IsCli {
      fmt.Println()
    }

    do = false
    <-ch
  }

  if a.isErr {
    syscall.Exit(1)
  }
  syscall.Exit(0)
}

func (a *mgr) waitingCompletion() chan struct{} {
  ch := make(chan struct{})
  chEx := make(chan struct{})
  chTm := make(chan struct{})
  l := len(a.chansExit)

  if l > 0 {
    gr := sync.WaitGroup{}
    gr.Add(l)

    for _, c := range a.chansExit {
      go func() {
        <-c
        gr.Done()
      }()
    }

    gr.Wait()
    close(chEx)
  }

  if a.ShutdownTimeoutMs > 0 {
    go func() {
      <-time.After(a.ShutdownTimeoutMs)
      close(chTm)
    }()
  }

  if l > 0 || a.ShutdownTimeoutMs > 0 {
    go func() {
      select {
      case <-chEx:
      case <-chTm:
      }
      close(ch)
    }()
  } else {
    close(ch)
  }

  return ch
}
