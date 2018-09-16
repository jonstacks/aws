package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

var rdsClient *rds.RDS

// RDSClient sets the client to be used by the models.
func RDSClient(client *rds.RDS) {
	rdsClient = client
}

// ReservedDBInstances returns a slice of active reserved db instances. We
// have to do client side filtering since the API doesn't support filters at
// this time.
func ReservedDBInstances() ([]*rds.ReservedDBInstance, error) {
	params := &rds.DescribeReservedDBInstancesInput{}
	resp, err := rdsClient.DescribeReservedDBInstances(params)
	ris := resp.ReservedDBInstances
	filtered := make([]*rds.ReservedDBInstance, 0)
	for _, ri := range ris {
		if aws.StringValue(ri.State) == "active" {
			filtered = append(filtered, ri)
		}
	}
	return filtered, err
}

// RunningDBInstances returns a slice of running db instances.
func RunningDBInstances() ([]*rds.DBInstance, error) {
	params := &rds.DescribeDBInstancesInput{}
	resp, err := rdsClient.DescribeDBInstances(params)
	instances := resp.DBInstances
	filtered := make([]*rds.DBInstance, 0)
	for _, i := range instances {
		status := aws.StringValue(i.DBInstanceStatus)
		if status == "available" || status == "backing-up" {
			filtered = append(filtered, i)
		}
	}
	return filtered, err
}

// DBSnapshots returns a slice of RDS Snapshots.
func DBSnapshots() ([]*rds.DBSnapshot, error) {
	snapshots := make([]*rds.DBSnapshot, 0)
	params := &rds.DescribeDBSnapshotsInput{}

	err := rdsClient.DescribeDBSnapshotsPages(params,
		func(page *rds.DescribeDBSnapshotsOutput, lastPage bool) bool {
			snapshots = append(snapshots, page.DBSnapshots...)
			return !lastPage
		})
	return snapshots, err
}
