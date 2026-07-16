package files

import (
	"fmt"
	"os"

	"PolyMixer/help"
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

// NOTE: this file has been mutilated way to many times remember to use the full file for embedding mp3
func Pdf_open(file *os.File) (objId, appendMp3At int, pdfCpy *os.File) {
	fileInfo, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	fmt.Printf("[PROCESS]Extracting data from %v\n", fileInfo.Name())
	xrefStartIdx, rawBytes := help.Find_xref(file)
	endObjIdx, objIdx := help.Find_last_obj_idx(xrefStartIdx, rawBytes)
	objId = help.Find_obj_id(objIdx, rawBytes)
	appendMp3At = help.Find_spot_for_new_obj(endObjIdx, xrefStartIdx, file)

	pdfCpy, err = os.Open(fileInfo.Name())
	if err != nil {
		messages.E_open_file(fileInfo.Name(), err)
	}

	return objId, appendMp3At, pdfCpy
}
