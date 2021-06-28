package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jonstacks/aws/pkg/cmd"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/sirupsen/logrus"
)

func usage() string {
	return `Usage:
  sync-rds-logs <dbInstanceIdentifier> <directory>

Example:
  sync-rds-logs some-identifier /my/log/directory
`
}

func main() {
	if len(os.Args) < 3 {
		cmd.HandleError(
			fmt.Errorf("Not enough arguments supplied. \n\n%s", usage()),
		)
	}

	dbIdentifier := os.Args[1]
	directory := os.Args[2]

	models.Init(models.DefaultSession())

	logFiles, err := models.DescribeDBLogFiles(dbIdentifier)
	cmd.HandleError(err)
	logrus.Infof("Found %d log files for DBInstanceIdentifier=%s", len(logFiles), dbIdentifier)

	cmd.HandleError(os.MkdirAll(directory, os.ModePerm))

	for _, logFile := range logFiles {
		fileName := aws.StringValue(logFile.LogFileName)
		remoteSize := aws.Int64Value(logFile.Size)
		localFilePath := filepath.Join(directory, filepath.Base(fileName))

		downloadLog := func(dbIdentifier, fileName string) {
			logrus.Infof("Downloading %s to '%s'", fileName, localFilePath)
			req, err := models.GetRDSLogDownloadURL(dbIdentifier, fileName)
			cmd.HandleError(err)
			cmd.HandleError(utils.DownloadFile(localFilePath, req.URL.String()))
			logrus.Infof("Successfully downloaded %s", fileName)
		}

		stat, err := os.Stat(localFilePath)
		if err == nil {
			// File Exists
			localSize := stat.Size()
			logrus.Infof("%s (Local Size: %d, Remote Size: %d)", fileName, localSize, remoteSize)
			// Now check if remote size is bigger
			if remoteSize > localSize {
				downloadLog(dbIdentifier, fileName)
			} else {
				logrus.Infof("Not downloading %s as filesize matches remote", fileName)
			}
		} else {
			downloadLog(dbIdentifier, fileName)
		}
	}
}
