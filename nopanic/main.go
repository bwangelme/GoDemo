package main

import (
	"errors"
	"fmt"
)

var (
	RetryEmailTaskErr *retryEmailTaskErr
)

type retryEmailTaskErr struct {
	msg string
}

func (e *retryEmailTaskErr) New(msg string) *retryEmailTaskErr {
	return &retryEmailTaskErr{
		msg: msg,
	}
}

func (e *retryEmailTaskErr) Is(err error) bool {
	var t = &retryEmailTaskErr{}
	return errors.Is(err, t)
}

func (e *retryEmailTaskErr) Error() string {
	return e.msg
}

// golang 中，无限递归导致的程序 OOM，panic 是捕获不到的
// 我们的监控程序，只能看到程序重启。但是具体重启原因，只能到 stderr 中查看进程的栈信息
func main() {
	defer func() {
		if res := recover(); res != nil {
			fmt.Println("PPPPanic, res", res)
		}
	}()

	t := RetryEmailTaskErr.New("abc")
	RetryEmailTaskErr.Is(t)
}
