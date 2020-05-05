package views

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

// SecurityGroupAudit is a view for auditing security groups.
type SecurityGroupAudit struct {
	securityGroups    []*ec2.SecurityGroup
	networkInterfaces []*ec2.NetworkInterface
}

// NewSecurityGroupAudit creates and initializes a new security group audit
func NewSecurityGroupAudit(securityGroups []*ec2.SecurityGroup, networkInterfaces []*ec2.NetworkInterface) *SecurityGroupAudit {
	return &SecurityGroupAudit{
		securityGroups:    securityGroups,
		networkInterfaces: networkInterfaces,
	}
}

func (sga *SecurityGroupAudit) interfaceHasSecurityGroup(ifc *ec2.NetworkInterface, sg *ec2.SecurityGroup) bool {
	if ifc == nil || sg == nil {
		return false
	}
	for _, group := range ifc.Groups {
		if aws.StringValue(group.GroupId) == aws.StringValue(sg.GroupId) {
			return true
		}
	}
	return false
}

// Print prints the SecurityGroupAudit view
func (sga *SecurityGroupAudit) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"VPC ID",
		"Group ID",
		"Group Name",
		"Usages(Number of Interfaces)",
	})

	for _, sg := range sga.securityGroups {
		usages := 0
		for _, ifc := range sga.networkInterfaces {
			if sga.interfaceHasSecurityGroup(ifc, sg) {
				usages++
			}
		}
		table.Append([]string{
			aws.StringValue(sg.VpcId),
			aws.StringValue(sg.GroupId),
			aws.StringValue(sg.GroupName),
			strconv.Itoa(usages),
		})
	}

	table.Render()
}
