package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
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
	creds := credentials.NewEnvCredentials()
	signer := v4.NewSigner(creds)

	if len(os.Args) < 3 {
		fatal(fmt.Errorf("Not enough arguments supplied. \n\n%s", usage()))
	}

	dbIdentifier := os.Args[1]
	fileName := os.Args[2]
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		fatal(fmt.Errorf("You must supply the environment var 'AWS_DEFAULT_REGION'. \n\n%s", usage()))
	}

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
