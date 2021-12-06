package mgr

import (
  "fmt"
  "os"
  "os/exec"
  "sync"
)

func (a *mgr) cli() {
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
    switch b[0] {
    case 153, 208, 185, 113, 81:
      a.ctxCancel()
    }
  }
}

// Ожидает выполнения задач запуска
func (a *mgr) performWaitTasks() {
  a.wg.Wait()
  a.tasksCompleted = true
  if a.ctx.Err() != nil {
    return
  }

  // Задачи добавленные после task
  var wg sync.WaitGroup

  for _, fn := range a.fnsAfterTasks {
    wg.Add(1)
    go func(fn func() error) {
      defer wg.Done()
      if fn() != nil {
        a.isErr = true
      }
    }(fn)
  }
  a.fnsAfterTasks = nil

  if a.timer != nil {
    a.timer.Stop()
    a.timer = nil
  }

  if a.isErr {
    a.fireFailed()
    a.ctxCancel()
    return
  }
  a.fireStarted()
  a.started = true
}
