package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
	"github.com/sirupsen/logrus"
)

const (
	KeyStatusInactive = "Inactive"
	KeyStatusActive   = "Active"
)

func main() {
	models.Init(models.DefaultSession())
	users, err := models.IAMUsers()
	utils.ExitErrorHandler(err)

	logrus.Infof("Found %d users in AWS account", len(users))
	//users = FilterConsoleUsers(users)
	//logrus.Infof("Filtered out console users. Considering %d users", len(users))

	userKeyMap := make(map[string][]*iam.AccessKeyMetadata)
	keyLastUsedMap := make(map[string]*iam.AccessKeyLastUsed)
	for _, user := range users {
		if user.UserName != nil {
			keys, err := models.IAMAccessKeysMeatadata(*user.UserName)
			if err != nil {
				utils.ExitErrorHandler(err)
			}
			keys = FilterInactiveKeys(keys)
			keys = FilterRecentlyCreatedKeys(keys, 10*24*time.Hour)

			userKeyMap[*user.UserName] = keys

			for _, key := range keys {
				resp, err := models.IAMAccessKeyLastUsed(*key.AccessKeyId)
				if err != nil {
					utils.ExitErrorHandler(err)
				}
				keyLastUsedMap[*key.AccessKeyId] = resp.AccessKeyLastUsed
			}
		}
	}

	v := views.IAMAccessKeysAudit{
		Users:    userKeyMap,
		LastUsed: keyLastUsedMap,
	}
	v.Print()
}

func FilterConsoleUsers(users []*iam.User) []*iam.User {
	filtered := make([]*iam.User, 0)
	for _, u := range users {
		if u.PasswordLastUsed == nil {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func FilterNonConsoleUsers(users []*iam.User) []*iam.User {
	filtered := make([]*iam.User, 0)
	for _, u := range users {
		if u.PasswordLastUsed != nil {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func FilterInactiveKeys(keys []*iam.AccessKeyMetadata) []*iam.AccessKeyMetadata {
	filtered := make([]*iam.AccessKeyMetadata, 0)
	for _, k := range keys {
		if k == nil {
			continue
		}

		if k.Status == nil {
			continue
		}

		if aws.StringValue(k.Status) == KeyStatusInactive {
			continue
		}

		filtered = append(filtered, k)
	}
	return filtered
}

func FilterRecentlyCreatedKeys(keys []*iam.AccessKeyMetadata, d time.Duration) []*iam.AccessKeyMetadata {
	filtered := make([]*iam.AccessKeyMetadata, 0)
	for _, k := range keys {
		if k == nil {
			continue
		}
		if k.CreateDate != nil && time.Since(*k.CreateDate) <= d {
			continue
		}
		filtered = append(filtered, k)
	}
	return filtered
}
