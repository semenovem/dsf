package cmdapi

import "errors"

var (
  ErrAlreadyStarted = errors.New("cmd: server already started")
)
