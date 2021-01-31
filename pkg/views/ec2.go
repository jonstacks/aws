package views

import (
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/networktree"
	"github.com/olekukonko/tablewriter"
)

// ReservationUtilization shows which instance types & families we are utilizing
// instances in.
type ReservationUtilization struct {
	Running                             []*ec2.Instance
	Reservations                        []*ec2.ReservedInstances
	InstanceTypeReservationUtilizations map[string]*InstanceTypeReservationUtilization

	opts ReservationUtilizationOptions
}

// ReservationUtilizationOptions are options which modify the
// ReservationUtilization view.
type ReservationUtilizationOptions struct {
	OnlyUnmatched bool
}

// NewReservationUtilization Creates a new view for the reserved utilization.
func NewReservationUtilization(running []*ec2.Instance, reservations []*ec2.ReservedInstances, opts ReservationUtilizationOptions) *ReservationUtilization {
	ru := ReservationUtilization{
		Running:                             running,
		Reservations:                        reservations,
		InstanceTypeReservationUtilizations: make(map[string]*InstanceTypeReservationUtilization),
		opts:                                opts,
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

	if opts.OnlyUnmatched {
		ru.pruneMatched()
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

// pruneMatched removes all the keys where the running count = reserved count
func (ru *ReservationUtilization) pruneMatched() {
	for k, v := range ru.InstanceTypeReservationUtilizations {
		// This instnace type is perfectly matched, lets delete it.
		if v.Unreserved() == 0 {
			delete(ru.InstanceTypeReservationUtilizations, k)
		}
	}
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

// InstancesBySubnet is a view for showing the instances grouped by subnet.
type InstancesBySubnet struct {
	subnets   map[string]*ec2.Subnet
	instances map[string][]*ec2.Instance
}

// NewInstancesBySubnet creates a new view from the instances and subnets
func NewInstancesBySubnet(instances []*ec2.Instance, subnets []*ec2.Subnet) *InstancesBySubnet {
	ibs := &InstancesBySubnet{
		subnets:   make(map[string]*ec2.Subnet),
		instances: make(map[string][]*ec2.Instance),
	}
	for _, i := range instances {
		ibs.AddInstance(i)
	}
	for _, s := range subnets {
		ibs.AddSubnet(s)
	}
	return ibs
}

// AddInstance adds a new instance to the view
func (ibs *InstancesBySubnet) AddInstance(i *ec2.Instance) {
	subnetID := aws.StringValue(i.SubnetId)
	if _, ok := ibs.instances[subnetID]; !ok {
		ibs.instances[subnetID] = make([]*ec2.Instance, 0)
	}
	ibs.instances[subnetID] = append(ibs.instances[subnetID], i)
}

// AddSubnet adds a new subnet to the view
func (ibs *InstancesBySubnet) AddSubnet(s *ec2.Subnet) {
	subnetID := aws.StringValue(s.SubnetId)
	ibs.subnets[subnetID] = s
}

// Print prints the InstancesBySubnet view
func (ibs *InstancesBySubnet) Print() {
	emptySubnets := make([]*ec2.Subnet, 0)
	instanceCount := 0

	subnetTemplate := func(s *ec2.Subnet) string {
		subnetID := aws.StringValue(s.SubnetId)
		name := utils.GetTagValue(s.Tags, "Name")
		cidr := aws.StringValue(s.CidrBlock)
		return fmt.Sprintf("%s [Name=%s][CIDR=%s]", subnetID, name, cidr)
	}

	for subnetID, subnet := range ibs.subnets {
		instances, ok := ibs.instances[subnetID]
		if !ok {
			// Subnet is Empty
			emptySubnets = append(emptySubnets, subnet)
			continue
		}

		fmt.Printf("--- %s ---\n", subnetTemplate(subnet))
		for _, i := range instances {
			instanceCount++
			fmt.Printf("   * %s (%s)\n", aws.StringValue(i.InstanceId), utils.GetTagValue(i.Tags, "Name"))
		}
		fmt.Println()
	}

	fmt.Println("------ Empty Subnets ----")
	for _, subnet := range emptySubnets {
		fmt.Printf("%s\n", subnetTemplate(subnet))
	}

	fmt.Println("----- Summary -----")
	fmt.Printf("%d Subnets\n", len(ibs.subnets))
	fmt.Printf("%d Empty Subnets\n", len(emptySubnets))
	fmt.Printf("%d Total Instances\n", instanceCount)
}

// EmptySubnets is a view which shows subnets that are empty
type EmptySubnets struct {
	subnets []*ec2.Subnet
}

// NewEmptySubnets creates a new empty subnets view
func NewEmptySubnets(subnets []*ec2.Subnet) *EmptySubnets {
	return &EmptySubnets{subnets: subnets}
}

// Render implements views.View
func (es *EmptySubnets) Render(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"ID",
		"Name",
		"CIDR",
		"Available IPs",
		"Subnet Size",
		"State",
		"VPC ID",
	})
	for _, s := range es.subnets {
		if utils.IsSubnetEmpty(s) {
			subnetSize, _ := utils.SubnetSize(aws.StringValue(s.CidrBlock))
			table.Append([]string{
				aws.StringValue(s.SubnetId),
				utils.GetTagValue(s.Tags, "Name"),
				aws.StringValue(s.CidrBlock),
				strconv.FormatInt(aws.Int64Value(s.AvailableIpAddressCount), 10),
				strconv.Itoa(subnetSize),
				aws.StringValue(s.State),
				aws.StringValue(s.VpcId),
			})
		}
	}

	table.Render()
}

// VPCFreeSubnets gives you available subnet ranges for a VPC
type VPCFreeSubnets struct {
	vpcSubnetMap map[*ec2.Vpc][]*ec2.Subnet
}

// NewVPCFreeSubnets creates a new view for showing free subnets in each VPC
func NewVPCFreeSubnets(vpcs []*ec2.Vpc, subnets []*ec2.Subnet) *VPCFreeSubnets {
	vfs := &VPCFreeSubnets{
		vpcSubnetMap: make(map[*ec2.Vpc][]*ec2.Subnet),
	}

	for _, vpc := range vpcs {
		for _, subnet := range subnets {
			if aws.StringValue(subnet.VpcId) == aws.StringValue(vpc.VpcId) {
				vfs.addSubnet(vpc, subnet)
			}
		}
	}

	return vfs
}

func (vfs *VPCFreeSubnets) addSubnet(v *ec2.Vpc, s *ec2.Subnet) {
	if _, ok := vfs.vpcSubnetMap[v]; !ok {
		vfs.vpcSubnetMap[v] = make([]*ec2.Subnet, 0)
	}
	vfs.vpcSubnetMap[v] = append(vfs.vpcSubnetMap[v], s)
}

// Render renders the view and implements view
func (vfs *VPCFreeSubnets) Render(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"VPC ID",
		"VPC Name",
		"VPC CIDR",
		"Available Subnets",
	})

	for vpc, subnets := range vfs.vpcSubnetMap {
		cidr := aws.StringValue(vpc.CidrBlock)
		tree, err := networktree.New(cidr)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			continue
		}
		for _, subnet := range subnets {
			subCidr := aws.StringValue(subnet.CidrBlock)
			_, sn, err := net.ParseCIDR(subCidr)
			if err != nil {
				continue
			}
			subtree := tree.Find(sn)
			if subtree != nil {
				subtree.MarkUsed()
			}
		}

		unusedRanges := tree.UnusedRanges()
		unusedCIDRs := make([]string, len(unusedRanges))
		for i, n := range unusedRanges {
			unusedCIDRs[i] = n.String()
		}

		table.Append([]string{
			aws.StringValue(vpc.VpcId),
			utils.GetTagValue(vpc.Tags, "Name"),
			aws.StringValue(vpc.CidrBlock),
			strings.Join(unusedCIDRs, "\n"),
		})
	}

	table.Render()
}
