package core

import (
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
	y := []string{"a", "b", "c"}
	t.Run("X is in Y", func(t *testing.T) {
		x := "a"
		if !IfXinY(x, y) {
			t.Fatalf(`Test x ('%s') should exist `, x)
		}
	})
	t.Run("X is not in Y", func(t *testing.T) {
		x := "d"
		if IfXinY(x, y) {
			t.Fatalf(`Test x ('%s') should not exist `, x)
		}
	})
}

func TestExists(t *testing.T) {

}
