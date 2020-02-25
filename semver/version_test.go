package semver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVersion(t *testing.T) {
	t.Run("valid semver", func(t *testing.T) {
		tests := []struct {
			semver   string
			expected []int64
		}{
			{
				"1.0.0",
				[]int64{1, 0, 0},
			},
			{
				"0.1.0",
				[]int64{0, 1, 0},
			},
			{
				"0.0.1",
				[]int64{0, 0, 1},
			},
			{
				"0.0.0",
				[]int64{0, 0, 0},
			},
			{
				"1.0",
				[]int64{1, 0, 0},
			},
			{
				"0.1",
				[]int64{0, 1, 0},
			},
			{
				"0.0",
				[]int64{0, 0, 0},
			},
			{
				"0",
				[]int64{0, 0, 0},
			},
		}

		for i := range tests {
			version, err := NewVersion(tests[i].semver)
			require.NoError(t, err)
			require.NotNil(t, version)
			require.Equal(t, tests[i].expected[0], version.major)
			require.Equal(t, tests[i].expected[1], version.minor)
			require.Equal(t, tests[i].expected[2], version.patch)
		}
	})

	t.Run("invalid semver", func(t *testing.T) {
		tests := []struct {
			semver string
		}{
			{
				"a.b.c",
			},
			{
				"1.a.c",
			},
			{
				"a.1.1",
			},
		}

		for i := range tests {
			version, err := NewVersion(tests[i].semver)
			require.Error(t, err)
			require.Nil(t, version)
		}
	})
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		compareFrom string
		compareTo   string
		result      int
	}{
		{
			"1.0.0",
			"1.0.0",
			0,
		},
		{
			"1.0.0",
			"0.1.0",
			1,
		},
		{
			"1.0.0",
			"2.0.0",
			-1,
		},
		{
			"0.0.1",
			"0.1.0",
			-1,
		},
		{
			"0.1",
			"0.1.0",
			0,
		},
	}
	for i := range tests {
		compareFrom, err := NewVersion(tests[i].compareFrom)
		require.NoError(t, err)
		compareTo, err := NewVersion(tests[i].compareTo)
		require.NoError(t, err)

		result := compareFrom.Compare(compareTo)
		require.Equal(t, tests[i].result, result)
	}
}
