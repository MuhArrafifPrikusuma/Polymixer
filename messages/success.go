package messages

import (
	"fmt"
	"os"
)

func S_open_file(file *os.File, fileType string) {
	if fileType == "" {
		fmt.Printf("[SUCCESS]Success opening %v\n", file.Name())
	}
	fmt.Printf("[SUCCESS]Success opening %v -> file format %v\n", file.Name(), fileType)
}

func S_file_size(fileType, extraInfo string, size float64) {
	switch {
	case size/1073741824 > 1:
		fmt.Printf("[INFO]%v %v size: %.5fGB\n", fileType, extraInfo, size/1073741824)
	case size/1048576 > 1:
		fmt.Printf("[INFO]%v %v size: %.5fMB\n", fileType, extraInfo, size/1048576)
	case size/1024 > 1:
		fmt.Printf("[INFO]%v %v size: %.5fKB\n", fileType, extraInfo, size/1024)
	default:
		fmt.Printf("[INFO]%v %v size: %.5fB\n", fileType, extraInfo, size)
	}
}

func S_found_at_index(target string, idx int) {
	fmt.Printf("[SUCCESS]Found '%v' at index %v\n", target, idx)
}

func S_found_id(id int) {
	fmt.Printf("[SUCCESS]Found obj id: %v\n", id)
}
