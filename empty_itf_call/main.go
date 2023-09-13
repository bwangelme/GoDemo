package main

import "context"

func main() {
	var cancel context.CancelFunc = func() {}
	cancel()
}
