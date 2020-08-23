package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/jonstacks/aws/pkg/models"
)

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func usage() string {
	return `Usage:
  rds-logs-download-url <dbInstanceIdentifier> <logName>

Example:
  curl -o postgresql.log.2018-04-05-15 $(rds-log-download-url fa16fmqt5yah7r8 error/postgresql.log.2018-04-05-15)`
}

func main() {
	sess := models.DefaultSession()
	signer := v4.NewSigner(sess.Config.Credentials)

	if len(os.Args) < 3 {
		fatal(fmt.Errorf("Not enough arguments supplied. \n\n%s", usage()))
	}

	dbIdentifier := os.Args[1]
	fileName := os.Args[2]

	region := aws.StringValue(sess.Config.Region)
	url := fmt.Sprintf(
		"https://rds.%s.amazonaws.com/v13/downloadCompleteLogFile/%s/%s",
		region,
		dbIdentifier,
		fileName,
	)

	request, _ := http.NewRequest("GET", url, nil)
	_, err := signer.Presign(request, nil, endpoints.RdsServiceID, region, 1*time.Hour, time.Now())
	if err != nil {
		fatal(err)
	}

	fmt.Println(request.URL)
}
