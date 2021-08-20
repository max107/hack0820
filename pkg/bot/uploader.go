package bot

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
)

func NewS3Uploader(bucketName string, svc *s3.S3) *s3Uploader {
	uploader := s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
	})

	return &s3Uploader{
		bucketName: bucketName,
		svc:        svc,
		s3uploader: uploader,
	}
}

type s3Uploader struct {
	bucketName string
	svc        *s3.S3
	s3uploader *s3manager.Uploader
}

func (r *s3Uploader) Upload(key string, file io.Reader) error {
	if file == nil {
		return fmt.Errorf("empty file reader")
	}

	input := &s3manager.UploadInput{
		Bucket:             aws.String(r.bucketName),
		Key:                aws.String(key),
		Body:               file,
		ContentDisposition: aws.String("attachment"),
	}

	resp, err := r.s3uploader.Upload(input)
	log.Info().Msgf("resp: %v %v", resp, err)
	return err
}
