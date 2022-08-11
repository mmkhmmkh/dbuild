package s3

import (
	"log"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewS3(endpoint, accessKeyID, secretAccessKey string, useSSL bool) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup
}
