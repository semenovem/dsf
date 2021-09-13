package lg

import (
  "github.com/sirupsen/logrus"
  "time"
)

var (
  defaultTimeFormat = time.RFC3339
  defaultDev        = false
  defaultMode       = ModeJson
  defaultLev        = logrus.TraceLevel
  defaultCli        = false
)

// SetDev флаг разработки
func SetDev(on bool) {
  defaultDev = on
}

func GetDev() bool {
  return defaultDev
}

// SetLev Уровень логгирования
func SetLev(l logrus.Level) {
  defaultLev = l
}

func GetDefLev() logrus.Level {
  return defaultLev
}

// SetMod Режим логов
func SetMod(m string) {
  if !isMod(m) {
    panic("The value is not a valid mode")
  }
  defaultMode = m
}

func GetDefMod() string {
  return defaultMode
}

// SetCli Режим вывода логов в консоль с подсветкой
func SetCli(on bool) {
  defaultCli = on
}

func GetDefCLI() bool {
  return defaultCli
}
