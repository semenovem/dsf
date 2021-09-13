package lg

import (
  "github.com/sirupsen/logrus"
  "strings"
)

type op struct {
  ent        *logrus.Entry
  dev        bool
  timeFormat string
  mode       string
  lev        logrus.Level
  pkgName    string
  cli        bool
}

const (
  ModeJson      = "json"
  ModeText      = "text"
  ModeShortText = "short"
)

func With(ent *logrus.Entry) *op {
  return &op{
    ent:        ent,
    timeFormat: defaultTimeFormat,
    mode:       defaultMode,
    lev:        defaultLev,
    dev:        defaultDev,
    cli:        defaultCli,
  }
}

func New() *logrus.Entry {
  return newFrom().Apl()
}

func newFrom() *op {
  return With(logrus.NewEntry(logrus.New()))
}

// Set первым аргументом - название системы
// второй и далее в любом порядке - уровень логирования / тип лога
func Set(args ...string) *logrus.Entry {
  n := newFrom()

  if len(args) == 0 {
    return n.Apl()
  }

  if args[0] != "" {
    n.Pkg(args[0])
  }

  for i := 1; i < len(args); i++ {
    lev, err := logrus.ParseLevel(args[i])
    if err == nil {
      n.setLev(lev)
      continue
    }
    mod := strings.ToLower(args[i])
    if isMod(strings.ToLower(mod)) {
      n.mode = mod
    }
  }

  return n.Apl()
}

func New1(pkg string) *logrus.Entry {
  return newFrom().Pkg(pkg).Apl()
}

func New2(pkg string, lev string) *logrus.Entry {
  return newFrom().Pkg(pkg).Lev(lev).Apl()
}

// Pkg Добавляет в лог имя пакета
func (o *op) Pkg(n string) *op {
  o.pkgName = n
  return o
}

// Lev Установить уровень логгирования
func (o *op) Lev(s string) *op {
  if l, err := logrus.ParseLevel(s); err == nil {
    o.lev = l
  }
  return o
}

// Lev Установить уровень логгирования
func (o *op) setLev(lev logrus.Level) {
  o.lev = lev
}

func (o *op) Apl() *logrus.Entry {
  o.ent.Logger.SetLevel(o.lev)
  o.ent.Logger.SetFormatter(o.formatter())

  if o.pkgName != "" {
    o.ent.Data[logAdditionalKey] = o.pkgName
  }

  return o.ent
}

func (o *op) Dev(on bool) *op {
  o.dev = on
  return o
}

func (o *op) Mod(m string) *op {
  o.mode = m
  return o
}
