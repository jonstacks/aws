package views

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstanceFamily(t *testing.T) {
	x := rdsInstanceType("db.m5.large")
	assert.Equal(t, "db.m5", x.family())
	assert.Equal(t, "large", x.size())
	units, ok := x.normalizedUnits()
	assert.True(t, ok)
	assert.Equal(t, 8, units)
}
