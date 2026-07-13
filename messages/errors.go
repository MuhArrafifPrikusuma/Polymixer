package messages

import "log"

func E_open_file(fileName string, err error) {
	log.Fatalf("Failed to open %v: %v\n", fileName, err)
}

func E_read(err error) {
	log.Fatalf("Failed to read %v\n", err)
}

func E_stat_read(err error) {
	log.Fatalf("Failed to retrieve file info %v\n", err)
}

func E_index(str string) {
	log.Fatalf("%v not found!\n", str)
}

func E_byte_slice_too_small() {
	log.Fatal("Byte slice is too small, maybe invalid file?")
}
