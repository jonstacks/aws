package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

// Dump IAM policies to stdout or redirect to file with shell redirect
func main() {
	models.Init(models.DefaultSession())

	policies, err := models.ListPolicies()
	utils.ExitErrorHandler(err)

	fmt.Printf("Got %d policies\n", len(policies))

	for _, policy := range policies {
		var detail *iam.GetPolicyVersionOutput
		for i := 0; i < 3; i++ {
			detail, err = models.DescribePolicy(*policy.Arn, *policy.DefaultVersionId)
			if err == nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
		utils.ExitErrorHandler(err)

		doc := detail.PolicyVersion.Document
		decodedValue, _ := url.QueryUnescape(*doc)
		if strings.Contains(decodedValue, "dynamo") || strings.Contains(decodedValue, "table") {
			fmt.Printf("=============================\n%s\n=============================\n", aws.StringValue(policy.Arn))
			fmt.Printf("%s\n", decodedValue)
		} else {
			fmt.Printf("No dynamo policy found in %s\n", aws.StringValue(policy.Arn))
		}
	}
}
