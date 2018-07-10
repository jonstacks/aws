package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

// Calculates the free ranges in a VPC that can be used to create new subnets.
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

	vpcs, err := models.VPCs()
	utils.ExitErrorHandler(err)

	subnets, err := models.Subnets()
	utils.ExitErrorHandler(err)

	view := views.NewVPCFreeSubnets(vpcs, subnets)
	view.Render(os.Stdout)
}
