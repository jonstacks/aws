package main

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {
	models.S3Client(s3.New(models.DefaultSession()))

	buckets, err := models.ListBuckets()
	utils.ExitErrorHandler(err)

	replications := make([]*s3.GetBucketReplicationOutput, len(buckets))
	for i, bucket := range buckets {
		replication, err := models.GetBucketReplication(bucket)
		if err != nil {
			replication = nil
		}
		replications[i] = replication
	}

	audit := views.NewS3ReplicationAudit(buckets, replications)
	audit.Render(os.Stdout)
}
