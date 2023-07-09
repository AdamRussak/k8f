package provider

import (
	"encoding/json"
	"io"
	"os"
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

func TestCheckIfStructInit(t *testing.T) {
	t.Run("FieldExists", func(t *testing.T) {
		// Test case: Field exists in the struct
		user := User{
			Exec:                  Exec{APIVersion: "1", Args: []string{"1", "2"}, Command: "noting", Env: "dev", ProvideClusterInfo: true},
			ClientCertificateData: "certData",
			ClientKeyData:         "clientKeyData",
			Token:                 "veryComplicatedToken",
		}

		result := checkIfStructInit(user, "exec")
		expected := true
		assert.Equal(t, expected, result, "Unexpected result for field existence")
	})

	t.Run("FieldOmitted", func(t *testing.T) {
		// Test case: Field is omitted in the struct
		user := User{}

		result := checkIfStructInit(user, "exec")
		expected := false
		assert.Equal(t, expected, result, "Unexpected result for field omission")
	})
	// Add more test cases for other scenarios
}

func TestCleanFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	// Write some content to the file
	content := []byte("Test content")
	_, err = tmpfile.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	// Close the file before cleaning it
	tmpfile.Close()

	// Create a CommandOptions instance with the temporary file path
	c := CommandOptions{Path: tmpfile.Name()}

	// Call the cleanFile method
	c.cleanFile()
	// Open the file again to check if it's empty
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	// Read the file content
	_, err = io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	// Check if the file is empty
	fileStat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	// Verify that the file size is 0 after cleaning
	if fileStat.Size() != 0 {
		t.Errorf("Expected file size after cleaning: 0, got: %d", fileStat.Size())
	}
}

func TestConfigCopy(t *testing.T) {
	t.Run("RegularFile", func(t *testing.T) {
		// Create a temporary file for testing
		tmpfile, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpfile.Name())

		// Create a CommandOptions instance with the temporary file path
		c := CommandOptions{Path: tmpfile.Name()}

		// Call the Configcopy method
		c.Configcopy()

		// Check if the backup file exists
		_, err = os.Stat(tmpfile.Name() + ".bak")
		if err != nil {
			t.Errorf("Expected backup file to exist, got error: %v", err)
		}
		os.RemoveAll(tmpfile.Name() + ".bak")
	})
}

func TestSplitAzIDAndGiveItem(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		input := "item1-item2-item3"
		separator := "-"
		index := 1
		expected := "item2"

		result := SplitAzIDAndGiveItem(input, separator, index)

		if result != expected {
			t.Errorf("Expected result: %s, got: %s", expected, result)
		}
	})
}
