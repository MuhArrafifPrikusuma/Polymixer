package files

import (
	"os"
)

func (arg *Arguments) ReadHead(b []byte) {
	bytes, err := arg.File1.ReadHead(b)
}
