package cmd

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
)

func (a *Api) routerHelpReply() string {
  keys := make([]string, 0)
  for it, _ := range a.routes {
    keys = append(keys, it)
  }
  return fmt.Sprintf("Available URL: [%s]", strings.Join(keys, ", "))
}

func (a *Api) Route(path string, handler RouteHandler) {
  _, has := a.routes[path]
  if has {
    a.log.Panic("URL path already exists")
  }

  if handler == nil {
    a.log.Panic("Handler cannot be nil")
  }

  a.routes[path] = &Route{
    Path:    path,
    Handler: handler,
  }
}

func (a *Api) DelRoute(path string) {
  delete(a.routes, path)
}

func (a *Api) HasRoute(path string) bool {
  _, has := a.routes[path]
  return has
}

func (a *Api) RouteReply(route string, fn func(*Reply) error) {
  a.Route(route, func(w http.ResponseWriter, r *http.Request) (*Reply, error) {
    reply := &Reply{}
    err := fn(reply)
    return reply, err
  })
}
func (a *Api) RouteFn(route string, fn func() error) {
  a.Route(route, func(w http.ResponseWriter, r *http.Request) (*Reply, error) {
    reply := &Reply{
      Ok: true,
    }
    err := fn()
    if err != nil {
      reply.Ok = false
      reply.ErrMsg = err.Error()
    }
    return reply, err
  })
}

func (a *Api) RouteReqReply(route string, fn func(*http.Request, *Reply) error) {
  a.Route(route, func(w http.ResponseWriter, r *http.Request) (*Reply, error) {
    reply := &Reply{}
    err := fn(r, reply)
    return reply, err
  })
}

func (a *Api) RouteText(route string, v string) {
  a.Route(route, func(_ http.ResponseWriter, _ *http.Request) (*Reply, error) {
    return &Reply{
      Ok:    true,
      Reply: v,
    }, nil
  })
}

func (a *Api) RouteHealth(route string) {
  a.Route(route, func(w http.ResponseWriter, r *http.Request) (*Reply, error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    return nil, nil
  })
}

func (a *Api) router(w http.ResponseWriter, r *http.Request) {
  a.log.Info("API IN", r.URL.Path)

  path := strings.ToLower(r.URL.Path)
  pathLen := len(path)

  var matchLen int
  var route *Route

  for p, r := range a.routes {
    if strings.HasPrefix(path, p) && len(p) > matchLen {
      if pathLen > len(p) && string(path[len(p)]) != "/" {
        continue
      }
      route = r
      matchLen = len(p)
    }
  }

  if route != nil {
    a.routeServe(w, r, route.Handler)

    return
  }

  a.NotFound(w, r, a.routerHelpReply())
}

func (a *Api) routeServe(w http.ResponseWriter, r *http.Request, h RouteHandler) {
  resp, err := h(w, r)

  if err != nil {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    body, _ := json.Marshal(&Reply{
      ErrMsg: err.Error(),
    })
    body = append(body, []byte("\n")...)
    _, err := w.Write(body)
    if err != nil {
      a.log.Error(err)
    }

    return
  }

  if resp != nil {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    body, _ := json.Marshal(resp)
    _, err := w.Write(body)
    if err != nil {
      a.log.Error(err)
    }

    _, err = w.Write([]byte("\n"))
    if err != nil {
      a.log.Error(err)
    }
  }
}
