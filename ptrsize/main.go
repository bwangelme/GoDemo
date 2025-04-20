package main

import (
	"fmt"
	"unsafe"
)

const PtrSize = 4 << (^uintptr(0) >> 63)

func main() {
	fmt.Printf("PtrSize = %d\n", PtrSize) // 输出 4 或 8
	fmt.Printf("unsafe.Sizeof(uintptr(0)) = %d\n", unsafe.Sizeof(uintptr(0)))
}
