package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

func hasCostTag(i *ec2.Instance) bool {
	return utils.GetTagValue(i.Tags, "cost") != ""
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
	s := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			},
		),
	)

	client := ec2.New(s)
	models.EC2Client(client)

	opts := models.RunningInstancesOpts{IncludeSpot: true}
	all := models.RunningInstances(opts)
	for _, i := range all {
		if !hasCostTag(i) {
			fmt.Println(identifier(i))
		}
	}
}
