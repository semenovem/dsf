package logger

import "github.com/sirupsen/logrus"

func (l *Logger) SetDefLevel(s string) error {
  p, err := logrus.ParseLevel(s)
  if err == nil {
    l.level = p
    l.distrDefLevel()
  }
  return err
}

func (l *Logger) SetLevel(name, lev string) error {
  v, err := logrus.ParseLevel(lev)
  if err == nil {
    sys, ok := l.listEntry[name]
    if !ok {
      return ErrSysNotFound
    }
    sys.isSetLev = true
    sys.ent.Logger.SetLevel(v)

    l.level = v
  }
  return err
}

func (l *Logger) GetDefLevel() logrus.Level {
  return l.level
}

func (l *Logger) distrDefLevel() {
  for _, it := range l.listEntry {
    if !it.isSetLev {
      it.ent.Logger.SetLevel(l.level)
    }
  }
}
