package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client is a S3 object storage client.
type Client struct {
	client *minio.Client
	bucket string
}

// Setup returns a client.
func (c *Client) Setup(endpoint string, region string, accessKeyId string, secretAccessKey string, bucket string) error {
	client, err := getClient(endpoint, region, accessKeyId, secretAccessKey)
	if err != nil {
		return err
	}

	c.bucket = bucket
	c.client = client

	return nil
}

func getClient(endpoint string, region string, accessKeyId string, secretAccessKey string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Region: region,
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
