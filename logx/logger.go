package logx

import (
	"fmt"
	"github.com/lamgor666/goboot-common/util/errorx"
	"github.com/sirupsen/logrus"
	"strings"
)

type Logger interface {
	Log(level interface{}, args ...interface{})
	Logf(level interface{}, format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type loggerImpl struct {
	channel string
	logger  *logrus.Logger
}

func (l *loggerImpl) Log(level interface{}, args ...interface{}) {
	logLevel := logrus.TraceLevel

	switch t := level.(type) {
	case logrus.Level:
		logLevel = t
	case string:
		lv := strings.ToLower(t)

		switch strings.ToLower(lv) {
		case "debug":
			logLevel = logrus.DebugLevel
		case "info":
			logLevel = logrus.InfoLevel
		case "warn":
			logLevel = logrus.WarnLevel
		case "error":
			logLevel = logrus.ErrorLevel
		case "panic":
			logLevel = logrus.PanicLevel
		case "fatal":
			logLevel = logrus.FatalLevel
		}
	}

	argList := make([]interface{}, 0)

	for _, v := range args {
		if err, ok := v.(error); ok {
			argList = append(argList, errorx.Stacktrace(err))
			continue
		}

		argList = append(argList, v)
	}

	entry := l.logger.WithField("channel", l.channel)
	entry.Writer()

	switch logLevel {
	case logrus.DebugLevel:
		entry.Debug(argList...)
	case logrus.InfoLevel:
		entry.Info(argList...)
	case logrus.WarnLevel:
		entry.Warn(argList...)
	case logrus.ErrorLevel:
		entry.Error(argList...)
	case logrus.PanicLevel:
		entry.Panic(argList...)
	case logrus.FatalLevel:
		entry.Fatal(argList...)
	default:
		entry.Trace(argList...)
	}
}

func (l *loggerImpl) Logf(level interface{}, format string, args ...interface{}) {
	var msg string

	if len(args) < 1 {
		msg = format
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	l.Log(level, msg)
}

func (l *loggerImpl) Trace(args ...interface{}) {
	l.Log("trace", args...)
}

func (l *loggerImpl) Tracef(format string, args ...interface{}) {
	l.Logf("trace", format, args...)
}

func (l *loggerImpl) Debug(args ...interface{}) {
	l.Log("debug", args...)
}

func (l *loggerImpl) Debugf(format string, args ...interface{}) {
	l.Logf("debug", format, args...)
}

func (l *loggerImpl) Info(args ...interface{}) {
	l.Log("info", args...)
}

func (l *loggerImpl) Infof(format string, args ...interface{}) {
	l.Logf("info", format, args...)
}

func (l *loggerImpl) Warn(args ...interface{}) {
	l.Log("warn", args...)
}

func (l *loggerImpl) Warnf(format string, args ...interface{}) {
	l.Logf("warn", format, args...)
}

func (l *loggerImpl) Error(args ...interface{}) {
	l.Log("error", args...)
}

func (l *loggerImpl) Errorf(format string, args ...interface{}) {
	l.Logf("error", format, args...)
}

func (l *loggerImpl) Panic(args ...interface{}) {
	l.Log("panic", args...)
}

func (l *loggerImpl) Panicf(format string, args ...interface{}) {
	l.Logf("panic", format, args...)
}

func (l *loggerImpl) Fatal(args ...interface{}) {
	l.Log("fatal", args...)
}

func (l *loggerImpl) Fatalf(format string, args ...interface{}) {
	l.Logf("fatal", format, args...)
}
