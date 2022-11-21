package logger

import (
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Logger struct {
	listEntry  map[string]*sys
	level      logrus.Level
	mode       modeOut
	dest       destOut
	timeFormat string
	sysName    string // Имя поля с названием системы
}

type sys struct {
	isSetLev bool // Установлен ли определенный уровень логирования
	mode     modeOut
	dest     destOut
	ent      *logrus.Entry
}

type modeOut uint
type destOut uint

const (
	modeUnkn modeOut = iota
	modeJson
	modeText
	modeShort

	defTimeFormat   = time.RFC3339
	defLevel        = logrus.TraceLevel
	defMode         = modeJson
	defDest         = destFile
	defSysFieldName = "sys"
)

const (
	destUnkn destOut = iota
	destFile
	destCli
)

var modeKeyVal = map[string]modeOut{
	modeValKey[modeJson]:  modeJson,
	modeValKey[modeText]:  modeText,
	modeValKey[modeShort]: modeShort,
}

var modeValKey = map[modeOut]string{
	modeJson:  "json",
	modeText:  "text",
	modeShort: "short",
}

var destKeyVal = map[string]destOut{
	destValKey[destFile]: destFile,
	destValKey[destCli]:  destCli,
}

var destValKey = map[destOut]string{
	destFile: "file",
	destCli:  "console",
}

var (
	ErrParseMode   = errors.New("logger: ошибка парсинга режима вывода логов")
	ErrParseDest   = errors.New("logger: ошибка парсинга места записи логов")
	ErrSysNotFound = errors.New("logger: система не найдена")
)

func New() *Logger {
	return &Logger{
		listEntry:  map[string]*sys{},
		timeFormat: defTimeFormat,
		mode:       defMode,
		dest:       defDest,
		level:      defLevel,
		sysName:    defSysFieldName,
	}
}

func (l *Logger) GetLog(n string) *logrus.Entry {
	return l.getOrCreateSys(n).ent
}

func (l *Logger) getOrCreateSys(n string) *sys {
	it, ok := l.listEntry[n]
	if !ok {
		it = &sys{
			ent: l.createEntry(),
		}
		l.listEntry[n] = it
		it.ent.Logger.SetFormatter(l.formatter(it))
		it.ent.Logger.SetLevel(l.level)
	}
	if n != "" {
		it.ent.Data[l.sysName] = n
	}
	return it
}

func (l *Logger) createEntry() *logrus.Entry {
	return logrus.NewEntry(logrus.New())
}

// aplDef применить дефолтные установки
func (l *Logger) aplDef() {
	for _, it := range l.listEntry {
		it.ent.Logger.SetFormatter(l.formatter(it))
	}
}

func (l *Logger) SetLevel(name, lev string) error {
	if lev == "" {
		return nil
	}
	v, err := logrus.ParseLevel(lev)
	if err == nil {
		it := l.getOrCreateSys(name)
		it.isSetLev = true
		it.ent.Logger.SetLevel(v)
	}
	return err
}

func (l *Logger) distrDefLevel() {
	for _, it := range l.listEntry {
		if !it.isSetLev {
			it.ent.Logger.SetLevel(l.level)
		}
	}
}

func (l *Logger) SetDef(lev, mode, dest string) error {
	var (
		err0 error
		apl  bool
	)
	if lev != "" {
		p, err := logrus.ParseLevel(lev)
		if err == nil {
			l.level = p
			l.distrDefLevel()
		} else {
			err0 = wrapErr(err0, err.Error())
		}
	}
	if mode != "" {
		m, err := parseMode(mode)
		if err == nil {
			l.mode = m
			apl = true
		} else {
			err0 = wrapErr(err0, err.Error())
		}
	}
	if dest != "" {
		d, err := parseDest(dest)
		if err == nil {
			l.dest = d
			apl = true
		} else {
			err0 = wrapErr(err0, err.Error())
		}
	}
	if apl {
		for _, it := range l.listEntry {
			it.ent.Logger.SetFormatter(l.formatter(it))
		}
	}
	return err0
}

func (l *Logger) Set(name, lev, mode, dest string) error {
	var (
		err0 error
		apl  bool
	)

	it, ok := l.listEntry[name]
	if !ok {
		return ErrSysNotFound
	}

	if lev != "" {
		l, err := logrus.ParseLevel(lev)
		if err == nil {
			it.ent.Logger.SetLevel(l)
			it.isSetLev = true
		} else {
			err0 = wrapErr(err0, err.Error())
		}
	}
	if mode != "" {
		m, err := parseMode(mode)
		if err == nil {
			it.mode = m
			apl = true
		} else {
			err0 = wrapErr(err0, err.Error())
		}
	}
	if dest != "" {
		d, err := parseDest(dest)
		if err == nil {
			it.dest = d
			apl = true
		} else {
			err0 = wrapErr(err0, err.Error())
		}
	}
	if apl {
		it.ent.Logger.SetFormatter(l.formatter(it))
	}
	return err0
}
