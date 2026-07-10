package main

import (
	"mp3ToPDF/files"
)

func main() {
	args := &files.Arguments{}

	if args.File1 != nil {
		defer args.File1.Close()
	}
	if args.File2 != nil {
		defer args.File2.Close()
	}

	files.TakeArg(args)

}
