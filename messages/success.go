package messages

import (
	"fmt"
	"os"
)

type FileSize int

const (
	BYTE FileSize = iota
	KB
	MB
	GB
)

func S_open_file(file *os.File, fileType string) {
	if fileType == "" {
		fmt.Printf("Success opening %v\n", file.Name())
	}
	fmt.Printf("Success opening %v -> file format %v\n", file.Name(), fileType)
}

func S_file_size(fileType, extraInfo string, size float64, viewSizeIn FileSize) {
	switch viewSizeIn {
	case BYTE:
		fmt.Printf("%v %v size: %.5fB\n", fileType, extraInfo, size)
	case KB:
		fmt.Printf("%v %v size: %.5fKB\n", fileType, extraInfo, size/1024)
	case MB:
		fmt.Printf("%v %v size: %.5fMB\n", fileType, extraInfo, size/1048576)
	case GB:
		fmt.Printf("%v %v size: %.5fGB\n", fileType, extraInfo, size/1073741824)
	}
}
