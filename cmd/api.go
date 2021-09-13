package cmd

import (
  "context"
  "errors"
  "github.com/sirupsen/logrus"
  "net/http"
  "sync"
)

type Cmd struct {
  port   int
  server *http.Server
  routes map[string]*Route
  mx     sync.Mutex
  log    *logrus.Entry
  cxt    context.Context
}

var (
  ErrPathExist         = errors.New("cmd: URL path already exists")
  ErrEmptyHandler      = errors.New("cmd: Handler cannot be nil")
  ErrInvalidDataFormat = errors.New("cmd: Invalid data format")
)

func New(ctx context.Context) *Cmd {
  a := &Cmd{
    routes: make(map[string]*Route),
    cxt:    ctx,
  }
  a.Route("/help", a.handlerHelp)
  return a
}

func (a *Cmd) Ready() bool {
  return a.server != nil
}

func (a *Cmd) SetLogger(l *logrus.Entry) {
  a.log = l
}
