package files

import (
	"bytes"
	"log"
	"os"

	"PolyMixer/messages"
)

type Arguments struct {
	File1, File2 *os.File
}

func TakeArg(arg *Arguments) {
	if len(os.Args) < 2 {
		log.Fatal("Error: Please insert files")
	}
	if len(os.Args) > 3 {
		log.Fatal("Cannot process more than 2 files... yet")
	}

	var lastObjId, appendToIdx int
	var ptrMp3, ptrPdf *os.File
	var fileType string
	var objRefPtr *Xref_ObjMap_t

	if len(os.Args) >= 3 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			messages.E_open_file(os.Args[1], err)
		}

		arg.File1 = file
		fileType, lastObjId, appendToIdx, ptrMp3, ptrPdf, objRefPtr = getHeader(arg.File1)
		messages.S_open_file(arg.File1, fileType)

		file, err = os.Open(os.Args[2])
		if err != nil {
			arg.File1.Close()
			messages.E_open_file(os.Args[2], err)
		}
		arg.File2 = file
		fileType, tmplastObjId, tmpappendToIdx, tmpptrMp3, tmpptrPdf, tmpobjRefPtr := getHeader(arg.File2)

		if fileType == "PDF" {
			lastObjId = tmplastObjId
			appendToIdx = tmpappendToIdx
			ptrPdf = tmpptrPdf
			objRefPtr = tmpobjRefPtr
		} else if fileType == "MP3" {
			ptrMp3 = tmpptrMp3
		}
		messages.S_open_file(arg.File2, fileType)
		Create_audio_object(objRefPtr, ptrPdf, ptrMp3, appendToIdx, lastObjId)
	}
	Mix_MP3_and_PDF(ptrPdf, ptrMp3, appendToIdx, lastObjId)
}

func getHeader(file *os.File) (strReturn string, lastObjId, appendToIdx int, ptrMp3, ptrPdf *os.File, objRefPtr *Xref_ObjMap_t) {
	buffer := make([]byte, 4)

	_, err := file.ReadAt(buffer, 0)
	if err != nil {
		messages.E_read(err)
	}
	if bytes.HasPrefix(buffer, []byte("ID3")) {
		ptrMp3 = mp3_get_body(file)
		strReturn = "MP3"
		return strReturn, 0, 0, ptrMp3, nil, nil

	} else if bytes.HasPrefix(buffer, []byte("%PDF")) {
		lastObjId, appendToIdx, ptrPdf, objRefPtr = Pdf_open(file)
		strReturn = "PDF"
		return strReturn, lastObjId, appendToIdx, nil, ptrPdf, objRefPtr
	}
	return "unknown file type", 0, 0, nil, nil, nil
}
