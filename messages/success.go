package messages

import (
	"fmt"
	"os"
)

func S_open_file(file *os.File, fileType string) {
	if fileType == "" {
		fmt.Printf("Success opening %v\n", file.Name())
	}
	fmt.Printf("Success opening %v -> file format %v\n", file.Name(), fileType)
}
