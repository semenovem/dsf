package cmd

import (
  "fmt"
  "net/http"
  "strconv"
  "time"
)

func (a *Cmd) SetPort(port int) {
  a.port = port
}

func (a *Cmd) Start() error {
  a.mx.Lock()
  defer a.mx.Unlock()

  if a.server != nil {
    a.log.Error(ErrAlreadyStarted)
    return ErrAlreadyStarted
  }

  if a.port < 1 {
    msg := fmt.Sprintf("Invalid value 'port': %d", a.port)
    a.log.Error(msg)
    return fmt.Errorf(msg)
  }

  a.server = &http.Server{
    Addr:    ":" + strconv.Itoa(a.port),
    Handler: http.HandlerFunc(a.router),
  }

  a.log.Infof("API server started on port %d", a.port)

  var err error

  go func() {
    defer a.Stop()

    if err = a.server.ListenAndServe(); err != nil {
      if err == http.ErrServerClosed {
        err = nil
      } else {
        a.log.Error("API server error: ", err)
      }
    }
  }()

  select {
  case <-time.After(time.Millisecond * 500):
  case <-a.ctx.Done():
  }

  if err == nil {
    go func() {
      <-a.ctx.Done()
      a.Stop()
    }()
  }

  return err
}

func (a *Cmd) Stop() {
  a.mx.Lock()
  defer a.mx.Unlock()
  a.stop()
}

func (a *Cmd) stop() {
  if a.server != nil {
    err := a.server.Close()
    a.server = nil
    if err != nil {
      a.log.Error("API server stop error")
    }
    a.log.Info("API server stopped")
  }
}
