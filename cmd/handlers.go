package cmd

import (
  "encoding/json"
  "fmt"
  "net/http"
)

func (a *Cmd) NotFound(w http.ResponseWriter, r *http.Request, reply string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusNotFound)

  resp := &Reply{
    ErrMsg: fmt.Sprintf("URL path '%s' not found", r.URL.Path),
    Reply:  reply,
  }
  body, _ := json.Marshal(resp)

  body = append(body, []byte("\n")...)

  _, err := w.Write(body)
  if err != nil {
    a.log.Error(err)
  }
}

func (a *Cmd) ServiceUnavailable(w http.ResponseWriter, errMsg string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusServiceUnavailable)

  resp := &Reply{
    ErrMsg: errMsg,
  }
  body, _ := json.Marshal(resp)
  body = append(body, []byte("\n")...)

  _, err := w.Write(body)
  if err != nil {
    a.log.Error(err)
  }
}

func (a *Cmd) handlerHelp(_ http.ResponseWriter, _ *http.Request) (*Reply, error) {
  return &Reply{
    Ok:    true,
    Reply: a.routerHelpReply(),
  }, nil
}
