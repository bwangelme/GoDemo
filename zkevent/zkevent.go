package main

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

func main() {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, 10*time.Minute)
	if err != nil {
		log.Fatalln("connect zk failed")
	}

	path := "/xyd/test_sub"

	for {
		log.Printf("Start to watch %s\n", path)

		sub, _, ch, err := conn.ChildrenW(path)
		if err != nil {
			log.Fatalf("watch %s failed %v\n", path, err)
		}
		log.Println("Sub is", sub)

		select {
		case e := <-ch:
			log.Printf("Receive event %s\n", e)
		}
	}
}
