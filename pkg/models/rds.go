package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
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
	ris := make([]*rds.ReservedDBInstance, 0)
	params := &rds.DescribeReservedDBInstancesInput{}

	err := rdsClient.DescribeReservedDBInstancesPages(params,
		func(page *rds.DescribeReservedDBInstancesOutput, lastPage bool) bool {
			ris = append(ris, page.ReservedDBInstances...)
			return !lastPage
		})

	if err != nil {
		return ris, err
	}

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
	instances := make([]*rds.DBInstance, 0)
	params := &rds.DescribeDBInstancesInput{}

	err := rdsClient.DescribeDBInstancesPages(params,
		func(page *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			instances = append(instances, page.DBInstances...)
			return !lastPage
		})

	if err != nil {
		return instances, err
	}

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

// GetRDSLogDownloadURL returns a signed request for the given DB Instance identifier
// and filename
func GetRDSLogDownloadURL(dbInstanceIdentifier string, fileName string) (*http.Request, error) {
	sess := DefaultSession()
	signer := v4.NewSigner(sess.Config.Credentials)
	region := aws.StringValue(sess.Config.Region)
	url := fmt.Sprintf(
		"https://rds.%s.amazonaws.com/v13/downloadCompleteLogFile/%s/%s",
		region,
		dbInstanceIdentifier,
		fileName,
	)

	request, _ := http.NewRequest("GET", url, nil)
	_, err := signer.Presign(request, nil, endpoints.RdsServiceID, region, 1*time.Hour, time.Now())
	return request, err
}

// DescribeDBLogFiles returns details about the log files for a given DB Instance identifier
func DescribeDBLogFiles(dbInstanceIdentifier string) ([]*rds.DescribeDBLogFilesDetails, error) {
	params := &rds.DescribeDBLogFilesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	}
	logFiles := make([]*rds.DescribeDBLogFilesDetails, 0)

	err := rdsClient.DescribeDBLogFilesPages(params,
		func(page *rds.DescribeDBLogFilesOutput, lastPage bool) bool {
			logFiles = append(logFiles, page.DescribeDBLogFiles...)
			return !lastPage
		})

	return logFiles, err
}
