package watch_test

import (
	"testing"

	"github.com/ymzuiku/gojest/expect"
)

func TestGojestExample1(t *testing.T) {
	t.Run("error 1", func(t *testing.T) {
		expect.Equal(t, 1, 1)
	})

	t.Run("error 2", func(t *testing.T) {
		expect.Equal(t, 2, 2)
	})

	t.Run("error 3", func(t *testing.T) {
		expect.Equal(t, 3, 3)
	})

	t.Run("error 4", func(t *testing.T) {
		expect.Equal(t, 20, 4)
	})

	t.Run("error 5", func(t *testing.T) {
		expect.Equal(t, 20, 5)
	})
}
