package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testErrorMessage = "Unexpected result for "

func TestAwsArgs(t *testing.T) {
	const testRegion = "eu-west-1"
	const testRegionFlag = "--region"
	const testGetTokenFlag = "get-token"
	const testClusterArn = "arn:aws:eks:eu-west-1:123456789:cluster/testCluster"
	const testClusterNameFlag = "--cluster-name"
	const testClusterName = "testCluster"
	// the testcase should have the CommandOptions struct, the region string, the clustername string and the arn string and expected array
	testCases := []struct {
		name        string
		Command     CommandOptions
		region      string
		clusterName string
		arn         string
		expected    []string
	}{
		{
			name:        "Test_AwsArgs_only_role",
			Command:     CommandOptions{AwsRoleString: "testRole", AwsAuth: false},
			region:      testRegion,
			clusterName: testClusterName,
			arn:         testClusterArn,
			expected:    []string{testRegionFlag, testRegion, "eks", testGetTokenFlag, testClusterNameFlag, testClusterName, "--role-arn", "arn:aws:iam::123456789:role/testRole"},
		},
		{
			name:        "Test_AwsArgs_only_role_with_aws_auth",
			Command:     CommandOptions{AwsRoleString: "testRole", AwsAuth: true},
			region:      testRegion,
			clusterName: testClusterName,
			arn:         testClusterArn,
			expected:    []string{"token", "-i", testClusterName, "--role-arn", "arn:aws:iam::123456789:role/testRole"},
		},
		{
			name:        "Test_AwsArgs_without_role",
			Command:     CommandOptions{AwsRoleString: "", AwsAuth: false},
			region:      testRegion,
			clusterName: testClusterName,
			arn:         testClusterArn,
			expected:    []string{testRegionFlag, testRegion, "eks", testGetTokenFlag, testClusterNameFlag, testClusterName},
		},
		{
			name:        "Test_AwsArgs_without_role_with_aws_auth",
			Command:     CommandOptions{AwsRoleString: "", AwsAuth: true},
			region:      testRegion,
			clusterName: testClusterName,
			arn:         testClusterArn,
			expected:    []string{testRegionFlag, testRegion, "eks", testGetTokenFlag, testClusterNameFlag, testClusterName},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var c CommandOptions = CommandOptions{AwsRoleString: tc.Command.AwsRoleString, AwsAuth: tc.Command.AwsAuth}
			args := c.AwsArgs(tc.region, tc.clusterName, tc.arn)
			assert.Equal(t, tc.expected, args, testErrorMessage+tc.name)
		})
	}
}

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
			assert.Equal(t, tc.expected, result, testErrorMessage+tc.name)
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
			assert.Equal(t, tc.expected, result, testErrorMessage+tc.name)
		})
	}
}
