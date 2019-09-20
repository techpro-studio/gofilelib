package file

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"os"
	"time"
)

type S3Storage struct {
	AccessKey       string
	AccessSecretKey string
	Endpoint *string
	PresignURLTimeoutSeconds uint8
	Region string
	Bucket  string
}


func NewS3Storage(accessKey string, accessSecretKey string, bucket string, region string, presignURLTimeoutSeconds uint8, endpoint *string) *S3Storage {
	return &S3Storage{AccessKey: accessKey, AccessSecretKey: accessSecretKey,
		Bucket: bucket, Region: region, PresignURLTimeoutSeconds:presignURLTimeoutSeconds, Endpoint:endpoint}
}


func (m *S3Storage) GetSession() *session.Session {
	s, err := session.NewSession(&aws.Config{
		Endpoint:m.Endpoint,
		Region: aws.String(m.Region),
		Credentials: credentials.NewStaticCredentials(
			m.AccessKey,
			m.AccessSecretKey,
			""),
	})
	if err != nil {
		panic(err)
	}
	return s
}

func (m *S3Storage) Download(ctx context.Context, path string, file *os.File)  error {
	downloader := s3manager.NewDownloader(m.GetSession())
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(m.Bucket),
			Key:    aws.String(path),
		})
	return err
}

func (m *S3Storage) Upload(ctx context.Context, file *InputFile, path string)error{
	size := file.Size
	buffer := make([]byte, size)
	_, err := file.Source.Read(buffer)
	if err != nil {
		panic(err)
	}

	_, err = s3.New(m.GetSession()).PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(m.Bucket),
		Key:                aws.String(path),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *S3Storage) GetDownloadLink(path string) (string, error) {
	svc := s3.New(m.GetSession())
	result, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(path),
	})
	return result.Presign(time.Second * time.Duration(m.PresignURLTimeoutSeconds))
}
