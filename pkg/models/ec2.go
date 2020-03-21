package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ec2Client *ec2.EC2

// EC2Client sets the client to be used by the models.
func EC2Client(client *ec2.EC2) {
	ec2Client = client
}

// Instances retrieves a list of instances by their instance IDs and returns
// an error if one occured.
func Instances(IDs []string) (instances []*ec2.Instance, err error) {
	var resp *ec2.DescribeInstancesOutput

	params := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(IDs),
	}
	resp, err = ec2Client.DescribeInstances(params)
	if err != nil {
		return
	}
	instances = make([]*ec2.Instance, 0)
	for _, r := range resp.Reservations {
		instances = append(instances, r.Instances...)
	}
	return
}

// ReservedInstances returns a slice of reserved instances.
func ReservedInstances() ([]*ec2.ReservedInstances, error) {
	params := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("state"),
				Values: []*string{
					aws.String("active"),
				},
			},
		},
	}
	resp, err := ec2Client.DescribeReservedInstances(params)
	return resp.ReservedInstances, err
}

// RunningInstancesOpts are options that can be passed to the running instances
// call.
type RunningInstancesOpts struct {
	IncludeSpot bool
}

// RunningInstances returns a slice of running instances
func RunningInstances(opts RunningInstancesOpts) []*ec2.Instance {
	instances := make([]*ec2.Instance, 0)
	params := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
	}

	ec2Client.DescribeInstancesPages(params,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, r := range page.Reservations {
				for _, i := range r.Instances {
					if *i.State.Name != "running" {
						continue
					}
					instances = append(instances, i)
				}
			}
			return !lastPage
		})

	if !opts.IncludeSpot {
		filtered := make([]*ec2.Instance, 0)
		for _, i := range instances {
			if aws.StringValue(i.InstanceLifecycle) != "spot" {
				filtered = append(filtered, i)
			}
		}
		instances = filtered
	}
	return instances
}

// SpotInstanceRequests returns a slice of spot instance requests by their IDs
// and an error if one occurs.
func SpotInstanceRequests(requestIDs []string) (reqs []*ec2.SpotInstanceRequest, err error) {
	var resp *ec2.DescribeSpotInstanceRequestsOutput

	params := &ec2.DescribeSpotInstanceRequestsInput{
		SpotInstanceRequestIds: aws.StringSlice(requestIDs),
	}
	resp, err = ec2Client.DescribeSpotInstanceRequests(params)
	reqs = resp.SpotInstanceRequests
	return
}

// Subnets returns a list of subnets
func Subnets() ([]*ec2.Subnet, error) {
	params := &ec2.DescribeSubnetsInput{}
	resp, err := ec2Client.DescribeSubnets(params)
	return resp.Subnets, err
}

// VPCs returns a list of VPCs
func VPCs() ([]*ec2.Vpc, error) {
	params := &ec2.DescribeVpcsInput{}
	resp, err := ec2Client.DescribeVpcs(params)
	return resp.Vpcs, err
}
