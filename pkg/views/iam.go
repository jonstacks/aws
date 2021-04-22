package views

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/olekukonko/tablewriter"
)

type IAMAccessKeysAudit struct {
	Users    map[string][]*iam.AccessKeyMetadata
	LastUsed map[string]*iam.AccessKeyLastUsed
}

// Print prints the table to stdout
func (v *IAMAccessKeysAudit) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"User Name",
		"Access Key ID",
		"Status",
		"Age Days",
		"Last Used",
		"Last Used Service",
	})

	for username, keys := range v.Users {
		for _, key := range keys {
			if key != nil {
				accessKeyID := aws.StringValue(key.AccessKeyId)
				keyAge := time.Since(*key.CreateDate)
				var lastUsed string
				var lastUsedService string
				if v.LastUsed[accessKeyID] != nil && v.LastUsed[accessKeyID].LastUsedDate != nil {
					lastUsed = v.LastUsed[accessKeyID].LastUsedDate.String()
					lastUsedService = aws.StringValue(v.LastUsed[accessKeyID].ServiceName)
				}
				table.Append([]string{
					username,
					accessKeyID,
					aws.StringValue(key.Status),
					fmt.Sprintf("%0.2f", (keyAge.Hours() / 24)),
					lastUsed,
					lastUsedService,
				})
			}
		}

	}

	table.Render()
}
