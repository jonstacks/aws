package main

import (
	"os"

	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

// Calculates the free ranges in a VPC that can be used to create new subnets.
func main() {
	models.Init(models.DefaultSession())

	vpcs, err := models.VPCs()
	utils.ExitErrorHandler(err)

	subnets, err := models.Subnets()
	utils.ExitErrorHandler(err)

	view := views.NewVPCFreeSubnets(vpcs, subnets)
	view.Render(os.Stdout)
}
