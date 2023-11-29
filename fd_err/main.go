package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9125}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	statsdData := []byte("dae.nezha.reqs:1|c|@0.25*4\n")
	conn.Write(statsdData)
	conn.Close()
	fmt.Printf("<%s>\n", conn.RemoteAddr())
}
