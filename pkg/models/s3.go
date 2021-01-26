package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

// S3Client sets the client to be used by the models.
func S3Client(client *s3.S3) {
	s3Client = client
}

type S3Objects []*s3.Object

// TotalSize returns the total size in bytes
func (objects *S3Objects) TotalSize() uint64 {
	var total uint64
	for _, object := range *objects {
		total += uint64(aws.Int64Value(object.Size))
	}
	return total
}

// GetS3ObjectsWithPrefix returns a slice of objects for the bucket with
// the given prefix
func GetS3ObjectsWithPrefix(bucket, prefix string) (S3Objects, error) {
	objects := make([]*s3.Object, 0)
	listParams := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	if prefix != "" {
		listParams.Prefix = aws.String(prefix)
	}

	err := s3Client.ListObjectsV2Pages(listParams, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, o := range page.Contents {
			objects = append(objects, o)
		}
		return !lastPage
	})

	return objects, err
}

func GetS3ObjectsWithPrefixChan(bucket, prefix string, objects chan<- *s3.Object) error {
	listParams := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	if prefix != "" {
		listParams.Prefix = aws.String(prefix)
	}

	err := s3Client.ListObjectsV2Pages(listParams, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, o := range page.Contents {
			objects <- o
		}
		return !lastPage
	})

	close(objects)
	return err
}
