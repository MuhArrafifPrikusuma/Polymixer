package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"unsafe"

	"PolyMixer/messages"
)

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

		objMap.objIdx_and_ID[ffieldIndex] = id
		// +6 so that it doesn't find 'obj' end'obj' <- from here
		objMap.endobjId[id] = endObjIdx + 6

		searchStart = endObjIdx + 6
	}
}

// map this (self explanatory)
type Objref_Table_t struct {
	objbyte_offset []byte
	refbyte_offset int
	genNumber      []byte
	marker         []byte
}
type Xref_ObjMap_t struct {
	xref_BObjoffset map[int]Objref_Table_t
}

// consume byte slice slice -> starting point of the reference table and
// number of index, last line feed (relative index) <- add with bsXref_startingpoint to turn into abs
func read_xref_data(bsfXref *[]byte) (startP, numOIdx, lastlf int) {
	fulldata := *bsfXref
	target := []byte("\n")
	firstlf := bytes.Index(fulldata, target)
	nextlf := fulldata[firstlf+1:]
	lastlf = bytes.Index(nextlf, target)
	if firstlf == -1 || lastlf == -1 {
		messages.E_index("line feed")
	}
	fmt.Println("[PROCESS]Looking objects info")
	lastlf += firstlf + 1

	lookFor_Fields := fulldata[firstlf:lastlf]
	fields := bytes.Fields(lookFor_Fields)
	if fields == nil {
		messages.E_cannot_find_fields(fields)
	}

	sprt_field := make([]int, 2)
	for i := range fields {
		split_field, err := strconv.Atoi(string(fields[i]))
		if err != nil {
			messages.E_strconv_atoi(err)
		}
		sprt_field[i] = split_field
	}
	return sprt_field[0], sprt_field[1], lastlf
}

func Find_ID_reference(bsfXref *[]byte, objMap *ObjMap_t, bsXref_startp int) (ptrXrefDat *Xref_ObjMap_t, firstFoundIDoffset int) {
	refID := 0

	fulldata := *bsfXref
	fmt.Printf("[PROCESS START]Find ID reference\n")

	refStart, numsO, startP := read_xref_data(bsfXref)
	messages.S_found_xref_data(refStart, numsO, startP, bsXref_startp)
	objRefTable := &Objref_Table_t{
		objbyte_offset: []byte(""),
		refbyte_offset: -1,
		genNumber:      []byte(""),
		marker:         []byte(""),
	}
	XrefMapping := &Xref_ObjMap_t{
		xref_BObjoffset: make(map[int]Objref_Table_t),
	}

	for {
		prepareField := fulldata[startP+1:]
		target := []byte("\n")
		nextlfIndex := bytes.Index(prepareField, target) + startP
		makeField := fulldata[startP:nextlfIndex]

		table_fields := bytes.Fields(makeField)
		if len(table_fields) <= 1 {
			break
		}
		if table_fields == nil {
			messages.E_cannot_find_fields(table_fields)
		}

		basePtr := uintptr(unsafe.Pointer(&fulldata[0]))
		var byteIndex uintptr
		var field []byte
		objbyte_off_int, objgenNumber_int := 0, 0

		for _, field = range table_fields {
			var err error
			fieldPtr := uintptr(unsafe.Pointer(&field[0]))

			byteIndex = fieldPtr - basePtr
			objbyte_off_int, err = strconv.Atoi(string(table_fields[0]))
			objgenNumber_int, err = strconv.Atoi(string(table_fields[1]))
			objRefTable.objbyte_offset = table_fields[0]
			objRefTable.refbyte_offset = int(byteIndex) + bsXref_startp
			objRefTable.genNumber = table_fields[1]
			objRefTable.marker = table_fields[2] // <- this is decimal
			if err != nil {
				messages.E_strconv_atoi(err)
			}
		}
		id := objMap.objIdx_and_ID[objbyte_off_int]
		if objbyte_off_int < firstFoundIDoffset || firstFoundIDoffset == 0 {
			firstFoundIDoffset = objbyte_off_int
		}
		XrefMapping.xref_BObjoffset[id] = *objRefTable
		fmt.Printf("[ALLOCATE]sizeof %vB for -> objID %v\n", unsafe.Sizeof(*objRefTable), id)
		fmt.Printf("[STORE]ref offset: %v\n[STORE]byte offset: %v\n[STORE]genNumber: %v\n[STORE]marker: %v\n",
			objRefTable.refbyte_offset, objbyte_off_int, objgenNumber_int, objRefTable.marker)
		// id := objMap.objIdx_and_ID[]

		startP = nextlfIndex + 1
		refID++
	}
	fmt.Println("[PROCESS END]Found! and store all value")
	return XrefMapping, firstFoundIDoffset
}

// CUT the head to then append the mp3 objstream right after head
func Cut_HEAD_to(objMapData *ObjMap_t, file *os.File, firstFoundID int) int {
	fileStat, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	buf := make([]byte, fileStat.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}

	findLastLineFeed := buf[:firstFoundID]
	fmt.Println(firstFoundID)

	// this will go to the previously filled position by id X
	cutTo := bytes.LastIndex(findLastLineFeed, []byte("\n"))
	if cutTo == -1 {
		messages.E_index("line feed")
	}
	messages.S_found_at_index("spot to append at", cutTo)

	return cutTo
}

// take replace startxref with []new data
func StartXref_refOffset(bsfXref *[]byte, added int) *[]byte {
	fulldata := *bsfXref

	fstartxref := bytes.Index(fulldata, []byte("startxref"))
	flf := bytes.IndexByte(fulldata[fstartxref:], '\n')

	flf += fstartxref + 1
	fnlf := bytes.IndexByte(fulldata[flf:], '\n')
	if flf == -1 || fnlf == -1 || fstartxref == -1 {
		messages.E_index("startxref")
	} else {
		fnlf += flf
	}
	startxref_valF := bytes.Fields(fulldata[flf:fnlf])

	startRefInt, err := strconv.Atoi(string(startxref_valF[0]))
	if err != nil {
		messages.E_strconv_atoi(err)
	}

	new_startRefByte := strconv.Itoa(startRefInt + added)

	var buf bytes.Buffer
	buf.Write(fulldata[:flf])
	buf.Write([]byte(new_startRefByte))
	buf.Write(fulldata[fnlf:])

	newXref := buf.Bytes()

	return &newXref
}
