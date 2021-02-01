package views

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/olekukonko/tablewriter"
)

// S3ReplicationAudit is a S3ReplicationAudit View
type S3ReplicationAudit struct {
	buckets      []*s3.Bucket
	replications []*s3.GetBucketReplicationOutput
}

// NewS3ReplicationAudit initializes the S3 Replication Audit from the
// buckets.
func NewS3ReplicationAudit(buckets []*s3.Bucket, replications []*s3.GetBucketReplicationOutput) S3ReplicationAudit {
	return S3ReplicationAudit{
		buckets:      buckets,
		replications: replications,
	}
}

// Render implments view.View and renders the table
func (v *S3ReplicationAudit) Render(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Name",
		"Cross Region Replication",
	})

	for i, bucket := range v.buckets {
		replication := v.replications[i]
		replicationStatus := ""
		if replication != nil {
			replicationBuckets := s3ReplicationRules(replication.ReplicationConfiguration.Rules).Buckets()
			replicationStatus = strings.Join(replicationBuckets, ",")
		}
		table.Append([]string{
			aws.StringValue(bucket.Name),
			replicationStatus,
		})
	}
	table.Render()
}

// S3VersioningAudit is a view for checking the status of versioning on buckets
type S3VersioningAudit struct {
	buckets        []*s3.Bucket
	versioningInfo []*s3.GetBucketVersioningOutput
}

// NewS3VersioningAudit initializes the S3 Versioning Audit from the
// buckets.
func NewS3VersioningAudit(buckets []*s3.Bucket, versioningInfo []*s3.GetBucketVersioningOutput) S3VersioningAudit {
	return S3VersioningAudit{
		buckets:        buckets,
		versioningInfo: versioningInfo,
	}
}

// Render implments view.View and renders the table
func (v *S3VersioningAudit) Render(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Name",
		"Versioning Enabled",
		"MFA Delete",
	})

	for i, bucket := range v.buckets {
		versioning := v.versioningInfo[i]
		status := ""
		mfa := ""
		if versioning != nil {
			status = aws.StringValue(versioning.Status)
			mfa = aws.StringValue(versioning.MFADelete)
		}
		table.Append([]string{
			aws.StringValue(bucket.Name),
			status,
			mfa,
		})
	}
	table.Render()
}

type s3ReplicationRules []*s3.ReplicationRule

func (rules s3ReplicationRules) Buckets() []string {
	buckets := make([]string, 0)
	for _, rule := range rules {
		buckets = append(buckets, utils.RemoveARN(rule.Destination.Bucket))
	}
	return buckets
}
