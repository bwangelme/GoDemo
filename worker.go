package main

// 这段代码展示了 waitGroup 的典型用法
// 子进程主动结束，主进程被动等待

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func worker(wg *sync.WaitGroup, dir string, file os.FileInfo) {
	var lines int
	absName := filepath.Join(dir, "", file.Name())
	fd, err := os.Open(absName)
	if err != nil {
		wg.Done()
		log.Println(err)
		return
	}

	reader := bufio.NewReader(fd)
	for {
		_, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		lines += 1
	}

	fmt.Printf("%-20s\t%-10d\t%-10d\n", file.Name(), file.Size(), lines)
	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	dir := "ml-latest-small"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	count := len(files)
	wg.Add(count)

	fmt.Printf("%-20s\t%-10s\t%-10s\n", "文件名", "文件大小", "文件行数")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		go worker(&wg, dir, file)
	}

	wg.Wait()
}
