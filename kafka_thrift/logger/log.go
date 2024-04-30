package logger

import (
	"github.com/sirupsen/logrus"
)

var (
	L *logrus.Logger
)

func init() {
	L = logrus.New()
	L.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-02/15:04:05",
		FullTimestamp:   true,
	})
}
