package messages

import "log"

func E_open_file(fileName string, err error) {
	log.Fatalf("Failed to open %v: %v", fileName, err)
}

func E_read(err error) {
	log.Fatalf("Failed to read %v", err)
}
