package help

import (
	"PolyMixer/messages"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

// NOTE: if fail this is no longer needed
func Find_xref(f *os.File) (xrefStrtIdx int, bs *[]byte) {
	fileStat, err := f.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	messages.S_file_size("PDF", "full", float64(fileStat.Size()), messages.MB)
	buf := make([]byte, fileStat.Size())

	_, err = f.ReadAt(buf, 0)
	if err != nil && err != io.EOF {
		messages.E_read(err)
	}

	target := []byte("xref")
	xref_startIdx := bytes.Index(buf, target)
	if xref_startIdx == -1 {
		messages.E_byte_slice_too_small(xref_startIdx)
	}
	messages.S_found_at_index(string(target), xref_startIdx)
	byteSlice := buf[0:xref_startIdx]
	messages.S_file_size("PDF", "without xref", float64(len(byteSlice)), messages.MB)
	return xref_startIdx, &byteSlice
}

// return the last object index
// NOTE: if fail might have to change it to find every single object index
func Find_last_obj_idx(xref_start_idx int, byteSlice *[]byte) (lastEndObjIdx, lastObjIdx int) {

	searchEnd := xref_start_idx
	if searchEnd > len(*byteSlice) {
		searchEnd = len(*byteSlice)
	}

	searchZone := (*byteSlice)[0:searchEnd]

	lastEndObjIdx = bytes.LastIndex(searchZone, []byte("endobj"))
	if lastEndObjIdx == -1 {
		messages.E_index("last endobj")
	}
	messages.S_found_at_index("endobj", lastEndObjIdx)
	searchZone = searchZone[0:lastEndObjIdx]
	lastObjIdx = bytes.LastIndex(searchZone, []byte("obj"))
	if lastObjIdx == -1 {
		messages.E_index("last object")
	}
	messages.S_found_at_index(" obj", lastObjIdx)
	lastObjIdx = lastObjIdx + 1
	return lastEndObjIdx, lastObjIdx
}

// NOTE: if fail might have to change it to find every single id
func Find_obj_id(lastObjIdx int, byteSlice *[]byte) int {
	searchZone := (*byteSlice)[0:lastObjIdx]

	var line_feed_AfterId, line_feed_BeforeId int
	lineFeedAsStartIndex := bytes.LastIndex(searchZone, []byte("\n"))
	messages.S_found_at_index("line feed", lineFeedAsStartIndex)
	for i := range 3 {
		whiteSpaceBeforeObjIdx := bytes.LastIndex(searchZone, []byte(" "))
		if whiteSpaceBeforeObjIdx == -1 {
			messages.E_index("white space")
		} else if lineFeedAsStartIndex == -1 {
			messages.E_index("line feed")
		}
		messages.S_found_at_index("white space", whiteSpaceBeforeObjIdx)
		if i < 2 {
			line_feed_AfterId = whiteSpaceBeforeObjIdx
			searchZone = searchZone[0 : line_feed_AfterId-(i+1)]
		}
		line_feed_BeforeId = lineFeedAsStartIndex
	}

	searchZone = searchZone[line_feed_BeforeId:line_feed_AfterId]
	objIdStr := string(bytes.TrimSpace(searchZone))

	id, err := strconv.Atoi(objIdStr)
	if err != nil {
		messages.E_strconv_atoi(err)
	}
	messages.S_found_id(id)

	return id
}

// NOTE: if fails this will also need to change ofcourse
func Find_spot_for_new_obj(lastObjIdx, xrefStrtIdx int, file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	buf := make([]byte, fileStat.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	findLastLineFeed := buf[lastObjIdx:xrefStrtIdx]
	appendToIdx := bytes.LastIndex(findLastLineFeed, []byte("\n")) + 1
	if appendToIdx == -1 {
		messages.E_index("line feed")
	}
	messages.S_found_at_index("line feed", appendToIdx)

	return appendToIdx
}

func Mix_MP3_and_PDF(filePdf, fileMp3 *os.File, appendToIdx, lastObjId int) (buf []byte) {
	fmt.Printf("[PROCESS] Mixing files\n")
	fileStat, err := filePdf.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	buf = make([]byte, fileStat.Size())
	return buf
}
