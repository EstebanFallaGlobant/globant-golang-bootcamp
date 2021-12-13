package tests

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/methods-types-structures/mts/math"
)

var target math.ErrNegativeSqrt = math.ErrNegativeSqrt(0)

func Test_Sqrt_WithValidNumber(t *testing.T) {
	number, root := 121, 11

	//t.Logf("Original number: %d, Expected root value: %d", number, root)
	value, err := math.Sqrt(number)

	if err != nil {
		t.Fatal(err)
	}

	if value != float64(root) {
		t.Fatalf("The square root of %d is not %.1f", root, value)
	}
}

func Test_Sqrt_WithNegativeNumber(t *testing.T) {
	number := -4

	_, err := math.Sqrt(number)

	switch v := err.(type) {
	default:
		t.Fatalf("The error produced is not the expected. Got: %T, Expected: %T", v, target)
	case nil:
		t.Fatalf("There was no error. Expected: %T", v)
	case math.ErrNegativeSqrt:
		t.Logf("The error: %T, is of the expected type: %T", v, target)
	}
}

func Test_Sqrt_WitNonNumberValue(t *testing.T) {
	testValue := "This is an invalid value"

	_, err := math.Sqrt(testValue)

	switch v := err.(type) {
	default:
		t.Fatalf("The error produced is not the expected. Got: %T, Expected: %T", v, target)
	case nil:
		t.Fatalf("There was no error. Expected: %T", v)
	case math.ErrNegativeSqrt:
		t.Logf("The error: %T, is of the expected type: %T", v, target)
	}
}
