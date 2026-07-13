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

	Find_last_obj_idx_and_last_obj_id(xref_startIdx, &byteSlice)
}

func Find_last_obj_idx_and_last_obj_id(xref_start_idx int, byteSlice *[]byte) {
	// [int] is for the count of how many obj been found, and int for that idx
	objFound := make(map[int]int, 0)
	var countFoundObj int
	target := []byte("obj")
	sliceLen := len(*byteSlice)
	n := 100
	sliceStart := sliceLen - n
	currentSlice := (*byteSlice)[sliceStart:sliceLen]
	// NOTE: i think i won't use this but im not sure yet
	//	currentLen := len(currentSlice)

	// TODO: find a way to search for target and make sure that there is only 2 obj at max with one as the starting obj and other endobj right before xref
	for {
		objIdx := bytes.Index(currentSlice, target)
		if objIdx != 0 {
			countFoundObj++
			objFound[countFoundObj] = objIdx

		}
		if countFoundObj > 2 {
			currentSlice = currentSlice[objFound[countFoundObj]+1 : sliceLen]

			for i := range countFoundObj {
				objFound[i-1] = objFound[countFoundObj]
			}
		}
	}
}

// TODO: find last xref then iterate back to find last endobj then find the last obj id then append
// to last_endobj_idx + 1 start with 1010 and end with 1010 places all of the data inside obj stream and
// don't reference it on xref
func Find_places_for_newObj(last_endobj_idx, last_obj_id int) {

}
