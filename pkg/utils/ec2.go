package utils

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetInstanceName returns the name for an ec2.Instance. If there is no
// associated name tag, it returns an empty string.
func GetInstanceName(i *ec2.Instance) string {
	for _, t := range i.Tags {
		if aws.StringValue(t.Key) == "Name" {
			return aws.StringValue(t.Key)
		}
	}
	return ""
}

// ExitErrorHandler exits the program with a non-zero exit status if err != nil
func ExitErrorHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
