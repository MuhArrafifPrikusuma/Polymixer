package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"PolyMixer/messages"
)

// NOTE: if fail this is no longer needed
func Find_xref(f *os.File) (bs *[]byte) {
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
	return &byteSlice
}

// return the last object index
// NOTE: if fail might have to change it to find every single object index
func Find_all_obj(byteSlice *[]byte, objMap *ObjMap) {
	fullData := *byteSlice
	searchStart := 0

	for {
		if searchStart >= len(fullData) {
			break
		}
		currentZone := fullData[searchStart:]
		fmt.Println("search start from", searchStart)

		// get Object starting Index
		relative_ObjIdx := bytes.Index(currentZone, []byte("obj"))
		if relative_ObjIdx == -1 {
			if searchStart != 0 {
				fmt.Println("[TEMPORARY] either all objects found or corrupted file or just my stupid code breaks")
				break
			}
			messages.E_index("obj")
		}
		objIdx := searchStart + relative_ObjIdx
		messages.S_found_at_index("obj", objIdx)

		// get endobj starting index
		relative_EndObjIdx := bytes.Index(currentZone[relative_ObjIdx:], []byte("endobj"))
		if relative_EndObjIdx == -1 {
			messages.E_index("endobj")
		}
		endObjIdx := searchStart + relative_EndObjIdx + relative_ObjIdx
		messages.S_found_at_index("endobj", endObjIdx)

		// get ID from current obj in scope
		searchZone_ID := fullData[:objIdx]
		lineFeedIdx := bytes.LastIndex(searchZone_ID, []byte("\n"))
		if lineFeedIdx == -1 {
			messages.E_index("line feed")
		}
		messages.S_found_at_index("line feed", lineFeedIdx)

		objID_searchArea := fullData[lineFeedIdx:objIdx]
		idFields := bytes.Fields(objID_searchArea)
		if idFields == nil {
			fmt.Println("[ERROR] TEMPORARY PLACEHOLDER")
		}
		messages.S_found_in_field(idFields)

		id, err := strconv.Atoi(string(idFields[0]))
		if err != nil {
			messages.E_strconv_atoi(err)
		}
		messages.S_found_id(id)

		objMap.obj_and_id[objIdx] = id
		// +6 so that it doesn't find the end'obj' <- from here
		objMap.endobjId[id] = endObjIdx + 6

		searchStart = endObjIdx + 6
	}
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

// NOTE: Save for later when find all object is fixed

func Find_spot_for_new_obj(objMapData *ObjMap, file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	buf := make([]byte, fileStat.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	// temporary data
	findLastLineFeed := buf[0:]
	// this is the actual one but need change -> buf[objidx:xrefstrtidx]
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
