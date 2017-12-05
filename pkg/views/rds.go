package views

import (
	"os"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/olekukonko/tablewriter"
)

// RDSReservationUtilization shows which instance types & families we are
// utilizing instances in.
type RDSReservationUtilization struct {
	Running                             []*rds.DBInstance
	Reservations                        []*rds.ReservedDBInstance
	InstanceTypeReservationUtilizations map[string]*InstanceTypeReservationUtilization
}

// NewRDSReservationUtilization Creates a new view for the reserved utilization.
func NewRDSReservationUtilization(running []*rds.DBInstance, reservations []*rds.ReservedDBInstance) *RDSReservationUtilization {
	ru := RDSReservationUtilization{
		Running:                             running,
		Reservations:                        reservations,
		InstanceTypeReservationUtilizations: make(map[string]*InstanceTypeReservationUtilization),
	}

	for _, i := range ru.Running {
		itype := aws.StringValue(i.DBInstanceClass)
		iru := ru.getOrInitializeITypeReservation(itype)
		iru.NumRunning++
	}

	for _, r := range ru.Reservations {
		itype := aws.StringValue(r.DBInstanceClass)
		iru := ru.getOrInitializeITypeReservation(itype)
		iru.NumReserved += int(*r.DBInstanceCount)
	}

	return &ru
}

func (ru *RDSReservationUtilization) getOrInitializeITypeReservation(s string) *InstanceTypeReservationUtilization {
	_, ok := ru.InstanceTypeReservationUtilizations[s]
	if !ok {
		ru.InstanceTypeReservationUtilizations[s] = &InstanceTypeReservationUtilization{
			InstanceType: s,
		}
	}
	return ru.InstanceTypeReservationUtilizations[s]
}

// SortedInstanceTypes returns a sorted slice of the instance types
func (ru *RDSReservationUtilization) SortedInstanceTypes() []string {
	types := make([]string, 0)
	for k := range ru.InstanceTypeReservationUtilizations {
		types = append(types, k)
	}
	sort.Strings(types)
	return types
}

// Print prints the table to stdout
func (ru *RDSReservationUtilization) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"DB Instance Class",
		"Running Count",
		"Reserved Count",
		"Has Unused",
		"Should be reserved?",
	})

	for _, k := range ru.SortedInstanceTypes() {
		iru := ru.getOrInitializeITypeReservation(k)
		extra := ""
		if iru.HasUnused() {
			extra = "X"
		}
		table.Append([]string{
			k,
			strconv.Itoa(iru.NumRunning),
			strconv.Itoa(iru.NumReserved),
			extra,
			strconv.Itoa(iru.Unreserved()),
		})
	}

	table.Render()
}
