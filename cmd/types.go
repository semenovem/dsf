package cmd

import (
  "net/http"
)

type Reply struct {
  Ok     bool   `json:"ok"`
  ErrMsg string `json:"errMsg"`
  Reply  string `json:"reply"`
}

type RouteHandler func(w http.ResponseWriter, r *http.Request) (*Reply, error)

type Route struct {
  Path    string
  Handler RouteHandler
}
