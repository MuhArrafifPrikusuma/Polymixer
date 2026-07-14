package messages

import "log"

func E_open_file(fileName string, err error) {
	log.Fatalf("[ERROR]Failed to open %v: %v\n", fileName, err)
}

func E_read(err error) {
	log.Fatalf("[ERROR]Failed to read %v\n", err)
}

func E_stat_read(err error) {
	log.Fatalf("[ERROR]Failed to retrieve file info %v\n", err)
}

func E_index(str string) {
	log.Fatalf("[ERROR]%v not found!\n", str)
}

func E_byte_slice_too_small(size int) {
	log.Fatalf("[ERROR]Byte slice %v is too small, maybe invalid file?\n", size)
}

func E_strconv_atoi(err error) {
	log.Fatal("[ERROR]", err)
}
