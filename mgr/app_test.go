package mgr

import (
  "fmt"
  "github.com/sirupsen/logrus"
  "testing"
  "time"
)

func testNew() *mgr {
  return New(logrus.NewEntry(logrus.New()))
}

func TestTask(t  *testing.T) {
  man := testNew()


  task := func() (chan struct{}, error) {
    time.Sleep(time.Millisecond * 10)
    ch := make(chan struct{})

    go func() {
      <- man.Ctx.Done()
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
  man := testNew()

  task := func() (chan struct{}, error) {
    time.Sleep(time.Millisecond * 10)
    ch := make(chan struct{})

    go func() {
      <- man.Ctx.Done()
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
