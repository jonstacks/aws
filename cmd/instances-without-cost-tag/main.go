package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

var costTag string

func init() {
	flag.StringVar(&costTag, "cost-tag", "cost", "The tag key for determining cost")
	flag.Parse()
}

func hasCostTag(i *ec2.Instance) bool {
	return utils.GetTagValue(i.Tags, costTag) != ""
}

func identifier(i *ec2.Instance) string {
	id := aws.StringValue(i.InstanceId)
	name := utils.GetInstanceName(i)
	if name != "" {
		return fmt.Sprintf("%s (%s)", id, name)
	}
	return id
}

func main() {
	models.Init(models.DefaultSession())

	opts := models.RunningInstancesOpts{IncludeSpot: true}
	all := models.RunningInstances(opts)
	for _, i := range all {
		if !hasCostTag(i) {
			fmt.Println(identifier(i))
		}
	}
}
