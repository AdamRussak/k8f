package provider

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestEvaluateVersion(t *testing.T) {
	t.Run("EmptyList", func(t *testing.T) {
		// Test case: Empty list
		list := []string{}
		result := evaluateVersion(list)
		expected := ""
		assert.Equal(t, expected, result, "Unexpected result for empty list")
	})

	t.Run("SingleVersion", func(t *testing.T) {
		// Test case: Single version in the list
		list := []string{"1.2.3"}
		result := evaluateVersion(list)
		expected := "1.2.3"
		assert.Equal(t, expected, result, "Unexpected result for single version")
	})

	t.Run("MultipleVersions", func(t *testing.T) {
		// Test case: Multiple versions in the list
		list := []string{"1.0.0", "1.2.3", "2.0.0", "0.1.0"}
		result := evaluateVersion(list)
		expected := "2.0.0"
		assert.Equal(t, expected, result, "Unexpected result for multiple versions")
	})
	// Add more test cases for other scenarios
}

func TestMicrosoftSupportedVersion(t *testing.T) {
	t.Run("SameMajorVersion", func(t *testing.T) {
		// Test case: Same major version
		latest := "10.0.0"
		current := "10.0.0"
		result := microsoftSupportedVersion(latest, current)
		expected := "OK"
		assert.Equal(t, expected, result, "Unexpected result for same major version")
	})

	t.Run("MinorVersionDifference", func(t *testing.T) {
		// Test case: Minor version difference within range
		latest := "10.0.0"
		current := "10.1.0"
		result := microsoftSupportedVersion(latest, current)
		expected := "OK"
		assert.Equal(t, expected, result, "Unexpected result for minor version difference within range")
	})

	t.Run("MinorVersionWarning", func(t *testing.T) {
		// Test case: Minor version difference exceeding range, but within warning range
		latest := "10.10.0"
		current := "10.8.0"
		result := microsoftSupportedVersion(latest, current)
		expected := "Warning"
		assert.Equal(t, expected, result, "Unexpected result for minor version warning")
	})

	t.Run("MinorVersionCritical", func(t *testing.T) {
		// Test case: Minor version difference exceeding range, critical status
		latest := "10.6.0"
		current := "10.0.0"
		result := microsoftSupportedVersion(latest, current)
		expected := "Critical"
		assert.Equal(t, expected, result, "Unexpected result for minor version critical")
	})

	t.Run("DifferentMajorVersion", func(t *testing.T) {
		// Test case: Different major versions
		latest := "11.0.0"
		current := "10.1.0"
		result := microsoftSupportedVersion(latest, current)
		expected := "Unknown"
		assert.Equal(t, expected, result, "Unexpected result for different major versions")
	})
	// Add more test cases for other scenarios
}

func TestHowManyVersionsBack(t *testing.T) {
	versionsList := []string{"v1.0", "v1.1", "v1.2", "v1.3", "v1.4", "v1.5"}

	t.Run("PerfectMatch", func(t *testing.T) {
		// Test case: Perfect match, current version is the latest
		currentVersion := "v1.0"
		result := HowManyVersionsBack(versionsList, currentVersion)
		expected := "Perfect"
		assert.Equal(t, expected, result, "Unexpected result for perfect match")
	})

	t.Run("OKMatch", func(t *testing.T) {
		// Test case: OK match, current version is within 3 versions from the latest
		currentVersion := "v1.3"
		result := HowManyVersionsBack(versionsList, currentVersion)
		expected := "OK"
		assert.Equal(t, expected, result, "Unexpected result for OK match")
	})

	t.Run("WarningMatch", func(t *testing.T) {
		// Test case: Warning match, current version is more than 3 versions back
		currentVersion := "v1.5"
		result := HowManyVersionsBack(versionsList, currentVersion)
		expected := "Warning"
		assert.Equal(t, expected, result, "Unexpected result for warning match")
	})

	t.Run("CriticalMatch", func(t *testing.T) {
		// Test case: Critical match, current version is not found in the list
		currentVersion := "v0.5"
		result := HowManyVersionsBack(versionsList, currentVersion)
		expected := "Critical"
		assert.Equal(t, expected, result, "Unexpected result for critical match")
	})
	// Add more test cases for other scenarios
}

func TestPrintoutResults(t *testing.T) {
	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John Doe",
		Age:  30,
	}

	t.Run("JSONOutput", func(t *testing.T) {
		// Test case: JSON output
		outputType := "json"
		result := PrintoutResults(data, outputType)
		expected, _ := json.Marshal(data)
		assert.Equal(t, string(expected), result, "Unexpected result for JSON output")
	})

	t.Run("YAMLOutput", func(t *testing.T) {
		// Test case: YAML output
		outputType := "yaml"
		result := PrintoutResults(data, outputType)
		expected, _ := yaml.Marshal(data)
		assert.Equal(t, string(expected), result, "Unexpected result for YAML output")
	})

	t.Run("UnsupportedOutput", func(t *testing.T) {
		// Test case: Unsupported output type
		outputType := "csv"
		result := PrintoutResults(data, outputType)
		expected := "Requested Output is not supported"
		assert.Equal(t, expected, result, "Unexpected result for unsupported output type")
	})
	// Add more test cases for other scenarios
}

func TestCountTotal(t *testing.T) {
	accounts := []Account{
		{TotalCount: 10},
		{TotalCount: 20},
		{TotalCount: 30},
	}

	t.Run("MultipleAccounts", func(t *testing.T) {
		// Test case: Multiple accounts
		result := countTotal(accounts)
		expected := 60
		assert.Equal(t, expected, result, "Unexpected total count for multiple accounts")
	})

	t.Run("EmptyAccounts", func(t *testing.T) {
		// Test case: Empty accounts
		result := countTotal([]Account{})
		expected := 0
		assert.Equal(t, expected, result, "Unexpected total count for empty accounts")
	})
	// Add more test cases for other scenarios
}
