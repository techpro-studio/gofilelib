package file

import (
	"context"
	"fmt"
	"github.com/techpro-studio/gohttplib"
	"time"
)

type UseCase interface {
	Upload(ctx context.Context, file *InputFile, directory string)(string, error)
	GetDownloadLink(ctx context.Context, fileId string)(string, error)
	GetFiles(ctx context.Context, fileIds []string) ([]Exported, error)
}

type DefaultUseCase struct {
	storage Storage
	repository Repository
}

func (d *DefaultUseCase) GetFiles(ctx context.Context, fileIds []string) ([]Exported, error) {
	files := d.repository.GetFiles(ctx, fileIds)
	if len(files) != len(fileIds){
		return nil, gohttplib.HTTP400("Mismatch files and fileIds length")
	}
	return files, nil
}

func (d *DefaultUseCase) Upload(ctx context.Context, file *InputFile, directory string) (string, error) {
	path := fmt.Sprintf("%s/%d.%s", directory, time.Now().Unix(), file.Name)
	err := d.storage.Upload(ctx, *file, path)
	if err != nil{
		return "", err
	}
	id := d.repository.CreateFile(ctx, file.Name, file.Size, path)
	return id, nil
}

func (d *DefaultUseCase) GetDownloadLink(ctx context.Context, fileId string) (string, error) {
	path := d.repository.GetFilePath(ctx, fileId)
	if path == nil{
		return "", gohttplib.HTTP404(fileId)
	}
	return d.storage.GetDownloadLink(*path)
}

func NewDefaultUseCase(storage Storage, repository Repository) *DefaultUseCase {
	return &DefaultUseCase{storage: storage, repository:repository}
}
