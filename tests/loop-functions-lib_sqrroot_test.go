package tests

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/loop-functions-lib/sqrroot"
)

func Test_SqrtFullRun_ValidNumber(t *testing.T) {
	initialNumber, expectedResult := float64(4), float64(2)

	_, result := sqrroot.SqrtFullRun(initialNumber)

	if result != expectedResult {
		t.Fatalf("Expected: %g. Got: %g", expectedResult, result)
	}
}

func TestSqrtFullRun_InvalidNumber(t *testing.T) {
	initialNumber := float64(-4)

	defer func() {
		if r := recover().(error); r == nil {
			t.Fatal()
		}
	}()

	_, result := sqrroot.SqrtFullRun(initialNumber)

	if result != 0 {
		t.Fatalf("The function din't panic. Got: %g. Expected: panic", result)
	}
}
