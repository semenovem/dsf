package cmd

import (
  "github.com/sirupsen/logrus"
  "net/http"
  "sync"
)

type Api struct {
  port   int
  server *http.Server
  routes map[string]*Route
  mx     sync.Mutex
  log    *logrus.Entry
}

func Create() *Api {
  a := &Api{
    routes: make(map[string]*Route),
  }

  a.Route("/help", a.handlerHelp)

  return a
}

func (a *Api) Ready() bool {
  return a.server != nil
}

func (a *Api) SetLogger(l *logrus.Entry) {
  a.log = l
}
