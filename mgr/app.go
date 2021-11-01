package mgr

import (
  "context"
  "fmt"
  "github.com/sirupsen/logrus"
  "os"
  "os/exec"
  "os/signal"
  "sync"
  "syscall"
  "time"
)

type mgr struct {
  Timeout           time.Duration // Таймаут ожидания 0 - без таймаута
  log               *logrus.Entry
  Ctx               context.Context
  IsCli             bool          // Вывод в консоль
  ShutdownTimeoutMs time.Duration // Время ожидания закрытия приложения
  ctxCancel         context.CancelFunc
  tasksCompleted    bool
  failFn            func() // Функция при ошибке в одной из задач
  wait              bool   // Флаг начала ожидания работы приложения
  isErr             bool   // Одна из задач завершилась с ошибкой
  started           bool   // Флаг запуска всех задач
  wg                sync.WaitGroup
  mx                sync.Mutex
  timer             *time.Timer
  fnsStarted        []func()
  fnsAfterTasks     []func()
  fnsFailed         []func() // Подписка на событие ошибки запуска
  sig               chan os.Signal
  chansExit         []chan struct{} // Каналы ожидания завершения работы
}

func New(l *logrus.Entry) *mgr {
  ctx, cancel := context.WithCancel(context.Background())

  o := &mgr{
    Ctx:               ctx,
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

func (a *mgr) Exit() {
  a.ctxCancel()
}

func (a *mgr) Ready() bool {
  return a.started
}

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
      a.log.Errorf("Task launch error: %v", err)
    }
    a.wg.Done()
  }()
}

// Run Добавляет задачу на запуск
func (a *mgr) Run(fn func() (chan struct{}, error)) {
  a.mx.Lock()
  if a.wait {
    a.log.Panic(ErrAfterWait)
  }

  l := len(a.chansExit)
  fmt.Println(">>>> l = ", l)
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
    go func() {
      // TODO change to receive a scan code of a button
      err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
      if err != nil {
        a.log.Warn(err)
        return
      }
      err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
      if err != nil {
        a.log.Warn(err)
        return
      }
      var b = make([]byte, 1)
      for {
        _, err = os.Stdin.Read(b)
        if err != nil {
          a.log.Warn(err)
          return
        }

        if b[0] == 10 {
          fmt.Println()
          continue
        }

        a.ctxCancel()

        // временно отключим
        //switch b[0] {
        //case 153, 208, 185, 113, 81:
        //  a.ctxCancel()
        //}
      }
    }()
  }

  go func() {
    a.wg.Wait()
    a.tasksCompleted = true
    if a.Ctx.Err() != nil {
      return
    }

    if a.timer != nil {
      a.timer.Stop()
      a.timer = nil
    }

    if a.isErr {
      a.fireFailed()
      a.Exit()
      return
    }
    a.fireStarted()
    a.started = true
  }()

  <-a.Ctx.Done()
  a.log.Info("Application stopping")
  a.log.Info("Application stopping = " , a.ShutdownTimeoutMs)

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
      fmt.Println("**************** +++++++++ task")
      case <-chTm:
      fmt.Println("**************** +++++++++ timeout")
      }
      close(ch)
    }()
  } else {
    close(ch)
  }

  return ch
}
