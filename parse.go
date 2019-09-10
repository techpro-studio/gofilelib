package file

import (
	"github.com/techpro-studio/gohttplib"
	"net/http"
)


func ParseMultipartAndGetFile(r* http.Request, key string, maxMultipartSize int64)(*InputFile, error){
	err := r.ParseMultipartForm(maxMultipartSize)
	if err != nil{
		return nil, gohttplib.HTTP400(err.Error())
	}
	file, header, err := r.FormFile(key)
	if err != nil{
		return nil, gohttplib.HTTP400(err.Error())
	}
	return &InputFile{Source: file, Name:header.Filename, Size:header.Size}, nil
}
