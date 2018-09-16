package main

import (
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {
	models.Init(models.DefaultSession())

	ris, err := models.ReservedDBInstances()
	utils.ExitErrorHandler(err)

	dbs, err := models.RunningDBInstances()
	utils.ExitErrorHandler(err)

	v := views.NewRDSReservationUtilization(dbs, ris)
	v.Print()
}
