package logger

import (
	"github.com/sirupsen/logrus"
)

var (
	L *logrus.Logger
)

func init() {
	L = logrus.New()
	L.SetFormatter(&logrus.TextFormatter{})
}
