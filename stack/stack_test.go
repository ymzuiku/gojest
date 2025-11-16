package stack_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/gojest/stack"
)

var err1 = errors.New("hello err")

func useErr1() error {
	return stack.New(err1)
}

func useErr2() error {
	return stack.New(useErr1())
}

func TestStack(t *testing.T) {
	t.Run("", func(t *testing.T) {
		assert.ErrorIs(t, useErr2(), err1)
	})
}

func useErrNil() error {
	return stack.New(nil)
}

func TestStackNil(t *testing.T) {
	t.Run("", func(t *testing.T) {
		assert.ErrorIs(t, useErrNil(), nil)
	})
}
