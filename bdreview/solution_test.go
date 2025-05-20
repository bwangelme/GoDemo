package bdreview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_solution(t *testing.T) {
	res := solution([]int{1, 2, 3, 4, 5, 6}, 12)
	assert.Equal(t, res, [][]int{{3, 4, 5}, {2, 4, 6}, {1, 5, 6}})
}
