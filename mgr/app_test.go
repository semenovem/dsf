package mgr

import (
  "context"
  "fmt"
  "github.com/sirupsen/logrus"
  "testing"
  "time"
)

func TestTask(t *testing.T) {
  ctx, cancel := context.WithCancel(context.Background())
  man := New(ctx, cancel, logrus.NewEntry(logrus.New()))

  task := func() (chan struct{}, error) {
    time.Sleep(time.Millisecond * 10)
    ch := make(chan struct{})

    go func() {
      <-ctx.Done()
      fmt.Println("___ function started exiting")
      time.Sleep(time.Millisecond * 3000)
      close(ch)
    }()

    return ch, nil
  }

  man.RegisterStarted(func() {
    fmt.Println(">>>>>>>>>> tasks is done")
  })

  man.Run(task)

  man.Wait()
}

func TestWaitingCompletion(t *testing.T) {
  ctx, cancel := context.WithCancel(context.Background())
  man := New(ctx, cancel, logrus.NewEntry(logrus.New()))

  task := func() (chan struct{}, error) {
    time.Sleep(time.Millisecond * 10)
    ch := make(chan struct{})

    go func() {
      <-ctx.Done()
      fmt.Println("___ function started exiting")
      time.Sleep(time.Millisecond * 3000)
      close(ch)
    }()

    return ch, nil
  }

  man.RegisterStarted(func() {
    fmt.Println(">>>>>>>>>> tasks is done")
  })

  man.Run(task)

  man.Wait()
}
