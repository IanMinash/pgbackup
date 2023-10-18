package pgbackup

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const RetainDays = 7

// GetS3Client returns a minio.Client that has been preconfigured with
// settings from the environment.
func GetS3Client() *minio.Client {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}

// UploadArchive takes in a filePath and uploads it to an s# bucket
func UploadArchive(ctx context.Context, filePath string) (minio.UploadInfo, error) {
	s3Client := GetS3Client()
	bucketName := os.Getenv("S3_BUCKET_NAME")
	objectName := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
	info, err := s3Client.FPutObject(ctx, bucketName, objectName, filePath,
		minio.PutObjectOptions{
			ContentType:     "application/zip",
			Mode:            minio.Governance,
			RetainUntilDate: time.Now().Add(time.Duration(RetainDays * 24 * time.Hour)),
		})
	if err != nil {
		return info, err
	}
	log.Printf("Successfully uploaded %s. Total size uploaded: %.2f MB. \n", objectName, float64(info.Size)/(1<<20))
	return info, nil
}
