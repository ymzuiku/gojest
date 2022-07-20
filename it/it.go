package it

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/gojest/stack"
)

/*
更便于开发断言

1. 约定不要有多余的message参数
2. 会使用stack包裹, 高亮显示, 方便定位
*/

var fails = 0
var UseFailNow = false

func fail(t *testing.T) {
	stack.Debug = true
	fails += 1
	t.Logf("Faile:%d %s", fails, stack.Red(stack.FileLine(3)))
	if UseFailNow {
		t.FailNow()
	}
}

func Equal(t *testing.T, a, b any) {
	if !assert.Equal(t, a, b) {
		fail(t)
	}
}

func NotEqual(t *testing.T, a, b any) {
	if !assert.NotEqual(t, a, b) {
		fail(t)
	}
}

func True(t *testing.T, a bool) {
	if !assert.True(t, a) {
		fail(t)
	}
}
func False(t *testing.T, a bool) {
	if !assert.False(t, a) {
		fail(t)
	}
}

func ErrorIs(t *testing.T, a, b error) {
	if !assert.ErrorIs(t, a, b) {
		fail(t)
	}
}

func Nil(t *testing.T, a any) {
	if !assert.Nil(t, a) {
		fail(t)
	}
}

func NotNil(t *testing.T, a any) {
	if !assert.NotNil(t, a) {
		fail(t)
	}
}

func Contains(t *testing.T, a, contains any) {
	if !assert.Contains(t, a, contains) {
		fail(t)
	}
}

func Empty(t *testing.T, a any) {
	if !assert.Empty(t, a) {
		fail(t)
	}
}

func EqualError(t *testing.T, a error, errString string) {
	if !assert.EqualError(t, a, errString) {
		fail(t)
	}
}

func ElementsMatch(t *testing.T, a, b any) {
	if !assert.ElementsMatch(t, a, b) {
		fail(t)
	}
}

func Zero(t *testing.T, a any) {
	if !assert.Zero(t, a) {
		fail(t)
	}
}

func NotZero(t *testing.T, a any) {
	if !assert.NotZero(t, a) {
		fail(t)
	}
}

func Subset(t *testing.T, list any, subset any) {
	if !assert.Subset(t, list, subset) {
		fail(t)
	}
}

func DirExists(t *testing.T, path string) {
	if !assert.DirExists(t, path) {
		fail(t)
	}
}

func ErrorContains(t *testing.T, a error, contains string) {
	if !assert.ErrorContains(t, a, contains) {
		fail(t)
	}
}

func ErrorAs(t *testing.T, a error, b any) {
	if !assert.ErrorAs(t, a, b) {
		fail(t)
	}
}

func FailNow(t *testing.T, failureMessage string) {
	if !assert.FailNow(t, failureMessage) {
		fail(t)
	}
}

func Fail(t *testing.T, failureMessage string) {
	if !assert.Fail(t, failureMessage) {
		fail(t)
	}
}
func FileExists(t *testing.T, path string) {
	if !assert.FileExists(t, path) {
		fail(t)
	}
}
