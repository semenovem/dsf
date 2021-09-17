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
  ctx    context.Context
}

var (
  ErrPathExist         = errors.New("cmd: URL path already exists")
  ErrEmptyHandler      = errors.New("cmd: Handler cannot be nil")
  ErrInvalidDataFormat = errors.New("cmd: Invalid data format")
)

func New(ctx context.Context, l *logrus.Entry) *Cmd {
  a := &Cmd{
    routes: make(map[string]*Route),
    ctx:    ctx,
    log:    l,
  }
  a.Route("/help", a.handlerHelp)

  return a
}

func (a *Cmd) Ready() bool {
  return a.server != nil
}
