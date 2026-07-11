package files

import (
	"fmt"
	"log"
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
			log.Fatalf("Failed to open file 1: %v\n", err)
		}
		fmt.Println("Succesfully opened:", arg.File1.Name())

		file, err = os.Open(os.Args[2])
		if err != nil {
			arg.File1.Close()
			log.Fatalf("Failed to open file 2: %v\n", err)
		}

		arg.File2 = file
		fmt.Println("Succesfully opened:", arg.File2.Name())

	}
}

func getHeader() {

}
