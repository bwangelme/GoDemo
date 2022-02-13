package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.Infof("YES")
}
