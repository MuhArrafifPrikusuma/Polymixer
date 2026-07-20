package files

import (
	"fmt"
	"os"

	"PolyMixer/messages"
)

// basically cut his head we don't need it
func mp3_get_body(file *os.File) *os.File {
	fileInfo, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	fmt.Printf("[PROCESS]Extracting data from %v\n", fileInfo.Name())
	messages.S_file_size("MP3", "full", float64(fileInfo.Size()))
	buf := make([]byte, fileInfo.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	mp3Body := buf[16:]
	messages.S_file_size("MP3", "headless", float64(len(mp3Body)))

	mp3Cpy, err := os.Open(fileInfo.Name())
	if err != nil {
		messages.E_open_file(fileInfo.Name(), err)
	}
	return mp3Cpy
}

type ObjMap struct {
	// contain [index]id
	obj_and_id map[int]int
	// contain [id]index
	endobjId map[int]int
}

// NOTE: this file has been mutilated way to many times remember to use the full file for embedding mp3
func Pdf_open(file *os.File) (objId, appendMp3At int, pdfCpy *os.File) {
	objMap := &ObjMap{
		obj_and_id: make(map[int]int),
		endobjId:   make(map[int]int),
	}
	fileInfo, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	fmt.Printf("[PROCESS]Extracting data from %v\n", fileInfo.Name())
	byteSlice_toXref, byteSlice_fromXref := Find_xref(file)
	Find_all_obj(byteSlice_toXref, objMap)
	Find_cross_reference_byID(byteSlice_fromXref)
	appendMp3At = Find_spot_for_new_obj(objMap, file)

	pdfCpy, err = os.Open(fileInfo.Name())
	if err != nil {
		messages.E_open_file(fileInfo.Name(), err)
	}

	return objId, appendMp3At, pdfCpy
}
