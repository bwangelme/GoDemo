package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// 使用 pkg/errors 追中 error 的栈信息
func main() {
	err := topFunction()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func topFunction() error {
	err := middleFunction()
	if err != nil {
		return errors.WithMessage(err, "error in topFunction")
	}
	return nil
}

func middleFunction() error {
	err := deeperFunction()
	if err != nil {
		return errors.WithMessage(err, "error in middleFunction")
	}
	return nil
}

func deeperFunction() error {
	return errors.New("original error")
}
