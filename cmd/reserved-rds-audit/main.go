package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {
	s := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			},
		),
	)

	client := rds.New(s)
	models.RDSClient(client)

	ris, err := models.ReservedDBInstances()
	utils.ExitErrorHandler(err)

	dbs, err := models.RunningDBInstances()
	utils.ExitErrorHandler(err)

	v := views.NewRDSReservationUtilization(dbs, ris)
	v.Print()
}
