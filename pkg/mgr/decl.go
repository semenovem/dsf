package mgr

import (
  "errors"
  "time"
)

var (
  ErrAfterWait = errors.New("man: task cannot be added after Wait")
  ErrWaitOnce  = errors.New("man: Using more than one method call")
  ErrHandler   = errors.New("man: Handler cannot be added after Wait")
)

const delayRepeatDot = 200
const defShutdownTimeout = time.Millisecond * 1000
