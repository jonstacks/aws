package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

func sgDirectoryPath(basePath string, sg ec2.SecurityGroup) string {
	var vpcID string
	if sg.VpcId != nil {
		vpcID = aws.StringValue(sg.VpcId)
	}
	if vpcID == "" {
		vpcID = "vpc-default"
	}

	return filepath.Join(basePath, "aws-sg-backup", vpcID)
}

func saveSecurityGroupToFile(path string, contents []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(contents)
	return err
}

func main() {
	models.Init(models.DefaultSession())

	basePath, err := os.Getwd()
	utils.ExitErrorHandler(err)

	sgs, err := models.SecurityGroups()
	utils.ExitErrorHandler(err)

	for _, sg := range sgs {
		sgBasePath := sgDirectoryPath(basePath, *sg)
		err = os.MkdirAll(sgBasePath, 0755)
		utils.ExitErrorHandler(err)

		name := aws.StringValue(sg.GroupName)
		if name == "" {
			name = aws.StringValue(sg.GroupId)
		}
		name = strings.ReplaceAll(name, "/", "-slash-")

		objBytes, err := json.Marshal(sg)
		utils.ExitErrorHandler(err)

		fullPath := fmt.Sprintf("%s.json", filepath.Join(sgBasePath, name))
		err = saveSecurityGroupToFile(fullPath, objBytes)
		utils.ExitErrorHandler(err)
	}
}
