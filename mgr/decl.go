package mgr

import (
  "errors"
  "time"
)

var (
  ErrAfterWait = errors.New("mgr: task cannot be added after Wait")
  ErrWaitOnce = errors.New("mgr: Использование более одного вызова метода")
)

const delayRepeatDot = 200
const defShutdownTimeout = time.Millisecond * 1000
