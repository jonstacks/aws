package views

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/stretchr/testify/assert"
)

func makeRDSInstance(instanceIdentifier, instanceType, engine string, multiAZ bool) *rds.DBInstance {
	return &rds.DBInstance{
		DBInstanceIdentifier: aws.String(instanceIdentifier),
		DBInstanceClass:      aws.String(instanceType),
		Engine:               aws.String(engine),
		MultiAZ:              aws.Bool(multiAZ),
	}
}

func makeRDSReservation(instanceType, productDescription string, count int) *rds.ReservedDBInstance {
	return &rds.ReservedDBInstance{
		DBInstanceClass:    aws.String(instanceType),
		DBInstanceCount:    aws.Int64(int64(count)),
		ProductDescription: aws.String(productDescription),
		State:              aws.String("active"),
	}
}

func TestGetInstanceFamily(t *testing.T) {
	x := rdsInstanceType("db.m5.large")
	assert.Equal(t, "db.m5", x.family())
	assert.Equal(t, "large", x.size())
	units, ok := x.normalizedUnits()
	assert.True(t, ok)
	assert.Equal(t, 4.0, units)
}

func TestRDSResrvationUtilization(t *testing.T) {
	utilization := NewRDSReservationUtilization(
		[]*rds.DBInstance{
			makeRDSInstance("db-1", "db.m5.large", "postgresql", true),
			makeRDSInstance("db-2", "db.m5.large", "postgresql", false),
			makeRDSInstance("db-3", "db.m5.large", "postgresql", false),
			makeRDSInstance("db-4", "db.m5.large", "aurora-postgresql", true),
		},
		[]*rds.ReservedDBInstance{
			makeRDSReservation("db.m5.large", "postgresql", 3),
			makeRDSReservation("db.r6.2xlarge", "aurora-postgresql", 2),
		},
	)

	// Should have 16 running units. 2 normal larges = 2 * 4 = 8. And 1 MultiAZ large = 2 * 4 = 8, so 16 units total.
	assert.Equal(t, 16.0, utilization.InstanceTypeReservationUtilizations["postgresql/db.m5"].NumRunning)
	assert.Equal(t, 12.0, utilization.InstanceTypeReservationUtilizations["postgresql/db.m5"].NumReserved)
	assert.Equal(t, 4.0, utilization.InstanceTypeReservationUtilizations["postgresql/db.m5"].Unreserved())
	assert.Equal(t, false, utilization.InstanceTypeReservationUtilizations["postgresql/db.m5"].HasUnused())

	// Again, here we expect 8 normalized units, 4 for each large since its MultiAZ.
	assert.Equal(t, 8.0, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.m5"].NumRunning)
	assert.Equal(t, 0.0, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.m5"].NumReserved)
	assert.Equal(t, 8.0, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.m5"].Unreserved())
	assert.Equal(t, false, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.m5"].HasUnused())

	assert.Equal(t, 0.0, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.r6"].NumRunning)
	assert.Equal(t, 32.0, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.r6"].NumReserved)
	assert.Equal(t, -32.0, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.r6"].Unreserved())
	assert.Equal(t, true, utilization.InstanceTypeReservationUtilizations["aurora-postgresql/db.r6"].HasUnused())
}
