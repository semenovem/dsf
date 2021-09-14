package logger

func (l *Logger) SetDefCli(s bool) {
  l.cli = s
  l.distrDefCli()
}

func (l *Logger) GetDefCli() bool {
  return l.cli
}

func (l *Logger) distrDefCli() {
  for _, it := range l.listEntry {
    if !it.isSetCli {
      it.ent.Logger.SetFormatter(l.formatter(it))
    }
  }
}

func (l *Logger) SetCli(name string, b bool) error {
  sys, ok := l.listEntry[name]
  if !ok {
    return ErrSysNotFound
  }
  sys.isSetCli = true
  sys.ent.Logger.SetFormatter(l.formatter(sys))

  return nil
}
