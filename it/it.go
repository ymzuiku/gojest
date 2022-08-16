package it

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/gojest/stack"
)

var UseFailNow = false

func OnFailTrace(t *testing.T) {
	stack.Debug = true
	for i := 4; i < 99; i++ {
		s := stack.FileLine(i)
		if strings.Contains(s, "runtime.Caller out") {
			break
		}
		t.Logf("Trace %+v", stack.Red(s))
	}
	if UseFailNow {
		t.FailNow()
	}
}

var OnFail = func(t *testing.T) {
	if UseFailNow {
		t.FailNow()
	}
}

func fail(t *testing.T) {
	OnFail(t)
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

func NotContains(t *testing.T, a, contains any) {
	if !assert.NotContains(t, a, contains) {
		fail(t)
	}
}

func Len(t *testing.T, a any, length int) {
	if !assert.Len(t, a, length) {
		fail(t)
	}
}

func Empty(t *testing.T, a any) {
	if !assert.Empty(t, a) {
		fail(t)
	}
}

func NotEmpty(t *testing.T, a any) {
	if !assert.NotEmpty(t, a) {
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

func Less(t *testing.T, e1, e2 any) {
	if !assert.Less(t, e1, e2) {
		fail(t)
	}
}

func LessOrEqual(t *testing.T, e1, e2 any) {
	if !assert.LessOrEqual(t, e1, e2) {
		fail(t)
	}
}

func Greater(t *testing.T, e1, e2 any) {
	if !assert.Greater(t, e1, e2) {
		fail(t)
	}
}

func GreaterOrEqual(t *testing.T, e1, e2 any) {
	if !assert.GreaterOrEqual(t, e1, e2) {
		fail(t)
	}
}
