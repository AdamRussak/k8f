package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetClusterName(t *testing.T) {
	testCases := []struct {
		name           string
		AwsClusterName bool
		ClusterName    string
		expected       string
	}{
		{
			name:           "Full_Cluster_output",
			AwsClusterName: false,
			ClusterName:    "arn:aws:eks:ap-southeast-2:234432434234:cluster/testCLuster01",
			expected:       "234432434234:ap-southeast-2:testCLuster01",
		},
		{
			name:           "Short_ARN_output",
			AwsClusterName: true,
			ClusterName:    "arn:aws:eks:ap-southeast-2:234432434234:cluster/testCLuster02",
			expected:       "ap-southeast-2:testCLuster02",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := CommandOptions{AwsClusterName: tc.AwsClusterName}
			result := c.SetClusterName(&tc.ClusterName)
			assert.Equal(t, tc.expected, result, "Unexpected result for "+tc.name)
		})
	}
}

func TestRemoveString(t *testing.T) {
	testCases := []struct {
		name         string
		fullString   string
		whatToRemove string
		expected     string
	}{
		{
			name:         "Test_proifle_remove",
			fullString:   "profile test",
			whatToRemove: "profile",
			expected:     "test",
		}, {
			name:         "Test_default_without_profile",
			fullString:   "default",
			whatToRemove: "profile",
			expected:     "default",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeString(tc.whatToRemove, tc.fullString)
			assert.Equal(t, tc.expected, result, "Unexpected result for "+tc.name)
		})
	}
}
