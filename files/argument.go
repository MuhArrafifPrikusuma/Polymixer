package files

import (
	"bytes"
	"log"
	"mp3ToPDF/messages"
	"os"
)

type Arguments struct {
	File1, File2 *os.File
}

func TakeArg(arg *Arguments) {
	if len(os.Args) < 2 {
		log.Fatal("Error: Please provide a valid file path")
	}
	if len(os.Args) > 3 {
		log.Fatal("Error: maximum arguments exceeded!")
	}

	if len(os.Args) >= 3 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			messages.E_open_file(os.Args[1], err)
		}
		arg.File1 = file
		fileType := getHeader(arg.File1)
		messages.S_open_file(arg.File1, fileType)

		file, err = os.Open(os.Args[2])
		if err != nil {
			arg.File1.Close()
			messages.E_open_file(os.Args[2], err)
		}
		arg.File2 = file
		fileType = getHeader(arg.File2)
		messages.S_open_file(arg.File2, fileType)

	}
}

func getHeader(file *os.File) string {
	buffer := make([]byte, 4)

	_, err := file.ReadAt(buffer, 0)
	if err != nil {
		messages.E_read(err)
	}
	if bytes.HasPrefix(buffer, []byte("ID3")) {
		return "MP3"
	} else if bytes.HasPrefix(buffer, []byte("%PDF")) {
		return "PDF"
	}
	return "unknown file types"
}
