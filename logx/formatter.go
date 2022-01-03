package logx

import (
	"github.com/lamgor666/goboot-common/enum/DatetimeFormat"
	"github.com/lamgor666/goboot-common/util/castx"
	"github.com/sirupsen/logrus"
	"strings"
)

type formatter struct {
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	ts := entry.Time.Format(DatetimeFormat.Full)
	channel := castx.ToString(entry.Data["channel"])
	level := strings.ToLower(entry.Level.String())
	msg := strings.TrimSpace(entry.Message)

	contents := strings.Join([]string{
		"ts:" + ts,
		"channel:" + channel,
		"level:" + level,
		"msg:" + msg,
	}, fieldSep)

	return []byte(contents), nil
}
