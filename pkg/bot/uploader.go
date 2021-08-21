package bot

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewUploader(bucketName string, svc *s3.S3) *uploader {
	return &uploader{
		bucketName: bucketName,
		svc:        svc,
		s3uploader: s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
			u.PartSize = 10 * 1024 * 1024
		}),
	}
}

type uploader struct {
	bucketName string
	svc        *s3.S3
	s3uploader *s3manager.Uploader
}

func (r *uploader) Upload(key string, file io.Reader) error {
	if file == nil {
		return fmt.Errorf("empty file reader")
	}

	_, err := r.s3uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(r.bucketName),
		Key:                aws.String(key),
		Body:               file,
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
