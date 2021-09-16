package logger

import (
  "errors"
  "github.com/sirupsen/logrus"
  "time"
)

type Logger struct {
  listEntry  map[string]*sys2
  level      logrus.Level
  mode       ModeOut
  cli        bool
  timeFormat string
  sysName    string // Имя поля с названием системы
}

type sys2 struct {
  isSetLev  bool // Установлен ли определенный уровень логирования
  isSetMode bool // Установлен ли режим вывода логов
  mode      ModeOut
  isSetCli  bool // Установлен ли режим вывода в консоль
  cli       bool

  ent *logrus.Entry
}

type ModeOut int32

const (
  ModeIntJson ModeOut = iota
  ModeIntText
  ModeIntShort
  defTimeFormat   = time.RFC3339
  defLevel        = logrus.TraceLevel
  defMode         = ModeIntJson
  defSysFieldName = "sys"
)

var ModeKeyVal = map[string]ModeOut{
  ModeValKey[ModeIntJson]:  ModeIntJson,
  ModeValKey[ModeIntText]:  ModeIntText,
  ModeValKey[ModeIntShort]: ModeIntShort,
}

var ModeValKey = map[ModeOut]string{
  ModeIntJson:  "json",
  ModeIntText:  "text",
  ModeIntShort: "short",
}

var (
  ErrParseMode    = errors.New("logger: ошибка парсинга режима вывода логов")
  ErrEmptySysName = errors.New("logger: имя системы не может быть пустым")
  ErrSysNotFound  = errors.New("logger: система не найдена")
)

func New() *Logger {
  return &Logger{
    listEntry:  map[string]*sys2{},
    timeFormat: defTimeFormat,
    mode:       defMode,
    level:      defLevel,
    sysName:    defSysFieldName,
  }
}

func (l *Logger) GetLog(n string) *logrus.Entry {
  it, ok := l.listEntry[n]
  if !ok {
    it = &sys2{
      ent: l.createEntry(),
    }
    l.listEntry[n] = it
    it.ent.Logger.SetFormatter(l.formatter(it))
  }
  if n != "" {
    it.ent.Data[l.sysName] = n
  }
  return it.ent
}

func (l *Logger) createEntry() *logrus.Entry {
  return logrus.NewEntry(logrus.New())
}

// SysFieldName установить название поля с именем системы
func (l *Logger) SysFieldName(n string) {
  for _, sys := range l.listEntry {
    oldName, ok := sys.ent.Data[l.sysName]
    if ok {
      delete(sys.ent.Data, l.sysName)
      sys.ent.Data[n] = oldName
    }
  }
  l.sysName = n
}
