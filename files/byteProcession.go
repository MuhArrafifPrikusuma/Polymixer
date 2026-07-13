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
	fmt.Printf("%v File info: \n", fileInfo.Name())
	messages.S_file_size("MP3", "full size", float64(fileInfo.Size()), messages.MB)
	buf := make([]byte, fileInfo.Size())

	_, err = file.ReadAt(buf, 0)
	if err != nil {
		messages.E_read(err)
	}
	mp3Body := buf[16:]
	messages.S_file_size("MP3", "headless size", float64(len(mp3Body)), messages.MB)
	return mp3Body
}

// func pdf_mutilation(file *os.File) (head, metadata, body []byte) {
//
// }

func header_arrangement(b []byte) {

}
