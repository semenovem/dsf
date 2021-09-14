package logger

import (
  "github.com/sirupsen/logrus"
)

func (l *Logger) formatter(it *sys2) logrus.Formatter {
  cli := false

  if it.isSetCli {
    cli = it.cli
  } else {
    cli = l.cli
  }

  var mode ModeOut
  if it.isSetMode {
    mode = it.mode
  } else {
    mode = l.mode
  }

  switch mode {
  case ModeIntText:
    return &Formatter{
      TrimMessages:     true,
      HideKeys:         false,
      DisableTimestamp: false,
      FieldsOrder:      []string{defSysName},
      TimestampFormat:  l.timeFormat,
      NoColors:         !cli,
    }
  case ModeIntShort:
    return &Formatter{
      TrimMessages:     true,
      HideKeys:         false,
      DisableTimestamp: true,
      FieldsOrder:      []string{defSysName},
      TimestampFormat:  l.timeFormat,
      NoColors:         !cli,
    }
  }

  return &logrus.JSONFormatter{}
}
