package views

// InstanceTypeReservationUtilization keeps track of how many of a particular
// Instance type are running
type InstanceTypeReservationUtilization struct {
	InstanceType string
	NumReserved  int
	NumRunning   int
}

func (i *InstanceTypeReservationUtilization) String() string {
	return i.InstanceType
}

// HasUnused returns whether or not the instance type has reservations that
// aren't currently being used.
func (i *InstanceTypeReservationUtilization) HasUnused() bool {
	return i.NumReserved > i.NumRunning
}

// Unreserved returns the difference between the number of running instances
// and the number of reserved instances. A negative result implies there are
// more instances reserved than running.
func (i *InstanceTypeReservationUtilization) Unreserved() int {
	return i.NumRunning - i.NumReserved
}
