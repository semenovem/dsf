package cmd

import (
  "fmt"
  "net/http"
  "strconv"
)

func (a *Cmd) SetPort(port int) {
  a.port = port
}

func (a *Cmd) Start() error {
  a.mx.Lock()
  defer a.mx.Unlock()

  if a.server != nil {
    return nil
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

  go func() {
    defer a.Stop()

    if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      a.log.Error("API server error: ", err)
    }
  }()

  a.log.Infof("API server started on port %d", a.port)

  return nil
}

func (a *Cmd) Stop() {
  a.mx.Lock()
  defer a.mx.Unlock()

  if a.server != nil {
    err := a.server.Close()
    a.server = nil
    if err != nil {
      a.log.Error("API server stop error")
    }
    a.log.Info("API server stopped")
  }
}
