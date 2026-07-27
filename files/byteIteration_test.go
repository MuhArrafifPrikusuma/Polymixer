package files_test

import (
	"os"
	"testing"

	"PolyMixer/files"
)

func TestFind_xref(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		f     *os.File
		want  *[]byte
		want2 *[]byte
		want3 int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, got3 := files.Find_xref(tt.f)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Find_xref() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("Find_xref() = %v, want %v", got2, tt.want2)
			}
			if true {
				t.Errorf("Find_xref() = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestFind_all_obj(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		byteSlice *[]byte
		objMap    *files.ObjMap_t
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files.Find_all_obj(tt.byteSlice, tt.objMap)
		})
	}
}

func TestFind_ID_reference(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		bsfXref       *[]byte
		objMap        *files.ObjMap_t
		bsXref_startp int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files.Find_ID_reference(tt.bsfXref, tt.objMap, tt.bsXref_startp)
		})
	}
}

func TestFind_spot_for_new_obj(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		objMapData *files.ObjMap_t
		file       *os.File
		want       int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := files.Find_spot_for_new_obj(tt.objMapData, tt.file)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Find_spot_for_new_obj() = %v, want %v", got, tt.want)
			}
		})
	}
}
