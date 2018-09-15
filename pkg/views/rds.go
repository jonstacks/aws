package views

import (
	"io"
	"os"
	"sort"
	"strconv"
	"time"

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

// RDSSnapshotAudit gives an overview of the RDS snapshots, with their instances,
// and how much storage is being used
type RDSSnapshotAudit struct {
	Snapshots []*rds.DBSnapshot
	Instances []*rds.DBInstance
}

// NewRDSSnapshotAudit returns an audit view for comparing snapshots vs running
// instances
func NewRDSSnapshotAudit(snapshots []*rds.DBSnapshot, instances []*rds.DBInstance) *RDSSnapshotAudit {
	return &RDSSnapshotAudit{Snapshots: snapshots, Instances: instances}
}

// NumRunningInstances returns the number of running instances
func (audit *RDSSnapshotAudit) NumRunningInstances() int {
	return len(audit.Instances)
}

// NumSnapshots returns the number of snapshots
func (audit *RDSSnapshotAudit) NumSnapshots() int {
	return len(audit.Snapshots)
}

// TotalRunningStorageGB returns the total storage of running intances in GB
func (audit *RDSSnapshotAudit) TotalRunningStorageGB() int64 {
	var storage int64
	for _, i := range audit.Instances {
		storage += aws.Int64Value(i.AllocatedStorage)
	}
	return storage
}

// TotalVirtualSnapshotStorageGB returns the total storage of snapshots in GB.
// This is the "virtual" storage though as the snapshots only store deltas
func (audit *RDSSnapshotAudit) TotalVirtualSnapshotStorageGB() int64 {
	var storage int64
	for _, s := range audit.Snapshots {
		storage += aws.Int64Value(s.AllocatedStorage)
	}
	return storage
}

// OldInstancesWithSnapshots returns a map whose keys are DBInstanceIdentifiers
// which no longer exist. The values are slices of *rds.DBSnapshots
func (audit *RDSSnapshotAudit) OldInstancesWithSnapshots() map[string][]*rds.DBSnapshot {
	runningIdentifiers := make(map[string]bool)
	for _, i := range audit.Instances {
		dbIdentifier := aws.StringValue(i.DBInstanceIdentifier)
		runningIdentifiers[dbIdentifier] = true
	}

	old := make(map[string][]*rds.DBSnapshot)
	for _, snap := range audit.Snapshots {
		dbIdentifier := aws.StringValue(snap.DBInstanceIdentifier)
		if _, ok := runningIdentifiers[dbIdentifier]; !ok {
			if _, ok := old[dbIdentifier]; !ok {
				old[dbIdentifier] = make([]*rds.DBSnapshot, 0)
			}

			old[dbIdentifier] = append(old[dbIdentifier], snap)
		}
	}
	return old
}

// Render renders the table to the io.Writer
func (audit *RDSSnapshotAudit) Render(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"DB Instance Identifier",
		"Snapshot Identifier",
		"Created Time",
		"Size",
	})

	oldSnapMap := audit.OldInstancesWithSnapshots()
	for i := range oldSnapMap {
		for _, snap := range oldSnapMap[i] {
			table.Append([]string{
				i,
				aws.StringValue(snap.DBSnapshotIdentifier),
				aws.TimeValue(snap.SnapshotCreateTime).Format(time.UnixDate),
				strconv.Itoa(int(aws.Int64Value(snap.AllocatedStorage))),
			})
		}
	}
	table.Render()
}
