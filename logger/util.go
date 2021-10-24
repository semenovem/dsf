package logger

import (
  "github.com/pkg/errors"
  "strings"
)

func parseMode(m string) (modeOut, error) {
  v, ok := modeKeyVal[strings.ToLower(m)]
  if ok {
    return v, nil
  }

  return 0, ErrParseMode
}

func parseDest(m string) (destOut, error) {
  v, ok := destKeyVal[strings.ToLower(m)]
  if ok {
    return v, nil
  }

  return 0, ErrParseDest
}

func wrapErr(err error, msg string) error {
  if err == nil {
    return errors.New(msg)
  }

  return errors.Wrap(err, msg)
}
