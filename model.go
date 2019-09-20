package file

import (
	"io"
)

type InputFile struct{
	Source io.Reader
	Size int64
	Name string
}

type Exported struct {
	ID string
	Size int64
	Path string
	Name string
}