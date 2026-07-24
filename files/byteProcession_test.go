package files

import (
	"os"
	"testing"
)

func Test_mp3_get_body(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file *os.File
		want *os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mp3_get_body(tt.file)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("mp3_get_body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPdf_open(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file  *os.File
		want  int
		want2 int
		want3 *os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, got3 := Pdf_open(tt.file)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Pdf_open() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("Pdf_open() = %v, want %v", got2, tt.want2)
			}
			if true {
				t.Errorf("Pdf_open() = %v, want %v", got3, tt.want3)
			}
		})
	}
}
