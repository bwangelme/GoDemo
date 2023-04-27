package main

/*
 本例子展示了 logrus 如何添加预设字段
*/

import "github.com/sirupsen/logrus"

var (
	L *logrus.Logger
)

type customLogger struct {
	formatter logrus.Formatter
	domain    string
}

func (l *customLogger) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["domain"] = l.domain
	return l.formatter.Format(entry)
}

func main() {
	L = logrus.New()
	L.SetLevel(logrus.InfoLevel)
	L.SetFormatter(&customLogger{
		formatter: &logrus.JSONFormatter{},
		domain:    "custom-domain",
	})

	L.Info("this msg")
}
