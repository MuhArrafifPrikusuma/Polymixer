package help

import (
	"PolyMixer/messages"
	"bytes"
	"log"
	"math"
	"os"
	"strconv"
)

func find_xref(f *os.File) (bs *[]byte) {
	target := []byte("xref")
	fileStat, err := f.Stat()
	defer f.Close()
	if err != nil {
		messages.E_stat_read(err)
	}

	buf := make([]byte, fileStat.Size())

	xref_startIdx := bytes.Index(buf, target)
	byteSlice := buf[:xref_startIdx]

	find_last_obj_idx(xref_startIdx, &byteSlice)
	return &byteSlice
}

// return the last object index
func find_last_obj_idx(xref_start_idx int, byteSlice *[]byte) (lastEndObjIdx, lastObjIdx int) {
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

func find_obj_id(lastObjIdx int, byteSlice *[]byte) int {
	searchZone := (*byteSlice)[0:lastObjIdx]
	data := *byteSlice

	var line_feed_AfterId, line_feed_BeforeId int
	for i := range 2 {
		whiteSpaceBeforeObjIdx := bytes.LastIndex(searchZone, []byte("\n"))
		if i < 1 {
			line_feed_AfterId = whiteSpaceBeforeObjIdx
			searchZone = searchZone[0 : line_feed_AfterId-1]
		} else {
			line_feed_BeforeId = whiteSpaceBeforeObjIdx
		}
	}

	searchZone = data[line_feed_BeforeId:line_feed_AfterId]
	searchZone = bytes.TrimSpace(searchZone)
	ids := bytes.Fields(searchZone)

	objId := ids[0]

	id, err := strconv.Atoi(string(objId))
	if err != nil {
		log.Fatal("placeholder")
	}

	return id
}

func find_places_for_new_bin(lastEndObjIdx int) {

}
