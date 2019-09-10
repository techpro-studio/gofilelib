package file

import "mime/multipart"

type InputFile struct{
	Source multipart.File
	Size int64
	Name string
}

type Exported struct {
	ID string
	Size int64
	Path string
	Name string
}