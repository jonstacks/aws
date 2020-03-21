package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

func main() {
	spotInstaceRequestIDs := os.Args[1:]

	models.Init(models.DefaultSession())

	requests, err := models.SpotInstanceRequests(spotInstaceRequestIDs)
	utils.ExitErrorHandler(err)

	instanceIDs := make([]string, 0)

	for _, request := range requests {
		instanceIDs = append(instanceIDs, aws.StringValue(request.InstanceId))
	}

	instances, err := models.Instances(instanceIDs)
	utils.ExitErrorHandler(err)

	for _, i := range instances {
		if i != nil && i.PrivateIpAddress != nil {
			fmt.Println(aws.StringValue(i.PrivateIpAddress))
		}
	}
}
