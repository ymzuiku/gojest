package watch_test

import (
	"fmt"
	"testing"

	"github.com/ymzuiku/gojest/it"
)

func TestGojest(t *testing.T) {

	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("error %d", i), func(t *testing.T) {
			it.Equal(t, 20, 30)
		})
	}
}
