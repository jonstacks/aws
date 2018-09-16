package main

import (
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {
	models.Init(models.DefaultSession())

	subnets, err := models.Subnets()
	utils.ExitErrorHandler(err)

	opts := models.RunningInstancesOpts{IncludeSpot: true}
	instances := models.RunningInstances(opts)

	view := views.NewInstancesBySubnet(instances, subnets)
	view.Print()
}
