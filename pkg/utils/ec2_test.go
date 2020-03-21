package utils

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

var testInstance1 = &ec2.Instance{
	Tags: []*ec2.Tag{
		&ec2.Tag{
			Key:   aws.String("Name"),
			Value: aws.String("my-test-instance"),
		},
	},
}

var testInstance2 = &ec2.Instance{
	Tags: []*ec2.Tag{
		&ec2.Tag{
			Key:   aws.String("Expires"),
			Value: aws.String("2012-03-20"),
		},
	},
}

func TestGetInstanceName(t *testing.T) {
	assert.Equal(t, "my-test-instance", GetInstanceName(testInstance1))
	assert.Equal(t, "", GetInstanceName(testInstance2))
}

func TestGetTagValue(t *testing.T) {
	assert.Equal(t, "my-test-instance", GetTagValue(testInstance1.Tags, "Name"))
	assert.Equal(t, "", GetTagValue(testInstance1.Tags, "MissingTag"))

	assert.Equal(t, "", GetTagValue(testInstance2.Tags, "Name"))
	assert.Equal(t, "2012-03-20", GetTagValue(testInstance2.Tags, "Expires"))
}

func TestSubnetSize(t *testing.T) {
	size, err := SubnetSize("10.1.0.0/24")
	assert.Nil(t, err)
	assert.Equal(t, 256, size)

	size, err = SubnetSize("IDontParse")
	assert.NotNil(t, err)
	assert.Equal(t, 0, size)

	size, err = SubnetSize("10.1.0.0/0")
	assert.Nil(t, err)
	assert.Equal(t, 4294967296, size)

	size, err = SubnetSize("10.1.0.0/32")
	assert.Nil(t, err)
	assert.Equal(t, 1, size)
}

func TestStringSliceContains(t *testing.T) {
	testCases := []struct {
		slice          []string
		searchTerm     string
		expectedResult bool
	}{
		{[]string{"a", "b", "c", "d"}, "c", true},
		{[]string{"a", "b", "c", "d"}, "a", true},
		{[]string{}, "a", false},
		{[]string{}, "", false},
		{[]string{""}, "c", false},
		{[]string{"a", "b"}, "c", false},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedResult, StringSliceContains(tc.slice, tc.searchTerm))
	}
}
