package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
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

		obj := new(Object)

		obj.Key = object.Key
		obj.ETag = object.ETag

		ret = append(ret, obj)
	}

	return ret, nil
}
