package main

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
	"github.com/jonstacks/aws/pkg/views"
	"github.com/sirupsen/logrus"
)

func main() {
	models.Init(models.DefaultSession())
	users, err := models.IAMUsers()
	utils.ExitErrorHandler(err)

	logrus.Infof("Found %d users in AWS account", len(users))
	users = FilterConsoleUsers(users)
	logrus.Infof("Filtered out console users. Considering %d users", len(users))

	userKeyMap := make(map[string][]*iam.AccessKeyMetadata)
	keyLastUsedMap := make(map[string]*iam.AccessKeyLastUsed)
	for _, user := range users {
		if user.UserName != nil {
			keys, err := models.IAMAccessKeysMeatadata(*user.UserName)
			if err != nil {
				utils.ExitErrorHandler(err)
			}
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
