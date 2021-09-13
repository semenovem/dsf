package lg

import (
  "github.com/sirupsen/logrus"
  "strings"
)

func ParseLev(s string) logrus.Level {
  if l, err := logrus.ParseLevel(s); err == nil {
    return l
  }
  return defaultLev
}

func ParseMod(m string, defMod string) string {
  if isMod(strings.ToLower(m)) {
    return m
  }
  if isMod(strings.ToLower(defMod)) {
    return defMod
  }

  panic("The value is not a valid mode")
}

func isMod(m string) bool {
  switch m {
  case ModeJson, ModeText, ModeShortText:
    return true
  }

  return false
}
