package logx

import (
	"github.com/lamgor666/goboot-common/AppConf"
	"github.com/lamgor666/goboot-common/enum/RegexConst"
	"github.com/lamgor666/goboot-common/util/castx"
	"github.com/lamgor666/goboot-common/util/fsx"
	"github.com/lamgor666/goboot-common/util/slicex"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
)

var fieldSep = "~logrus.FieldSep~"
var logDir string
var globalAlyslsSettings *alyslsSettings
var loggers = map[string]Logger{}

func WithLogDir(dir string) {
	dir = fsx.GetRealpath(dir)
	
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		logDir = dir
	}
}

func WithAlyslsSettings(settings ...map[string]interface{}) {
	st := map[string]interface{}{}

	if len(settings) > 0 && len(settings[0]) > 0 {
		st = settings[0]
	}

	if len(st) < 1 {
		st = AppConf.GetMap("alysls")
	}

	globalAlyslsSettings = newAlyslsSettings(st)
}

func InitLoggers(defines ...[]map[string]interface{}) {
	entries := make([]map[string]interface{}, 0)

	if len(defines) > 0 && len(defines[0]) > 0 {
		entries = defines[0]
	}

	if len(entries) < 1 {
		entries = AppConf.GetMapSlice("logging.loggers")
	}

	_formater := &formatter{}
	_alyslsAppender := &alyslsAppender{}

	for _, entry := range entries {
		name := castx.ToString(entry["name"])

		if name == "" {
			continue
		}

		appenderList := make([]string, 0)

		if a1 := castx.ToStringSlice(entry["appenders"]); len(a1) > 0 {
			appenderList = a1
		} else if s1, ok := entry["appenders"].(string); ok && s1 != "" {
			re1 := regexp.MustCompile(RegexConst.CommaSep)
			appenderList = re1.FindStringSubmatch(s1)
		}

		appenders := make([]appender, 0)

		if len(appenderList) > 0 {
			if slicex.InStringSlice("both", appenderList) || slicex.InStringSlice("file", appenderList) {
				appenders = append(appenders, newFileAppender(map[string]interface{}{
					"channel":   name,
					"filepath":  entry["filepath"],
					"maxSize":   entry["maxSize"],
					"maxBackup": entry["maxBackup"],
				}))
			}

			if slicex.InStringSlice("both", appenderList) || slicex.InStringSlice("alysls", appenderList) {
				appenders = append(appenders, _alyslsAppender)
			}
		}

		minLevel := logrus.DebugLevel

		if lvl, err := logrus.ParseLevel(castx.ToString(entry["level"])); err == nil {
			minLevel = lvl
		}

		_logger := &logrus.Logger{
			Out:       &writer{appenders: appenders},
			Formatter: _formater,
			Level:     minLevel,
		}

		WithLogger(name, &loggerImpl{channel: name, logger: _logger})
	}
}

func WithLogger(name string, logger Logger) {
	loggers[name] = logger
}

func Channel(name string) Logger {
	logger := loggers[name]
	
	if logger == nil {
		logger = NewNoopLogger()
	}
	
	return logger
}

func Log(level interface{}, args ...interface{}) {
	Channel("runtime").Log(level, args...)
}

func Logf(level interface{}, format string, args ...interface{}) {
	Channel("runtime").Logf(level, format, args...)
}

func Trace(args ...interface{}) {
	Channel("runtime").Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	Channel("runtime").Tracef(format, args...)
}

func Debug(args ...interface{}) {
	Channel("runtime").Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Channel("runtime").Debugf(format, args...)
}

func Info(args ...interface{}) {
	Channel("runtime").Info(args...)
}

func Infof(format string, args ...interface{}) {
	Channel("runtime").Infof(format, args...)
}

func Warn(args ...interface{}) {
	Channel("runtime").Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Channel("runtime").Warnf(format, args...)
}

func Error(args ...interface{}) {
	Channel("runtime").Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Channel("runtime").Errorf(format, args...)
}

func Panic(args ...interface{}) {
	Channel("runtime").Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	Channel("runtime").Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	Channel("runtime").Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Channel("runtime").Infof(format, args...)
}
