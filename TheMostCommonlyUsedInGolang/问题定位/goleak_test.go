package goleak1

import (
	"go.uber.org/goleak"
	"testing"
)

func Test_leak(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"base case"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			leak()
		})
	}
}
