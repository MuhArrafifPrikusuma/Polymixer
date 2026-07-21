package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"PolyMixer/messages"
)

/* NOTE:
objects ID is actually not in order so i can just take the very last object ID then increment by one to create
new object with that id which will be placed in the very first slot after the header and fill it with mp3.
after that i need to recalculate all byte offset by running find_all_obj again and i need to find all new objects
byte offset and after that i will then need to make a function to trackdown each of this object id and which
index they are in before the mp3 append to determine which reference table should be modified later on.
and i will need to rerun that function after append to finally find the new byte offset of all object and replace
the old reference table with new one
*/

func Find_xref(f *os.File) (bs, bsfXref *[]byte, xref_start int) {
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

	from_start_to_xref := buf[:xref_startIdx]
	fromXref_to_end := buf[xref_startIdx:]
	messages.S_file_size("PDF", "without xref", float64(len(from_start_to_xref)))
	return &from_start_to_xref, &fromXref_to_end, xref_startIdx
}

// return the last object index
func Find_all_obj(byteSlice *[]byte, objMap *ObjMap_t) {
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
				fmt.Println("[PROCESS END]All objects found")
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

		// get ID from obj in current scope
		searchZone_ID := fullData[:objIdx]
		lineFeedIdx := bytes.LastIndex(searchZone_ID, []byte("\n"))
		if lineFeedIdx == -1 {
			messages.E_index("line feed")
		}
		messages.S_found_at_index("line feed", lineFeedIdx)

		objID_searchArea := fullData[lineFeedIdx:objIdx]
		idFields := bytes.Fields(objID_searchArea)
		if idFields == nil {
			messages.E_cannot_find_fields(idFields)
		}
		messages.S_found_in_field(idFields)

		ffieldIndex := lineFeedIdx + bytes.Index(objID_searchArea, idFields[0])
		if ffieldIndex == -1 {
			messages.E_index("index of ID")
		}
		messages.S_found_at_index("ID", ffieldIndex)

		id, err := strconv.Atoi(string(idFields[0]))
		if err != nil {
			messages.E_strconv_atoi(err)
		}
		messages.S_found_id(id)

		objMap.objIdx_and_ID[id] = ffieldIndex
		// +6 so that it doesn't find 'obj' end'obj' <- from here
		objMap.endobjId[id] = endObjIdx + 6

		searchStart = endObjIdx + 6
	}
}

type Xref_ObjMap_t struct {
	xref_boffset map[int]*ObjMap_t
}

// NOTE: First i need to create a function to detect all required data (starting point and amount of objects)
// before finally going to iterate through actual ID byte offset tables
// we can do that by looking at the largest index + 1 == amount of objects
func Find_ID_reference(bsfXref *[]byte, objMap *ObjMap_t, bsXref_startp int) {
	fmt.Printf("[PROCESS START]Find ID reference\n")
	fulldata := *bsfXref
	// xref_objmapping := Xref_ObjMap_t{
	// 	xref_boffset: make(map[int]*ObjMap_t),
	// }

	// FIX: save all of this for later use
	target := []byte("\n")
	rltv_firstlf_inTable := bytes.Index(fulldata, target)
	abs_firstlf_inTable := rltv_firstlf_inTable + bsXref_startp
	if rltv_firstlf_inTable == -1 {
		messages.E_index("what file is this bruh")
	}
	messages.S_found_at_index("line feed", abs_firstlf_inTable)

	// use this to find next lf
	tmpScope := fulldata[rltv_firstlf_inTable+1:]
	lastlf_inTable := bytes.Index(tmpScope, target)
	for {
		currentScope := fulldata[rltv_firstlf_inTable:lastlf_inTable]
		table_fields := bytes.Fields(currentScope)
		if table_fields == nil {
			messages.E_cannot_find_fields(table_fields)
		}
		messages.S_found_in_field(table_fields)
	}
}

// NOTE: Save for later when find all object is fixed

func Find_spot_for_new_obj(objMapData *ObjMap_t, file *os.File) int {
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
	messages.S_found_at_index("found spot to append at", appendToIdx)

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
