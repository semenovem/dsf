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
  Log               *logrus.Entry
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
  timer             *time.Timer
  fnsStarted        []func()
  fnsAfterTasks     []func()
  fnsFailed         []func() // Подписка на событие ошибки запуска
  sig               chan os.Signal
}

func New() *mgr {
  ctx, cancel := context.WithCancel(context.Background())

  o := &mgr{
    Ctx:               ctx,
    ctxCancel:         cancel,
    Log:               logrus.NewEntry(logrus.New()),
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
        a.Log.Warn("Application launch was stopped due to timeout")
        a.ctxCancel()
      }
    }
    a.timer = time.AfterFunc(a.Timeout, fn)
  }
}

// Task Добавляет задачу на запуск
func (a *mgr) Task(fn func() error) {
  if a.wait {
    a.Log.Panic(ErrAfterWait)
  }
  a.startTimeout()
  a.wg.Add(1)
  go func() {
    if err := fn(); err != nil {
      a.isErr = true
      a.Log.Errorf("Task launch error: %v", err)
    }
    a.wg.Done()
  }()
}

// Wait ожидание запуска задач
func (a *mgr) Wait() {
  if a.wait {
    a.Log.Panic(ErrWaitOnce)
  }

  a.wait = true

  if a.IsCli {
    go func() {
      // TODO change to receive a scan code of a button
      err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
      if err != nil {
        a.Log.Warn(err)
        return
      }
      err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
      if err != nil {
        a.Log.Warn(err)
        return
      }
      var b = make([]byte, 1)
      for {
        _, err = os.Stdin.Read(b)
        if err != nil {
          a.Log.Warn(err)
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

  if a.ShutdownTimeoutMs > 0 {
    a.Log.Infof("Application stopping")
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
    case <-time.After(a.ShutdownTimeoutMs):
    case <-a.sig:
    }

    if a.IsCli {
      fmt.Println()
    }

    do = false
    <-ch
  }
}
