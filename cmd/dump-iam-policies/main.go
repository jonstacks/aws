package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

// Calculates the free ranges in a VPC that can be used to create new subnets.
func main() {
	models.Init(models.DefaultSession())

	policies, err := models.ListPolicies()
	utils.ExitErrorHandler(err)

	fmt.Printf("Got %d policies\n", len(policies))

	for _, policy := range policies {
		fmt.Printf("=============================\n%s\n=============================\n", aws.StringValue(policy.Arn))
		detail, err := models.DescribePolicy(*policy.Arn, *policy.DefaultVersionId)
		utils.ExitErrorHandler(err)

		doc := detail.PolicyVersion.Document
		decodedValue, _ := url.QueryUnescape(*doc)
		fmt.Printf("%s\n", decodedValue)
		time.Sleep(2 * time.Second)
	}
}
