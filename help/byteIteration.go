package help

import (
	"PolyMixer/messages"
	"bytes"
	"os"
)

// "PolyMixer/messages"
// "bytes"
// "os"

// TODO: use this function later to cut the head of and metadata of pdf file and then also check the unused space between those two and use it to embed the whole mp3 binary there
// func search_and_count(lookFor, startFrom string, ignoreStartingValue bool, file *os.File) (offset int) {
// 	fileStat, err := file.Stat()
// 	if err != nil {
// 		messages.E_stat_read(err)
// 	}
//
// 	buf_pdf := make([]byte, fileStat.Size())
// 	_, err = file.ReadAt(buf_pdf, 0)
// 	if err != nil {
// 		messages.E_read(err)
// 	}
//
// 	startingIndex := bytes.Index(buf_pdf, []byte(startFrom))
// 	if startingIndex == -1 {
// 		messages.E_index("startingIndex")
// 	}
//
// 	sliced_buf_pdf := buf_pdf
// 	if ignoreStartingValue {
// 		sliced_buf_pdf = buf_pdf[startingIndex:]
// 	}
// 	endIndex := bytes.Index(sliced_buf_pdf, []byte(lookFor))
// 	if endIndex == -1 {
// 		messages.E_index("endIndex")
// 	}
//
// 	// TODO: remind me to count offset and empty spaces between startIndex and endIndex later on to be used as
// 	offset = 0
// 	for i := startingIndex; i < endIndex; offset++ {
// 	}
//
// 	return
// }

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
