package unittest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit2(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		// TODO: Add test cases.
		{"base case", args{"a:b:c", ":"}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split2(tt.args.s, tt.args.sep)
			assert.Equal(t, tt.wantResult, got)
			// if gotResult := Split2(tt.args.s, tt.args.sep); !reflect.DeepEqual(gotResult, tt.wantResult) {
			// 	t.Errorf("Split2() = %v, want %v", gotResult, tt.wantResult)
			// }
		})
	}
}
