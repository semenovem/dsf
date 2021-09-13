package cmd

import (
  "net/http"
)

type Reply struct {
  Ok     bool        `json:"ok"`
  ErrMsg string      `json:"errMsg"`
  Reply  interface{} `json:"reply"`
}

type RouteHandler func(w http.ResponseWriter, r *http.Request) (*Reply, error)

type Route struct {
  Path    string
  Handler RouteHandler
}

func (r *Reply) Err(err interface{}) {
  switch e := err.(type) {
  case string:
    r.ErrMsg = e
  case error:
    r.ErrMsg = e.Error()
  default:
    panic(ErrInvalidDataFormat)
  }
}

func (r *Reply) ErrCritical(err error) {
  r.ErrMsg = err.Error()
}

func (r *Reply) Payload(rep interface{}) {
  r.Reply = rep
}

func (r *Reply) Success() {
  r.Ok = true
}

func (r *Reply) SetOk(v bool) {
  r.Ok = v
}
