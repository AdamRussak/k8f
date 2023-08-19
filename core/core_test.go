package core

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnvVarOrSitIt(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "Test_CheckEnvVarOrSitIt_EnvVarSet",
			key:      "TEST",
			value:    "updated",
			expected: "SET",
		},
		{
			name:     "Test_CheckEnvVarOrSitIt_EnvVarNotSet",
			key:      "TEST",
			value:    "updated",
			expected: "updated",
		},
		{
			name:     "Test_CheckEnvVarOrSitIt_EnvVarSetAndNoChangeNeeded",
			key:      "TEST",
			value:    "SET",
			expected: "SET",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var origVal = "SET"
			if tc.expected == origVal {
				t.Setenv(tc.key, origVal)
			}
			CheckEnvVarOrSitIt(tc.key, tc.value)
			val, present := os.LookupEnv(tc.key)
			if !present || val != tc.expected {
				t.Fatalf(`Var(%s) = %s,should have been %s`, tc.key, val, tc.expected)
			}
		})
	}
}

func TestPrintOutStirng(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "EmptyArray",
			input:    []string{},
			expected: "",
		},
		{
			name:     "OneItemArray",
			input:    []string{"Test"},
			expected: " Test",
		},
		{
			name:     "TwoItemArray",
			input:    []string{"Test", "Test"},
			expected: " Test Test",
		},
		{
			name:     "ThreeItemArray",
			input:    []string{"Test", "Test", "Test"},
			expected: " Test Test Test",
		},
		{
			name:     "FourItemArray",
			input:    []string{"Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test",
		},
		{
			name:     "FiveItemArray",
			input:    []string{"Test", "Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test Test",
		},
		{
			name:     "SixItemArray",
			input:    []string{"Test", "Test", "Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test Test Test",
		},
		{
			name:     "SevenItemArray",
			input:    []string{"Test", "Test", "Test", "Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test Test Test Test",
		},
		{
			name:     "EightItemArray",
			input:    []string{"Test", "Test", "Test", "Test", "Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test Test Test Test Test",
		},
		{
			name:     "NineItemArray",
			input:    []string{"Test", "Test", "Test", "Test", "Test", "Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test Test Test Test Test Test",
		},
		{
			name:     "TenItemArray",
			input:    []string{"Test", "Test", "Test", "Test", "Test", "Test", "Test", "Test", "Test", "Test"},
			expected: " Test Test Test Test Test Test Test Test Test Test",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PrintOutStirng(tc.input)
			if result != tc.expected {
				t.Fatalf(`PrintOutStirng(%s) = %s,should have been %s`, tc.input, result, tc.expected)
			}
		})
	}
}

func TestIfXinY(t *testing.T) {
	testCases := []struct {
		name     string
		x        string
		y        []string
		expected bool
	}{
		{
			name:     "StringExistsInSlice",
			x:        "apple",
			y:        []string{"apple", "banana", "orange"},
			expected: true,
		},
		{
			name:     "StringDoesNotExistInSlice",
			x:        "grape",
			y:        []string{"apple", "banana", "orange"},
			expected: false,
		},
		{
			name:     "EmptySlice",
			x:        "apple",
			y:        []string{},
			expected: false,
		},
		{
			name:     "EmptyString",
			x:        "",
			y:        []string{"apple", "banana", "orange"},
			expected: false,
		},
		{
			name:     "EmptyStringAndEmptySlice",
			x:        "",
			y:        []string{},
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result bool = IfXinY(tc.x, tc.y)
			if result != tc.expected {
				t.Errorf("Expected %q in %q to be %v, but got %v", tc.x, tc.y, tc.expected, result)
			}
		})
	}
}

func TestExists(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "PathExists",
			path:     "existing_file.txt",
			expected: true,
		},
		{
			name:     "PathDoesNotExist",
			path:     "nonexistent_file.txt",
			expected: false,
		},
		{
			name:     "PermissionDenied",
			path:     "/root/somefile.txt",
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//create a temporary file
			if tc.expected {
				file, err := os.Create(tc.path)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					file.Close()
					os.Remove(tc.path)
				}()
			} else {
				os.Remove(tc.path)
			}
			var exists bool = Exists(tc.path)
			if exists != tc.expected {
				t.Errorf("Exists(%q) = %v; expected %v", tc.path, exists, tc.expected)
			}
		})
	}
}

