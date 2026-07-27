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
	if len(os.Args) > 4 {
		log.Fatal("Cannot process more than 2 files... yet")
	}

	var lastObjId, cutToIdx, objSize int
	var ptrMp, newObj, bsfXref, ptrPdf *[]byte
	var fileType, fname string

	if len(os.Args) >= 4 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			messages.E_open_file(os.Args[1], err)
		}

		arg.File1 = file
		fileType, lastObjId, cutToIdx, ptrMp, ptrPdf, _, bsfXref = getHeader(arg.File1)
		messages.S_open_file(arg.File1, fileType)

		file, err = os.Open(os.Args[2])
		if err != nil {
			arg.File1.Close()
			messages.E_open_file(os.Args[2], err)
		}
		arg.File2 = file
		fileType, tmplastObjId, tmpcutToIdx, tmpptrMp, tmpptrPdf, _, tmpbsxfXref := getHeader(arg.File2)

		fname = os.Args[3]

		if fileType == "PDF" {
			lastObjId = tmplastObjId
			cutToIdx = tmpcutToIdx // <- this will be used for mixing
			ptrPdf = tmpptrPdf
			bsfXref = tmpbsxfXref
		} else if fileType == "MP3" || fileType == "MP4" {
			ptrMp = tmpptrMp
		}
		messages.S_open_file(arg.File2, fileType)
		newObj, objSize = Create_audio_video_object(ptrMp, lastObjId)
	}
	newXrefSlice := StartXref_refOffset(bsfXref, objSize)
	Mix_MP3_and_PDF(ptrPdf, newXrefSlice, newObj, cutToIdx, fname)
}

func getHeader(file *os.File) (strReturn string, lastObjId, cutToIdx int, ptrMp, ptrPdf *[]byte, objRefPtr *Xref_ObjMap_t, bsfXref *[]byte) {
	buffer := make([]byte, 10)

	_, err := file.ReadAt(buffer, 0)
	if err != nil {
		messages.E_read(err)
	}
	isMp4 := bytes.Index(buffer[:], []byte("ftyp"))
	if bytes.HasPrefix(buffer, []byte("ID3")) {
		ptrMp = mp3_get_body(file)
		strReturn = "MP3"
		return strReturn, 0, 0, ptrMp, nil, nil, nil

	} else if bytes.HasPrefix(buffer, []byte("%PDF")) {
		lastObjId, cutToIdx, ptrPdf, objRefPtr, bsfXref = Pdf_open(file)
		strReturn = "PDF"
		return strReturn, lastObjId, cutToIdx, nil, ptrPdf, objRefPtr, bsfXref
	} else if isMp4 != -1 {
		ptrMp = mp4_read(file)
		strReturn = "MP4"
		return strReturn, 0, 0, ptrMp, nil, nil, nil
	}
	return "unknown file type", 0, 0, nil, nil, nil, nil
}
