package main

import (
	"container/ring"
	"fmt"
)

func main() {
	nr := ring.New(1)
	np := nr.Prev()
	np.Value = "233"
	np = np.Prev()
	np = np.Prev()
	np = np.Prev()
	fmt.Println(np.Value)
	np.Unlink(1)
}
