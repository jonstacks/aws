package main

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
)

func main() {

	onlyUnmatched := flag.Bool("only-unmatched", false,
		"Only show instance types that are unmatched in running & reserved instnace count.")

	flag.Parse()

	s := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			},
		),
	)

	client := ec2.New(s)
	models.EC2Client(client)

	ris, err := models.ReservedInstances()
	utils.ExitErrorHandler(err)

	opts := models.RunningInstancesOpts{IncludeSpot: false}
	all := models.RunningInstances(opts)

	viewOpts := views.ReservationUtilizationOptions{
		OnlyUnmatched: *onlyUnmatched,
	}
	v := views.NewReservationUtilization(all, ris, viewOpts)
	v.Print()
}
