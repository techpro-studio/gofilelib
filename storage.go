package file

import "context"

type Storage interface {
	Upload(ctx context.Context, file InputFile, path string) error
	GetDownloadLink(path string) (string, error)
}
