package log

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warning
	Error
)

func (l LogLevel) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func ParseLogLevel(s string) (LogLevel, error) {
	switch s {
	case "debug", "DEBUG":
		return Debug, nil
	case "info", "INFO":
		return Info, nil
	case "warn", "WARN":
		return Warning, nil
	case "error", "ERROR":
		return Error, nil
	default:
		return SelLogLevel, fmt.Errorf("unknown log level %s", s)
	}
}

func SetLogLevel(l LogLevel) {
	SelLogLevel = l
}

var SelLogLevel = Info

type Logger struct {
	Parent *Logger
	Name   string
}

func (l *Logger) prefix() []string {
	if l.Parent == nil {
		return []string{l.Name}
	}
	return append(l.Parent.prefix(), l.Name)
}

func (l *Logger) getPrefix() string {
	s := ""
	for _, p := range l.prefix() {
		if p == "" {
			continue
		}
		s += p + "."
	}
	if s == "" {
		return ""
	}
	return s[:len(s)-1] + " "
}

func (l *Logger) Logln(level LogLevel, s ...any) {
	if level < SelLogLevel {
		return
	}

	comb := fmt.Sprint(s...)
	fmt.Printf("%s%s %s\n", l.getPrefix(), level, comb)
}

func (l *Logger) Logf(level LogLevel, format string, s ...any) {
	if level < SelLogLevel {
		return
	}

	comb := fmt.Sprintf(format, s...)
	fmt.Printf("%s%s %s\n", l.getPrefix(), level, comb)
}

func (l *Logger) Fatalln(s ...any) {
	l.Logln(Error, s...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, s ...any) {
	l.Logf(Error, format, s...)
	os.Exit(1)
}

func (l *Logger) ChildLogger(name string) *Logger {
	return &Logger{
		Parent: l,
		Name:   name,
	}
}

var Default = &Logger{}

var Logln = Default.Logln
var Logf = Default.Logf
var Fatalln = Default.Fatalln
var Fatalf = Default.Fatalf
