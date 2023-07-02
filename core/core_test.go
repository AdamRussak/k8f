package core

import (
	"fmt"
	"os"
	"testing"
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
	t.Run("MergeMultipleINIFiles", func(t *testing.T) {
		// Test case: Merge multiple INI files
		inputPaths := []string{"../test/file1.ini", "../test/file2.ini", "../test/file0.ini"}

		reader, err := MergeINIFiles(inputPaths)

		if err != nil {
			t.Errorf("Error occurred while merging INI files: %v", err)
		}

		expectedOutput := `[DEFAULT]
Key1 = Value1
Key2 = Value2

[DEFAULT]
Key3 = Value3

[DEFAULT]
Key6 = Value6
`
		buffer := make([]byte, len(expectedOutput))
		_, err = reader.Read(buffer)
		if err != nil {
			t.Errorf("Error occurred while reading merged INI files: %v", err)
		}
		actualOutput := string(buffer)
		if actualOutput != expectedOutput {
			t.Errorf("Merged INI files do not match the expected output.\nExpected:\n%s\n\nActual:\n%s", expectedOutput, actualOutput)
		}
	})

	t.Run("MergeNoINIFiles", func(t *testing.T) {
		// Test case: No INI files to merge
		inputPaths := []string{}
		reader, _ := MergeINIFiles(inputPaths)
		if reader.Len() != 0 {
			t.Errorf("Merged INI files do not match the expected output.\nExpected:\n%s\n\nActual:\n%s", "", fmt.Sprint(reader.Len()))
		}
	})
}
