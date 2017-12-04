package views

import (
	"os"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

// ReservationUtilization shows which instance types & families we are utilizing
// instances in.
type ReservationUtilization struct {
	Running                             []*ec2.Instance
	Reservations                        []*ec2.ReservedInstances
	InstanceTypeReservationUtilizations map[string]*InstanceTypeReservationUtilization
}

// NewReservationUtilization Creates a new view for the reserved utilization.
func NewReservationUtilization(running []*ec2.Instance, reservations []*ec2.ReservedInstances) *ReservationUtilization {
	ru := ReservationUtilization{
		Running:                             running,
		Reservations:                        reservations,
		InstanceTypeReservationUtilizations: make(map[string]*InstanceTypeReservationUtilization),
	}

	for _, i := range ru.Running {
		itype := aws.StringValue(i.InstanceType)
		iru := ru.getOrInitializeITypeReservation(itype)
		iru.NumRunning++
	}

	for _, r := range ru.Reservations {
		itype := aws.StringValue(r.InstanceType)
		iru := ru.getOrInitializeITypeReservation(itype)
		iru.NumReserved += int(*r.InstanceCount)
	}

	return &ru
}

func (ru *ReservationUtilization) getOrInitializeITypeReservation(s string) *InstanceTypeReservationUtilization {
	_, ok := ru.InstanceTypeReservationUtilizations[s]
	if !ok {
		ru.InstanceTypeReservationUtilizations[s] = &InstanceTypeReservationUtilization{
			InstanceType: s,
		}
	}
	return ru.InstanceTypeReservationUtilizations[s]
}

// SortedInstanceTypes returns a sorted slice of the instance types
func (ru *ReservationUtilization) SortedInstanceTypes() []string {
	types := make([]string, 0)
	for k := range ru.InstanceTypeReservationUtilizations {
		types = append(types, k)
	}
	sort.Strings(types)
	return types
}

// Print prints the table to stdout
func (ru *ReservationUtilization) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Instance Type",
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
