package models

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

// S3Client sets the client to be used by the models.
func S3Client(client *s3.S3) {
	s3Client = client
}

// ListBuckets returns a list of buckets
func ListBuckets() ([]*s3.Bucket, error) {
	resp, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return []*s3.Bucket{}, err
	}
	return resp.Buckets, err
}

// GetBucketLocation gets the bucket's location
// func GetBucketLocation(bucket *s3.Bucket) (*s3.GetBucketLocationOutput, error) {
// 	input := &s3.GetBucketLocationInput{Bucket: bucket.Name}
// 	output, err := s3Client.GetBucketLocation(input)
// 	return output, err
// }

// GetBucketReplication returns the Replication status of the bucket
func GetBucketReplication(bucket *s3.Bucket) (*s3.GetBucketReplicationOutput, error) {
	input := &s3.GetBucketReplicationInput{Bucket: bucket.Name}
	output, err := s3Client.GetBucketReplication(input)
	return output, err
}
