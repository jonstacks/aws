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

	versioningInfo := make([]*s3.GetBucketVersioningOutput, len(buckets))
	for i, bucket := range buckets {
		versioning, err := models.GetBucketVersioning(bucket)
		if err != nil {
			versioning = nil
		}
		versioningInfo[i] = versioning
	}

	audit := views.NewS3VersioningAudit(buckets, versioningInfo)
	audit.Render(os.Stdout)
}
