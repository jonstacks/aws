package models

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// NetworkInterfaces returns all of the network interfaces or an error if one occured.
func NetworkInterfaces() ([]*ec2.NetworkInterface, error) {
	params := &ec2.DescribeNetworkInterfacesInput{}
	resp, err := ec2Client.DescribeNetworkInterfaces(params)
	return resp.NetworkInterfaces, err
}

// SecurityGroups returns all of the security groups or an error if one occured.
func SecurityGroups() ([]*ec2.SecurityGroup, error) {
	securityGroups := make([]*ec2.SecurityGroup, 0)
	params := &ec2.DescribeSecurityGroupsInput{GroupIds: []*string{}}
	err := ec2Client.DescribeSecurityGroupsPages(params,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			securityGroups = append(securityGroups, page.SecurityGroups...)
			return !lastPage
		})
	return securityGroups, err
}
