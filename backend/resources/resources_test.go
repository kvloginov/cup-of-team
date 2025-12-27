package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadHolidaysRu12(t *testing.T) {
	holidays, err := LoadHolidaysRu12()

	require.NoError(t, err, "LoadHolidaysRu12 should not return an error")
	require.NotNil(t, holidays, "Holidays should not be nil")
	assert.NotEmpty(t, holidays, "Holidays should not be empty")
}

func TestLoadNamedaysRu(t *testing.T) {
	namedays, err := LoadNamedaysRu()

	require.NoError(t, err, "LoadNamedaysRu should not return an error")
	require.NotNil(t, namedays, "Namedays should not be nil")
	assert.NotEmpty(t, namedays, "Namedays should not be empty")
}
