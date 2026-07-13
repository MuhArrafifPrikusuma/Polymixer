package help

import (
	"PolyMixer/messages"
	"bytes"
	"os"
)

func Find_xref(f *os.File) {
	target := []byte("xref")
	fileStat, err := f.Stat()
	defer f.Close()
	if err != nil {
		messages.E_stat_read(err)
	}

	buf := make([]byte, fileStat.Size())

	xref_startIdx := bytes.Index(buf, target)
	byteSlice := buf[:xref_startIdx]

	Find_last_obj_idx(xref_startIdx, &byteSlice)
}

// return the last object index
func Find_last_obj_idx(xref_start_idx int, byteSlice *[]byte) (lastEndObjIdx, lastObjIdx int) {
	if byteSlice == nil || len(*byteSlice) == 0 {
		messages.E_byte_slice_too_small()
	}

	searchEnd := xref_start_idx
	if searchEnd > len(*byteSlice) {
		searchEnd = len(*byteSlice)
	}

	searchZone := (*byteSlice)[0:searchEnd]

	lastEndObjIdx = bytes.LastIndex(searchZone, []byte("endobj"))
	if lastEndObjIdx == -1 {
		messages.E_index("last endobj")
	}
	searchZone = searchZone[0:lastEndObjIdx]
	lastObjIdx = bytes.LastIndex(searchZone, []byte(" obj"))
	if lastObjIdx == -1 {
		messages.E_index("last object")
	}
	lastObjIdx = lastObjIdx + 1
	return lastEndObjIdx, lastObjIdx
}

// TODO: find last xref then iterate back to find last endobj then find the last obj id then append
// to last_endobj_idx + 1 start with 1010 and end with 1010 places all of the data inside obj stream and
// don't reference it on xref
func Find_places_for_newObj(last_endobj_idx, last_obj_id int) {

}
