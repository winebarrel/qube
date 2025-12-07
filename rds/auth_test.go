package rds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube/rds"
)

func TestResolveCNAME(t *testing.T) {
	t.Run("rds.amazonaws.com suffix", func(t *testing.T) {
		host, err := rds.ResolveCNAME("mydb.123456789012.us-east-1.rds.amazonaws.com")
		require.NoError(t, err)
		assert.Equal(t, "mydb.123456789012.us-east-1.rds.amazonaws.com", host)
	})

	t.Run("resolve CNAME", func(t *testing.T) {
		host, err := rds.ResolveCNAME("www.google.com")

		if err != nil {
			t.Skipf("skipping test due to network error: %v", err)
		}

		require.NoError(t, err)
		assert.NotEmpty(t, host)
	})

	t.Run("lookup error", func(t *testing.T) {
		_, err := rds.ResolveCNAME("non-existent-domain.invalid.")
		require.Error(t, err)
	})
}
