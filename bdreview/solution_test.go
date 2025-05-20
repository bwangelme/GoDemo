package bdreview

import (
	"fmt"
	"testing"
)

func Test_solution(t *testing.T) {
	res := solution([]int{1, 2, 3, 4, 5, 6}, 12)
	fmt.Println("res", res)
	t.Fail()
}
