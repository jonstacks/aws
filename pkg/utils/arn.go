package utils

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

func RemoveARN(s *string) string {
	deRef := aws.StringValue(s)
	parts := strings.Split(deRef, `:::`)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
