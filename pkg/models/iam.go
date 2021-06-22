package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

var iamClient *iam.IAM

// IAMClient sets the client to be used by the models.
func IAMClient(client *iam.IAM) {
	iamClient = client
}

// IAM Users returns the list of IAM users
func IAMUsers() ([]*iam.User, error) {
	var err error
	input := &iam.ListUsersInput{}

	users := make([]*iam.User, 0)
	err = iamClient.ListUsersPages(input,
		func(page *iam.ListUsersOutput, lastPage bool) bool {
			users = append(users, page.Users...)
			return !lastPage
		})

	return users, err
}

// IAMAccessKeysMetadata returns all AccessKeyMetadata
func IAMAccessKeysMeatadata(username string) ([]*iam.AccessKeyMetadata, error) {
	var err error
	input := &iam.ListAccessKeysInput{
		UserName: aws.String(username),
	}

	keys := make([]*iam.AccessKeyMetadata, 0)
	err = iamClient.ListAccessKeysPages(input,
		func(page *iam.ListAccessKeysOutput, lastPage bool) bool {
			keys = append(keys, page.AccessKeyMetadata...)
			return !lastPage
		})
	return keys, err
}

func IAMAccessKeyLastUsed(accessKeyID string) (*iam.GetAccessKeyLastUsedOutput, error) {
	input := &iam.GetAccessKeyLastUsedInput{
		AccessKeyId: aws.String(accessKeyID),
	}
	return iamClient.GetAccessKeyLastUsed(input)
}
