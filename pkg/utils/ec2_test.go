package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
