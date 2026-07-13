package files

import (
	"PolyMixer/messages"
	"fmt"
	"os"
)

// basically cut his head we don't need it
func mp3_get_body(file *os.File) []byte {
	fileInfo, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}

	buf := make([]byte, fileInfo.Size())
	fmt.Printf("MP3 file size: %vMB\n", len(buf)/1048576)

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	mp3Body := buf[16:]
	fmt.Printf("Headless size: %vMB\n", len(mp3Body)/1048576)
	return mp3Body
}

// func pdf_mutilation(file *os.File) (head, metadata, body []byte) {
//
// }

func header_arrangement(b []byte) {

}
