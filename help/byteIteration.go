package help

import (
	"bytes"
	"io"
	"os"
	"strconv"

	"PolyMixer/messages"
)

// NOTE: if fail this is no longer needed
func Find_xref(f *os.File) (xrefStrtIdx int, bs *[]byte) {
	fileStat, err := f.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	messages.S_file_size("PDF", "full", float64(fileStat.Size()))
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
	messages.S_file_size("PDF", "without xref", float64(len(byteSlice)))
	return xref_startIdx, &byteSlice
}

type ObjMap struct {
	// contain [index]id
	obj_and_id map[int]int
	// contain [id]index
	endobjId map[int]int
}

// return the last object index
// NOTE: if fail might have to change it to find every single object index
func Find_all_obj(xref_start_idx int, byteSlice *[]byte) (endObjIdx, objIdx int) {
	searchEnd := xref_start_idx
	if searchEnd > len(*byteSlice) {
		searchEnd = len(*byteSlice)
	}

	searchZone := (*byteSlice)[0:searchEnd]
	for {

		endObjIdx = bytes.Index(searchZone, []byte("endobj"))
		endObjIdx += 7
		if endObjIdx == -1 {
			messages.E_index("obj")
		}
		messages.S_found_at_index("endobj", endObjIdx)

		objIdx = bytes.Index(searchZone, []byte("obj"))
		if objIdx == -1 {
			messages.E_index("obj")
		}
		messages.S_found_at_index("obj", objIdx)

		// find the id and id index
		line_feedIdx := bytes.Index(searchZone, []byte("\n"))
		white_spaceIdx := bytes.Index(searchZone, []byte(" "))
		if line_feedIdx == -1 || white_spaceIdx == -1 {
			if line_feedIdx == -1 {
				messages.E_index("line feed")
			}
			messages.E_index("white space")
		}
		messages.S_found_at_index("line feed", line_feedIdx)
		messages.S_found_at_index("white space", line_feedIdx)

		// FIXME: this currently read the first newline not the one right before
		// the object id so find a way to fix that
		find_full_obj_id_string := (*byteSlice)[line_feedIdx:white_spaceIdx]
		objIdStr := string(bytes.TrimSpace(find_full_obj_id_string))
		objIdIdx := int(line_feedIdx + 1)
		id, err := strconv.Atoi(objIdStr)
		if err != nil {
			messages.E_strconv_atoi(err)
		}
		messages.S_found_id(id)

		objMap := ObjMap{
			obj_and_id: make(map[int]int),
			endobjId:   make(map[int]int),
		}
		objMap.obj_and_id[objIdIdx] = id
		objMap.endobjId[id] = endObjIdx

		searchZone = searchZone[:xref_start_idx]
	}
	return endObjIdx, objIdx
}

// NOTE: if fail might have to change it to find every single id
//  func Find_obj_id(objIdx int, byteSlice *[]byte) int {
//  	searchZone := (*byteSlice)[0:objIdx]
//
//  	var line_feed_AfterId, line_feed_BeforeId int
//  	for i := range 3 {
//  		whiteSpaceBeforeObjIdx := bytes.Index(searchZone, []byte(" "))
//  		if whiteSpaceBeforeObjIdx == -1 {
//  			messages.E_index("white space")
//  		}
//  		messages.S_found_at_index("white space", whiteSpaceBeforeObjIdx)
//  		if i < 2 {
//  			line_feed_AfterId = whiteSpaceBeforeObjIdx
//  			searchZone = searchZone[0 : line_feed_AfterId-(i+1)]
//  		}
//  		line_feed_BeforeId = whiteSpaceBeforeObjIdx
//  	}
//
//  	searchZone = searchZone[line_feed_BeforeId:line_feed_AfterId]
//  	objIdStr := string(bytes.TrimSpace(searchZone))
//
//  	id, err := strconv.Atoi(objIdStr)
//  	if err != nil {
//  		messages.E_strconv_atoi(err)
//  	}
//  	messages.S_found_id(id)
//
//  	return id
//  }

// NOTE: if fails this will also need to change ofcourse
func Find_spot_for_new_obj(objIdx, xrefStrtIdx int, file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	buf := make([]byte, fileStat.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	findLastLineFeed := buf[objIdx:xrefStrtIdx]
	appendToIdx := bytes.Index(findLastLineFeed, []byte("\n")) + 1
	if appendToIdx == -1 {
		messages.E_index("line feed")
	}
	messages.S_found_at_index("line feed", appendToIdx)

	return appendToIdx
}

// NOTE : create mp3 object to mix
func create_mp3_obj(appendToIdx, lastObjId int, fileMp3 *os.File) {
}

func Mix_MP3_and_PDF(filePdf, fileMp3 *os.File, appendToIdx, lastObjId int) {
	//  fmt.Printf("[PROCESS] Mixing files\n")
	//  fileStatPdf, err := filePdf.Stat()
	//  if err != nil {
	//  	messages.E_stat_read(err)
	//  }
	//  fileStatMp3, err := fileMp3.Stat()
	//  if err != nil {
	//  	messages.E_stat_read(err)
	//  }
	//  // create buffer for newfile
	//  buf := make([]byte, fileStatPdf.Size()+fileStatMp3.Size())
	//  bufPdf := make([]byte, fileStatPdf.Size())
	//  bufMp3 := make([]byte, fileStatMp3.Size())
	//  _, err = filePdf.ReadAt(bufPdf, 0)
	//  if err != nil {
	//  	messages.E_read(err)
	//  }
	//  // mp3 goes after this
	//  pdfFileWindow := bufPdf[0:appendToIdx]
	//  create_mp3_obj(appendToIdx, lastObjId, fileMp3)
}
