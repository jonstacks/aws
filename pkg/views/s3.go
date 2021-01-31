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

type s3ReplicationRules []*s3.ReplicationRule

func (rules s3ReplicationRules) Buckets() []string {
	buckets := make([]string, 0)
	for _, rule := range rules {
		buckets = append(buckets, utils.RemoveARN(rule.Destination.Bucket))
	}
	return buckets
}
