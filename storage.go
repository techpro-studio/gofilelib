package file

import (
	"context"
	"os"
)

type Storage interface {
	Upload(ctx context.Context, file *InputFile, path string) error
	Download(ctx context.Context, path string, file *os.File)  error
	GetDownloadLink(path string) (string, error)
}
