package files

import (
	"PolyMixer/help"
	"PolyMixer/messages"
	"bytes"
	"log"
	"os"
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

	var lastObjId, appendToId int
	var ptrMp3, ptrPdf *os.File
	var fileType string
	if len(os.Args) >= 3 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			messages.E_open_file(os.Args[1], err)
		}
		arg.File1 = file
		fileType, lastObjId, appendToId, ptrMp3, ptrPdf = getHeader(arg.File1)
		messages.S_open_file(arg.File1, fileType)

		file, err = os.Open(os.Args[2])
		if err != nil {
			arg.File1.Close()
			messages.E_open_file(os.Args[2], err)
		}
		arg.File2 = file
		fileType, tmplastObjId, tmpappendToId, tmpptrMp3, tmpptrPdf := getHeader(arg.File2)
		if fileType == "PDF" {
			lastObjId = tmplastObjId
			appendToId = tmpappendToId
			ptrPdf = tmpptrPdf
		} else if fileType == "MP3" {
			ptrMp3 = tmpptrMp3
		}
		messages.S_open_file(arg.File2, fileType)
	}
	help.Mix_MP3_and_PDF(ptrPdf, ptrMp3, appendToId, lastObjId)
}

func getHeader(file *os.File) (strReturn string, lastObjId, appendToId int, ptrMp3, ptrPdf *os.File) {
	buffer := make([]byte, 4)

	_, err := file.ReadAt(buffer, 0)
	if err != nil {
		messages.E_read(err)
	}
	if bytes.HasPrefix(buffer, []byte("ID3")) {
		ptrMp3 = mp3_get_body(file)
		strReturn = "MP3"
		return strReturn, 0, 0, ptrMp3, nil
	} else if bytes.HasPrefix(buffer, []byte("%PDF")) {
		lastObjId, appendToId, ptrPdf = Pdf_open(file)
		strReturn = "PDF"
		return strReturn, lastObjId, appendToId, nil, ptrPdf
	}
	return "unknown file type", 0, 0, nil, nil
}
