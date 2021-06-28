package main

import (
	"fmt"
	"os"

	"github.com/jonstacks/aws/pkg/cmd"
	"github.com/jonstacks/aws/pkg/models"
)

func usage() string {
	return `Usage:
  rds-logs-download-url <dbInstanceIdentifier> <logName>

Example:
  curl -o postgresql.log.2018-04-05-15 $(rds-log-download-url fa16fmqt5yah7r8 error/postgresql.log.2018-04-05-15)`
}

func main() {
	if len(os.Args) < 3 {
		cmd.HandleError(fmt.Errorf("Not enough arguments supplied. \n\n%s", usage()))
	}

	dbIdentifier := os.Args[1]
	filename := os.Args[2]

	req, err := models.GetRDSLogDownloadURL(dbIdentifier, filename)
	cmd.HandleError(err)

	fmt.Println(req.URL)
}
