package apptool

import (
  "context"
  "github.com/sirupsen/logrus"
  "os"
  "os/signal"
  "sync"
  "syscall"
  "time"
)

type apptool struct {
  Timeout        time.Duration // Таймаут ожидания 0 - без таймаута
  Log            *logrus.Entry
  Ctx            context.Context
  ctxCancel      context.CancelFunc
  tasksCompleted bool
  failFn         func() // Функция при ошибке в одной из задач
  wait           bool   // Флаг начала ожидания работы приложения
  isErr          bool   // Одна из задач завершилась с ошибкой
  wg             sync.WaitGroup
  timer          *time.Timer
  fnsStarted     []func()
  fnsAfterTasks  []func()
  fnsFailed      []func() // подписка на событие ошибки запуска
}

func New() *apptool {
  ctx, cancel := context.WithCancel(context.Background())

  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sig
    cancel()
  }()

  return &apptool{
    Ctx:       ctx,
    ctxCancel: cancel,
    Log:       logrus.NewEntry(logrus.New()),
  }
}

func (a *apptool) Exit() {
  a.ctxCancel()
}

func (a *apptool) startTimeout() {
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
func (a *apptool) Task(fn func() error) {
  if a.wait {
    a.Log.Panic("Task cannot be added after Wait")
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
func (a *apptool) Wait() {
  a.wait = true

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
  }()

  <-a.Ctx.Done()
}
