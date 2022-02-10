package main

import (
	"io"
	"log"
	"net"
	"syscall"
)

func main() {
	const addr = "/tmp/uds.sock"
	syscall.Unlink(addr)

	l, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(c)
	}
}
