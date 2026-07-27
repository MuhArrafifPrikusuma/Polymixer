package files

import (
	"fmt"
	"os"

	"PolyMixer/messages"
)

// basically cut his head we don't need it
func mp3_get_body(file *os.File) *[]byte {
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
	mp3Body := buf[10:]
	messages.S_file_size("MP3", "headless", float64(len(mp3Body)))

	mp3Cpy, err := os.Open(fileInfo.Name())
	if err != nil {
		messages.E_open_file(fileInfo.Name(), err)
	}
	messages.S_open_file(mp3Cpy, "MP3")
	return &mp3Body
}

type ObjMap_t struct {
	// contain [IDx]index of ID
	objIdx_and_ID map[int]int
	// contain [ID]index of endobj
	endobjId map[int]int
}

func Pdf_open(file *os.File) (objId, cutTo int, pdfCpy *[]byte, objRefPtr *Xref_ObjMap_t, bsfXref *[]byte) {
	objMap := &ObjMap_t{
		objIdx_and_ID: make(map[int]int),
		endobjId:      make(map[int]int),
	}
	fileInfo, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}
	fmt.Printf("[PROCESS]Extracting data from %v\n", fileInfo.Name())

	byteSlice_toXref, byteSlice_fromXref, xref_start_at := Find_xref(file)

	Find_all_obj(byteSlice_toXref, objMap)

	objRefPtr, firstFoundID := Find_ID_reference(byteSlice_fromXref, objMap, xref_start_at)

	cutTo = Cut_HEAD_to(objMap, file, firstFoundID)

	return objId, cutTo, byteSlice_toXref, objRefPtr, byteSlice_fromXref
}

func Create_audio_object(ptrMP3 *[]byte, lastObjID int) (obj *[]byte, size int) {
	fmt.Println("[PROCESS START]creating new object...")
	defer fmt.Println("[PROCESS END]Object created successfully")
	content := *ptrMP3
	objID := lastObjID + 1
	obj_HEADER := fmt.Sprintf("%d 0 obj\n", objID)
	obj_METADATA := fmt.Sprintf("<</Length %d>>\n", len(content))

	mergeOBJ := []byte(obj_HEADER + obj_METADATA + "stream\n")
	mergeOBJ = append(mergeOBJ, content...)
	mergeOBJ = append(mergeOBJ, []byte("\nendstream\nendobj\n")...)

	objSize := len(mergeOBJ)

	return &mergeOBJ, objSize
}

func Mix_MP3_and_PDF(filePdf, bsfXref, mp3Obj *[]byte, cutToIDX int, newName string) {
	fullPDF := *filePdf
	pdfHEAD := fullPDF[:cutToIDX+1]
	pdfBODY := fullPDF[cutToIDX:]

	new_PDF := pdfHEAD
	new_PDF = append(new_PDF, *mp3Obj...)
	new_PDF = append(new_PDF, pdfBODY...)
	new_PDF = append(new_PDF, *bsfXref...)

	newfile, err := os.Create(newName)
	if err != nil {
		return
	}

	defer newfile.Close()

	_, err = newfile.Write(new_PDF)

	fmt.Println("[SUCCESS]Successfully merged both files")
}
