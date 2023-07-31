package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"k8f/core"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestReturnMinorDiff(t *testing.T) {
	testCases := []struct {
		name           string
		currentVersion []string
		latestVersion  []string
		expected       int
	}{
		{
			name:           "PerfectMatch",
			currentVersion: []string{"1", "27"},
			latestVersion:  []string{"1", "27"},
			expected:       0,
		}, {
			name:           "PerfectMatch-1",
			currentVersion: []string{"1", "26"},
			latestVersion:  []string{"1", "27"},
			expected:       1,
		}, {
			name:           "WarningMatch",
			currentVersion: []string{"1", "22"},
			latestVersion:  []string{"1", "27"},
			expected:       5,
		}, {
			name:           "CriticalMatch",
			currentVersion: []string{"1", "0"},
			latestVersion:  []string{"1", "27"},
			expected:       27,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test case: Perfect match, current version is the latest
			result := returnMinorDiff(tc.currentVersion, tc.latestVersion)
			assert.Equal(t, tc.expected, result, "Unexpected result for "+tc.name)

		})
	}
}
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
	setCurrent := "10.0.0"
	testCases := []struct {
		name     string
		latest   string
		current  string
		expected string
	}{
		{
			name:     "SameMajorVersion",
			latest:   setCurrent,
			current:  setCurrent,
			expected: "OK",
		},
		{
			name:     "MinorVersionDifference",
			latest:   "10.1.0",
			current:  setCurrent,
			expected: "OK",
		},
		{
			name:     "MinorVersionWarning",
			latest:   "10.10.0",
			current:  "10.8.0",
			expected: "Warning",
		},
		{
			name:     "MinorVersionCritical",
			latest:   "10.6.0",
			current:  setCurrent,
			expected: "Critical",
		},
		{
			name:     "DifferentMajorVersion",
			latest:   "11.0.0",
			current:  "10.1.0",
			expected: "Unknown",
		},
		// Add more test cases for other scenarios
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := microsoftSupportedVersion(tc.latest, tc.current)
			assert.Equal(t, tc.expected, result, "Unexpected result for "+tc.name)
		})
	}
}

func TestHowManyVersionsBack(t *testing.T) {
	versionsList := []string{"1.24", "1.23", "1.22", "1.21", "1.20", "1.27", "1.26", "1.25"}
	testCases := []struct {
		name           string
		currentVersion string
		expected       string
	}{
		{
			name:           "PerfectMatch",
			currentVersion: "1.27",
			expected:       "Perfect",
		}, {
			name:           "PerfectMatch-1",
			currentVersion: "1.26",
			expected:       "Perfect",
		}, {
			name:           "WarningMatch",
			currentVersion: "1.22",
			expected:       "Warning",
		}, {
			name:           "CriticalMatch",
			currentVersion: "0.5",
			expected:       "Critical",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test case: Perfect match, current version is the latest
			result := HowManyVersionsBack(versionsList, tc.currentVersion)
			assert.Equal(t, tc.expected, result, "Unexpected result for "+tc.name)

		})
	}
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

func TestGetBackupFileVersion(t *testing.T) {
	testCases := []struct {
		name     string
		current  int
		expected int
	}{
		{name: "NoBackups",
			current:  0,
			expected: 0},
		{name: "OneBackups",
			current:  1,
			expected: 1},
		{name: "FiveBackups",
			current:  4,
			expected: 5},
	}
	for _, tc := range testCases {
		directory := t.TempDir()
		for i := 0; i < tc.current; i++ {
			_, err := os.Create(directory + backupExtnesion + "." + fmt.Sprint(i))
			log.Debug(directory + backupExtnesion + "." + fmt.Sprint(i))
			core.FailOnError(err, filedtoCopyToTarget)
		}
		t.Run(tc.name, func(t *testing.T) {
			c := CommandOptions{Path: directory}
			// Test case: Perfect match, current version is the latest
			result := c.GetBackupFileVersion()
			assert.Equal(t, tc.expected, result, "Unexpected result for "+tc.name)
		})
	}
}
