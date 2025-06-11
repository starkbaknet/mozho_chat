package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// Service interface for S3 operations
type Service interface {
	UploadFile(file *multipart.FileHeader, entityType, entityId, fileType string, isPublic bool) (*UploadResult, error)
}

type UploadResult struct {
	Key      string
	URL      string
	Size     int64
	MimeType string
}

type S3Service struct {
	client    *s3.Client
	presigner *s3.PresignClient
	isMinio   bool
	public    string
	private   string
	region    string
	endpoint  string
}

func NewS3Service() (*S3Service, error) {
	err := loadEnv()
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     os.Getenv("S3_ACCESS_KEY"),
				SecretAccessKey: os.Getenv("S3_SECRET_KEY"),
				Source:          "env",
			}, nil
		})),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           os.Getenv("S3_ENDPOINT"),
				SigningRegion: os.Getenv("S3_REGION"),
			}, nil
		})),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)

	return &S3Service{
		client:    client,
		presigner: presigner,
		isMinio:   os.Getenv("IS_MINIO") == "true",
		public:    os.Getenv("S3_PUBLIC_BUCKET_NAME"),
		private:   os.Getenv("S3_PRIVATE_BUCKET_NAME"),
		region:    os.Getenv("S3_REGION"),
		endpoint:  os.Getenv("S3_ENDPOINT"),
	}, nil
}

func (s *S3Service) generateFileName(original string) string {
	ext := filepath.Ext(original)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

// UploadFile implements the Service interface
func (s *S3Service) UploadFile(file *multipart.FileHeader, entityType, entityId, fileType string, isPublic bool) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buffer, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	result, err := s.uploadFileBytes(buffer, file.Filename, entityType, entityId, fileType, isPublic)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Key:      result["key"],
		URL:      result["url"],
		Size:     int64(len(buffer)),
		MimeType: result["mimeType"],
	}, nil
}

func (s *S3Service) uploadFileBytes(buffer []byte, originalName, entityType, entityId, fileType string, isPublic bool) (map[string]string, error) {
	if len(buffer) == 0 {
		return nil, errors.New("empty file buffer")
	}

	fileName := s.generateFileName(originalName)
	objectKey := fmt.Sprintf("%s/%s/%s/%s", entityType, entityId, fileType, fileName)
	bucket := s.private
	if isPublic {
		bucket = s.public
	}

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s", s.endpoint, bucket, objectKey)
	if !s.isMinio {
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, s.region, objectKey)
	}

	result := map[string]string{
		"name":     originalName,
		"key":      objectKey,
		"url":      url,
		"public":   fmt.Sprintf("%v", isPublic),
		"mimeType": httpDetectMimeType(buffer),
		"size":     fmt.Sprintf("%d", len(buffer)),
	}

	if isPublic {
		delete(result, "key") // only return key for private
	}

	return result, nil
}

func (s *S3Service) UploadPhoto(buffer []byte, originalName, entityType, entityId, fileType string, isPublic bool, withThumbnail bool) (map[string]string, error) {
	original, err := s.uploadFileBytes(buffer, originalName, entityType, entityId, fileType, isPublic)
	if err != nil {
		return nil, err
	}

	result := map[string]string{
		"original": original["url"],
	}

	if withThumbnail {
		img, _, err := image.Decode(bytes.NewReader(buffer))
		if err != nil {
			return nil, err
		}
		thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
		var thumbBuffer bytes.Buffer
		err = imaging.Encode(&thumbBuffer, thumb, imaging.PNG)
		if err != nil {
			return nil, err
		}

		thumbRes, err := s.uploadFileBytes(thumbBuffer.Bytes(), "thumb-"+originalName, entityType, entityId, fileType, isPublic)
		if err != nil {
			return nil, err
		}
		result["thumbnail"] = thumbRes["url"]
	}

	return result, nil
}

func (s *S3Service) GetPrivateSignedUrl(key string, expiresIn time.Duration) (string, error) {
	bucket := s.private

	req, err := s.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

// simple MIME type detection
func httpDetectMimeType(buffer []byte) string {
	if len(buffer) >= 512 {
		return http.DetectContentType(buffer[:512])
	}
	return http.DetectContentType(buffer)
}

func loadEnv() error {
	if os.Getenv("S3_ACCESS_KEY") == "" {
		return godotenv.Load()
	}
	return nil
}
