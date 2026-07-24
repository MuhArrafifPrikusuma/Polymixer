package files_test

import (
	"testing"

	"PolyMixer/files"
)

func TestTakeArg(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		arg *files.Arguments
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files.TakeArg(tt.arg)
		})
	}
}
