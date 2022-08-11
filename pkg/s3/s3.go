package s3

import (
	"github.com/minio/minio-go"
	"strings"
)

const (
	BUCKET_PREFIX = "dbuild-"
)

func NewS3Client(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func Upload(id, bucketName, path string, client *minio.Client) error {
	bucketExists, err := client.BucketExists(bucketName)
	if err != nil {
		return err
	}
	if !bucketExists {
		err = client.MakeBucket(bucketName, "")
		if err != nil {
			return err
		}
	}

	parts := strings.Split(path, "/")

	_, err = client.FPutObject(bucketName, parts[len(parts)-1]+"-"+BUCKET_PREFIX+id, path, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return nil
}
