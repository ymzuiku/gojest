package gojest_test

import (
	"fmt"
	"testing"
)

func TestGojest(t *testing.T) {

	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("error %d", i), func(t *testing.T) {
			t.Fatalf("the %d is error", i)
			t.Fatal("a")
		})
	}
}
