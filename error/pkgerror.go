package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open failed")
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}
	return buf, nil
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, ".settings.xml"))
	return config, errors.WithMessage(err, "counld not read config")
}

func main() {
	_, err := ReadConfig()
	if err != nil {
		// Cause 会拿到最底层的错误
		fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		// %+v 会打印堆栈信息
		fmt.Printf("stack trace:\n%+v\n", err)
		os.Exit(1)
	}
}
