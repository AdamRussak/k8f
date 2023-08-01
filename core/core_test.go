package core

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnvVarOrSitIt(t *testing.T) {
	var key = "TEST"
	var value = "updated"
	t.Run("Var set and no change needed", func(t *testing.T) {
		var origVal = "SET"
		t.Setenv(key, origVal)
		CheckEnvVarOrSitIt(key, value)
		val, present := os.LookupEnv(key)
		if !present || val != "SET" {
			t.Fatalf(`Var(%s) = %s,should have been %s`, key, value, origVal)
		}
	})
	t.Run("Var Not set updated", func(t *testing.T) {
		CheckEnvVarOrSitIt(key, value)
		val, present := os.LookupEnv(key)
		if !present || val != value {
			t.Fatalf(`Var(%s) = %s,should have been %s`, key, value, value)
		}
	})
}

func TestPrintOutStirng(t *testing.T) {
	var shouldBe = " T e s t"

	t.Run("Array prints out correct string", func(t *testing.T) {
		var corerctArray = []string{"T", "e", "s", "t"}
		testString := PrintOutStirng(corerctArray)
		if testString != shouldBe {
			t.Fatalf(`Test String should be '%s' but was actuly '%s'`, shouldBe, testString)
		}
	})
}

func TestIfXinY(t *testing.T) {
	t.Run("StringExistsInSlice", func(t *testing.T) {
		// Test case: String exists in the slice
		x := "apple"
		y := []string{"apple", "banana", "orange"}

		result := IfXinY(x, y)

		if !result {
			t.Errorf("Expected %q to exist in the slice, but it doesn't", x)
		}
	})

	t.Run("StringDoesNotExistInSlice", func(t *testing.T) {
		// Test case: String does not exist in the slice
		x := "grape"
		y := []string{"apple", "banana", "orange"}

		result := IfXinY(x, y)

		if result {
			t.Errorf("Expected %q not to exist in the slice, but it does", x)
		}
	})

	t.Run("EmptySlice", func(t *testing.T) {
		// Test case: Empty slice
		x := "apple"
		y := []string{}

		result := IfXinY(x, y)

		if result {
			t.Errorf("Expected %q not to exist in the empty slice, but it does", x)
		}
	})
}

func TestExists(t *testing.T) {
	t.Run("PathExists", func(t *testing.T) {
		// Test case: Path exists
		path := "existing_file.txt"

		// Create a temporary file
		file, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			file.Close()
			os.Remove(path)
		}()

		exists := Exists(path)
		if !exists {
			t.Errorf("Exists(%q) = false; expected true", path)
		}
	})

	t.Run("PathDoesNotExist", func(t *testing.T) {
		// Test case: Path does not exist
		path := "nonexistent_file.txt"

		// Remove the file if it exists (cleanup from previous runs)
		os.Remove(path)

		exists := Exists(path)
		if exists {
			t.Errorf("Exists(%q) = true; expected false", path)
		}
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		// Test case: Permission denied to access path
		path := "/root/somefile.txt"

		exists := Exists(path)
		if exists {
			t.Errorf("Exists(%q) = true; expected false", path)
		}
	})
}

func TestCreateDirectory(t *testing.T) {
	t.Run("DirectoryDoesNotExist", func(t *testing.T) {
		// Test case: Directory already exists
		dir := "new_directory/"
		testPath := dir + "config"
		os.RemoveAll(dir)
		// Create a temporary directory
		defer func() {
			os.RemoveAll(dir)
		}()
		CreateDirectory(testPath)
		// Verify that the directory still exists
		_, err := os.Stat(dir)
		if err != nil {
			t.Errorf("Directory %q should exist, but got error: %v", dir, err)
		}
	})

	t.Run("DirectoryExists", func(t *testing.T) {
		// Test case: Directory already exists
		dir := "new_directory/"
		testPath := dir + "config"
		os.RemoveAll(dir)
		// Create a temporary directory
		err := os.Mkdir(dir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			os.RemoveAll(dir)
		}()
		CreateDirectory(testPath)
		// Verify that the directory still exists
		_, err = os.Stat(dir)
		if err != nil {
			t.Errorf("Directory %q should exist, but got error: %v", dir, err)
		}
	})
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
