package uploader

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"strings"
)

const (
	ProjectImageFolder = "projects"
	ProfileImageFolder = "profiles"
)

type s3Store struct {
	uploader *s3manager.Uploader
	region   string
	bucket   string
}

type S3StoreConfig struct {
	Region             string
	Bucket             string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

func NewS3Store(config *S3StoreConfig) (ImageUploader, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "s3_user",
		Config: aws.Config{
			Region:      aws.String(config.Region),
			Credentials: credentials.NewStaticCredentials(config.AwsAccessKeyID, config.AwsSecretAccessKey, ""),
			Logger:      aws.NewDefaultLogger(),
		},
	})

	if err != nil {
		return nil, err
	}

	return &s3Store{
		uploader: s3manager.NewUploader(sess),
		region:   config.Region,
		bucket:   config.Bucket,
	}, nil
}

func (s *s3Store) Upload(folder string, file *multipart.FileHeader) (string, error) {
	fileExtension := strings.Split(file.Filename, ".")[1]
	f, err := file.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open file")
		return "", err
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close file")
		}
	}(f)

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate UUID")
		return "", err
	}

	result, err := s.uploader.UploadWithContext(context.TODO(), &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.%s", folder, id.String(), fileExtension)),
		Body:   f,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to upload file")
		return "", err
	}

	return result.Location, nil
}
