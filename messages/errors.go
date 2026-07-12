package messages

import "log"

func E_open_file(fileName string, err error) {
	log.Fatalf("Failed to open %v: %v", fileName, err)
}

func E_read(err error) {
	log.Fatalf("Failed to read %v", err)
}

func E_stat_read(err error) {
	log.Fatalf("Failed to retrieve file info %v", err)
}

func E_index(str string) {
	log.Fatalf("%v not found!", str)
}
