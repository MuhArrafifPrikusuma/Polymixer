package files

import (
	"PolyMixer/messages"
	"os"
)

func mp3_get_body(file *os.File) []byte {
	fileInfo, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}

	buf := make([]byte, fileInfo.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	mp3Body := buf[16:]
	// cut header
	return mp3Body
}

// func pdf_mutilation(file *os.File) (head, metadata, body []byte) {
//
// }

func header_arrangement(b []byte) {

}
