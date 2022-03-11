package main

import (
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

type IamPolicyWithDocument struct {
	Policy   *iam.Policy
	Document string
}

// Dump IAM policies to stdout or redirect to file with shell redirect
func main() {
	models.Init(models.DefaultSession())

	policies, err := models.ListPolicies()
	utils.ExitErrorHandler(err)

	log.Printf("Got %d policies\n", len(policies))

	policyChan := make(chan *iam.Policy, 100)
	policyWg := sync.WaitGroup{}
	respChan := make(chan *IamPolicyWithDocument, 100)

	// Populate the policy channel the workers will pull from.
	go func() {
		for _, policy := range policies {
			policyChan <- policy
		}
		close(policyChan)
	}()

	// Start workers
	for i := 1; i <= 4; i++ {
		policyWg.Add(1)
		go func(workerNum int) {
			defer policyWg.Done()
			for policy := range policyChan {
				var detail *iam.GetPolicyVersionOutput
				for i := 0; i < 3; i++ {
					detail, err = models.DescribePolicy(*policy.Arn, *policy.DefaultVersionId)
					if err == nil {
						break
					}
					time.Sleep(time.Second * 5)
				}
				doc := detail.PolicyVersion.Document
				decodedDoc, err := url.QueryUnescape(*doc)
				if err != nil {
					log.Printf("[Worker %d] Error decoding document: %s\n", workerNum, err)
					continue // If we can't decode it, just log the message and move on for now
				}

				respChan <- &IamPolicyWithDocument{policy, decodedDoc}
			}
		}(i)
	}

	go func() {
		policyWg.Wait()
		close(respChan)
	}()

	// Process responses
	for resp := range respChan {
		if strings.Contains(resp.Document, "dynamo") || strings.Contains(resp.Document, "table") {
			log.Printf("=============================\n%s\n=============================\n", aws.StringValue(resp.Policy.Arn))
			log.Printf("%s\n", resp.Document)
		} else {
			log.Printf("No dynamo policy found in %s\n", aws.StringValue(resp.Policy.Arn))
		}
	}
}
