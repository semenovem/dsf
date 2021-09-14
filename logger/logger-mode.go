package logger

func (l *Logger) SetDefMode(s string) error {
  p, err := ParseMod(s)
  if err == nil {
    l.mode = p
    l.distrDefMode()
  }
  return err
}

func (l *Logger) GetDefMode() ModeOut {
  return l.mode
}

func (l *Logger) distrDefMode() {
  for _, it := range l.listEntry {
    if !it.isSetMode {
      it.ent.Logger.SetFormatter(l.formatter(it))
    }
  }
}
