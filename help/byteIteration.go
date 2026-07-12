package help

import (
	"PolyMixer/messages"
	"bytes"
	"os"
)

// TODO: use this function later to cut the head of and metadata of pdf file and then also check the unused space between those two and use it to embed the whole mp3 binary there
func search_and_count(lookFor, startFrom string, ignoreStartingValue bool, file *os.File) (offset int) {
	fileStat, err := file.Stat()
	if err != nil {
		messages.E_stat_read(err)
	}

	buf_pdf := make([]byte, fileStat.Size())
	_, err = file.ReadAt(buf_pdf, 0)
	if err != nil {
		messages.E_read(err)
	}

	startingIndex := bytes.Index(buf_pdf, []byte(startFrom))
	if startingIndex == -1 {
		messages.E_index("startingIndex")
	}

	sliced_buf_pdf := buf_pdf
	if ignoreStartingValue {
		sliced_buf_pdf = buf_pdf[startingIndex:]
	}
	endIndex := bytes.Index(sliced_buf_pdf, []byte(lookFor))
	if endIndex == -1 {
		messages.E_index("endIndex")
	}

	// TODO: remind me to count offset and empty spaces between startIndex and endIndex later on to be used as
	offset = 0
	space_between := 0
	for i := startingIndex; i < endIndex; offset++ {

	}

	return
}
