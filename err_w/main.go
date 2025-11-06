package main

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	OriginErr = errors.New("original error")
)

// 底层的 error，上层使用 errors.Is 直接判断就可以
func main() {
	err := topFunction()
	if errors.Is(err, OriginErr) {
		fmt.Printf("%+v\n", err)
	}
}

func topFunction() error {
	err := middleFunction()
	if err != nil {
		return fmt.Errorf("efg %w", err)
	}
	return nil
}

func middleFunction() error {
	err := deeperFunction()
	if err != nil {
		return fmt.Errorf("abc %w", err)
	}
	return nil
}

func deeperFunction() error {
	return OriginErr
}
