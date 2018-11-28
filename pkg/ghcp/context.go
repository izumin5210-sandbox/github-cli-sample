package ghcp

import (
	"io"
	"os"
)

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func StdIO() IO {
	return IO{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

type Ctx struct {
	IO IO
}
