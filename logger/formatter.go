package lg

import (
  "github.com/sirupsen/logrus"
)

const logAdditionalKey = "p"

var defLevel = logrus.InfoLevel

var FormatterDevCli = &logrus.TextFormatter{
  FullTimestamp:    false,
  DisableSorting:   false,
  DisableTimestamp: true,

  //TimestampFormat:        "15:04:05",
  DisableLevelTruncation: false,
  PadLevelText:           false,
  SortingFunc: func(s []string) {
    iPkg := -1
    iMsg := -1
    for i, v := range s {
      switch v {
      case logAdditionalKey:
        iPkg = i
      case "msg":
        iMsg = i
      }
    }

    if iPkg >= 0 && iMsg >= 0 {
      m := s[iMsg]
      s[iMsg] = s[iPkg]
      s[iPkg] = m
    }
  },
}

func (o *op) formatter() logrus.Formatter {
  switch o.mode {
  case ModeText:
    return &Formatter{
      TrimMessages:     true,
      HideKeys:         false,
      DisableTimestamp: false,
      FieldsOrder:      []string{"p"},
      TimestampFormat:  o.timeFormat,
      NoColors:         !o.cli,
    }
  case ModeShortText:
    return &Formatter{
      TrimMessages:     true,
      HideKeys:         false,
      DisableTimestamp: true,
      FieldsOrder:      []string{"p"},
      TimestampFormat:  o.timeFormat,
      NoColors:         !o.cli,
    }
  }

  return &logrus.JSONFormatter{}
}
