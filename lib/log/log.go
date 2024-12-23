package log

import (
	"fmt"
	"io"
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

var SelLogLevel = Warning

type Logger struct {
	Parent *Logger
	Name   string
	Writer io.Writer
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

func (l *Logger) log(str string, level LogLevel, wovr io.Writer) {
	if wovr != nil {
		fmt.Fprintf(wovr, "%s%s %s\n", l.getPrefix(), level, str)
	} else {
		fmt.Fprintf(l.Writer, "%s%s %s\n", l.getPrefix(), level, str)
	}
}

func (l *Logger) Logln(level LogLevel, s ...any) {
	if level < SelLogLevel {
		return
	}

	comb := fmt.Sprint(s...)
	l.log(comb, level, nil)
}

func (l *Logger) Logf(level LogLevel, format string, s ...any) {
	if level < SelLogLevel {
		return
	}

	comb := fmt.Sprintf(format, s...)
	l.log(comb, level, nil)
}

func (l *Logger) Fatalln(s ...any) {
	comb := fmt.Sprint(s...)
	l.log(comb, Error, os.Stderr)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, s ...any) {
	comb := fmt.Sprintf(format, s...)
	l.log(comb, Error, os.Stderr)
	os.Exit(1)
}

func (l *Logger) ChildLogger(name string) *Logger {
	nl := &Logger{
		Parent: l,
		Name:   name,
		Writer: l.Writer,
	}

	return nl
}

func (l *Logger) OverrideWriter(w io.Writer) *Logger {
	l.Writer = w
	return l
}

var Default = &Logger{Writer: os.Stdout}

var Logln = Default.Logln
var Logf = Default.Logf
var Fatalln = Default.Fatalln
var Fatalf = Default.Fatalf
