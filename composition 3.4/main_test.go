package main1

import (
	"testing"

	"go.uber.org/goleak"
)

func TestGenerate(t *testing.T) {
	defer goleak.VerifyNone(t)

	// cancel := make(chan struct{})
	// defer close(cancel)
	// c1 := generate(cancel)
	// c2 := takeUnique(cancel, c1)
	// c3_1 := reverse(cancel, c2)
	// c3_2 := reverse(cancel, c2)
	// c4 := merge(cancel, c3_1, c3_2)
	// print(cancel, c4, 10)

	for i := 0; i < 10; i++ {
		t.Run("test.name", func(t *testing.T) {
			cancel := make(chan struct{})
			defer close(cancel)
			c1 := generate(cancel)
			c2 := takeUnique(cancel, c1)
			c3_1 := reverse(cancel, c2)
			c3_2 := reverse(cancel, c2)
			c4 := merge(cancel, c3_1, c3_2)
			print(cancel, c4, 10)
		})
	}
}
