package main

import (
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {
	models.Init(models.DefaultSession())

	sgs, err := models.SecurityGroups()
	utils.ExitErrorHandler(err)

	ifcs, err := models.NetworkInterfaces()
	utils.ExitErrorHandler(err)

	v := views.NewSecurityGroupAudit(sgs, ifcs)
	v.Print()
}
