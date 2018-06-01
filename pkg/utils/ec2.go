package utils

import (
	"fmt"
	"math"
	"net"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetInstanceName returns the name for an ec2.Instance. If there is no
// associated name tag, it returns an empty string.
func GetInstanceName(i *ec2.Instance) string {
	for _, t := range i.Tags {
		if aws.StringValue(t.Key) == "Name" {
			return aws.StringValue(t.Key)
		}
	}
	return ""
}

// ExitErrorHandler exits the program with a non-zero exit status if err != nil
func ExitErrorHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

// GetTagValue iterates over the tags and returns the value of the given tag
// if it exists, otherwise ""
func GetTagValue(tags []*ec2.Tag, key string) string {
	for _, t := range tags {
		if aws.StringValue(t.Key) == key {
			return aws.StringValue(t.Value)
		}
	}
	return ""
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// SubnetSize calculates the number of addresses in a given CIDR
func SubnetSize(cidr string) (int, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, err
	}

	ones, bits := ipnet.Mask.Size()
	exponent := bits - ones
	return powInt(2, exponent), nil
}

// IsSubnetEmpty returns true/false depending on if the subnet is empty. This is
// tailored to Amazon's implementation where they reserve 5 IP addresses per
// subnet. So, we will consider the subnet empty if its available IP Addresses
// is equal to the subnet size minus 5
func IsSubnetEmpty(subnet *ec2.Subnet) bool {
	size, err := SubnetSize(aws.StringValue(subnet.CidrBlock))
	if err != nil {
		return false
	}

	emptySubnetSize := int64(size - 5)

	return aws.Int64Value(subnet.AvailableIpAddressCount) == emptySubnetSize
}
