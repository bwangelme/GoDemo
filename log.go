// github.com/op/go-logging 包的测试代码
package main

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

const LogFilename = "/tmp/abc.log"

var log = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type Password string

func (p *Password) Redacted() string {
	return logging.Redact(string(*p))
}

func main() {
	// 接口类型变量的初始值是 nil
	var l logging.LeveledBackend
	fmt.Println(l)

	// %.1s 输出第一位的字符
	fmt.Printf("%.1s\n", logging.CRITICAL)

	// TODO 为什么这两个日志输出前面会有个时间
	// 可能是 log.Flag 的原因
	log.Debugf("%s", "abc")
	log.Infof("%s", "abcd")

	// 创建日志写入的文件
	fd, _ := os.Create(LogFilename)
	defer fd.Close()

	backend1 := logging.NewLogBackend(fd, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)

	log.Debugf("debug %s", Password("secret"))
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
}
