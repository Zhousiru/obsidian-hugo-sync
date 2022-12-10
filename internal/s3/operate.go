package s3

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
)

const (
	ContentTypeWebP = "image/webp"
	ContentType
)

// List lists all objects with specific prefix.
func (c *Client) List(prefix string) ([]*Object, error) {
	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    prefix,
	}

	ret := []*Object{}

	for object := range c.client.ListObjects(context.Background(), c.bucket, opts) {
		if object.Err != nil {
			return nil, object.Err
		}

		if strings.HasSuffix(object.Key, "/") {
			// ignore directory object
			continue
		}

		obj := new(Object)

		obj.Key = object.Key
		obj.ETag = object.ETag

		ret = append(ret, obj)
	}

	return ret, nil
}

// Put uploads specified bytes object to bucket.
func (c *Client) Put(objectKey string, object []byte, contentType string) error {
	reader := bytes.NewReader(object)
	_, err := c.client.PutObject(context.Background(), c.bucket, objectKey, reader, int64(len(object)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

// Get downloads specified object.
func (c *Client) Get(objectKey string) ([]byte, error) {
	reader, err := c.client.GetObject(context.Background(), c.bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	ret, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Remove removes specified object.
func (c *Client) Remove(objectKey string) error {
	return c.client.RemoveObject(context.Background(), c.bucket, objectKey, minio.RemoveObjectOptions{})
}
