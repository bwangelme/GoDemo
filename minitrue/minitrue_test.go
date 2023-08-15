package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCond(t *testing.T) {
	assert.Equal(t, 1, Cond(true, 1, 2))
	assert.Equal(t, 2, Cond(false, 1, 2))
}

func TestOr(t *testing.T) {
	// 范型函数可以不传类型参数
	// 我在 1.18, 1.19, 1.20 上测试均可以工作
	assert.Equal(t, 0, Or[int]())
	assert.Equal(t, 0, Or(0))
	assert.Equal(t, 1, Or(1))
	assert.Equal(t, 2, Or(0, 2))
	assert.Equal(t, 3, Or(0, 0, 3))
}
