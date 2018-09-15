package main

import (
	"fmt"
	"os"

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

	dbs, err := models.RunningDBInstances()
	utils.ExitErrorHandler(err)

	snapshots := models.DBSnapshots()

	view := views.NewRDSSnapshotAudit(snapshots, dbs)

	fmt.Printf("Number of Running DBs: %d\n", view.NumRunningInstances())
	fmt.Printf("Number of DB Snapshots: %d\n", view.NumSnapshots())
	fmt.Printf("Total Running storage: %d GB\n", view.TotalRunningStorageGB())
	fmt.Printf("Total Virtual Snapshot storage: %d GB\n", view.TotalVirtualSnapshotStorageGB())

	view.Render(os.Stdout)
}