func TestCreateDirectory(t *testing.T) {
	const conf string = "/config"
	testCases := []struct {
		name          string
		MainPath      string
		SecondaryPath string
		expected      bool
	}{
		{
			name:          "DirectoryDoesNotExist",
			MainPath:      "new_directory",
			SecondaryPath: conf,
			expected:      true,
		},
		{
			name:          "DirectoryExists",
			MainPath:      "existing_directory",
			SecondaryPath: conf,
			expected:      true,
		},
		{
			name:          "PermissionDenied",
			MainPath:      "/root",
			SecondaryPath: "/somefile.txt",
			expected:      false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//create a temporary file
			if tc.expected {
				err := os.MkdirAll(tc.MainPath+tc.SecondaryPath, 0777)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				os.RemoveAll(tc.MainPath + tc.SecondaryPath)
				os.RemoveAll(tc.MainPath)
			}

			var exists bool = Exists(tc.MainPath + tc.SecondaryPath)
			if exists != tc.expected {
				t.Errorf("Exists(%q) = %v; expected %v", tc.MainPath+tc.SecondaryPath, exists, tc.expected)
			} else if tc.expected {
				os.RemoveAll(tc.MainPath + tc.SecondaryPath)
				os.RemoveAll(tc.MainPath)
			}
		})
	}
}

func TestMergeINIFiles(t *testing.T) {
	testCases := []struct {
		name         string
		inputPaths   []string
		expectedData string
	}{
		{
			name:         "Creds first and then configuration",
			inputPaths:   []string{"../test/credentials", "../test/config"},
			expectedData: "[default]\naws_access_key_id = myKey\naws_secret_access_key = myAccessKey\n\n[account1]\naws_access_key_id = account1Key\naws_secret_access_key = account1AccessKey\n\n[account2]\naws_access_key_id = account2Key\naws_secret_access_key = account2AccessKey\n\n",
		},
		{
			name:         "Configuration and then Creds",
			inputPaths:   []string{"../test/config", "../test/credentials"},
			expectedData: "[default]\nregion = eu-west-1\n\n[profile account1]\nrole_arn = arn:aws:iam::123456789012:role/my-role\nsource_profile = default\n\n[profile account2]\nrole_arn = arn:aws:iam::125456389012:role/my-role\nsource_profile = default\n\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MergeINIFiles(tc.inputPaths)
			if err != nil {
				t.Fatalf("Error occurred: %v", err)
			}
			var actualDataBuf bytes.Buffer
			if _, err := actualDataBuf.ReadFrom(result); err != nil {
				t.Fatalf("Error reading merged data: %v", err)
			}
			actualData := actualDataBuf.String()
			if actualData != tc.expectedData {
				t.Errorf("Test case '%s' failed.\nExpected:\n%s\nGot:\n%s", tc.name, tc.expectedData, actualData)
			}
		})
	}
}

func TestCheckIfConfigExist(t *testing.T) {
	testCases := []struct {
		name         string
		inputString  string
		newKey       string
		expectedData bool
	}{
		{
			name: "section dose not exist",
			inputString: `[Section1]
			key1 = value1
			[Section2]
			key2 = value2
			`,
			newKey:       "newAAccount3",
			expectedData: false,
		},
		{
			name: "section exist",
			inputString: `[Section1]
			key1 = value1
			[Section2]
			key2 = value2
			`,
			newKey:       "section1",
			expectedData: true,
		},
		{
			name:         "No file exist",
			inputString:  ``,
			newKey:       "section1",
			expectedData: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			existingIni := bytes.NewBufferString(tc.inputString)
			testIfExist := checkIfConfigExist(tc.newKey, *existingIni)
			assert.True(t, testIfExist == tc.expectedData, "Unexpected result for "+tc.name)
		})
	}
}
