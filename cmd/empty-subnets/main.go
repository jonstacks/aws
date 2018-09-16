package main

import (
	"os"

	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {
	models.Init(models.DefaultSession())

	subnets, err := models.Subnets()
	utils.ExitErrorHandler(err)

	var view views.View
	view = views.NewEmptySubnets(subnets)
	view.Render(os.Stdout)
}
