package tests

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/methods-types-structures/mts/math"
	"github.com/stretchr/testify/assert"
)

var target math.ErrNegativeSqrt = math.ErrNegativeSqrt(0)

func Test_Sqrt_WithValidNumber(t *testing.T) {
	number, root := 121, float64(11)

	value, err := math.Sqrt(number)

	assert.Nil(t, err)
	assert.EqualValues(t, root, value)
}

func Test_Sqrt_WithNegativeNumber(t *testing.T) {
	number := -4

	_, err := math.Sqrt(number)

	assert.ErrorIs(t, math.ErrNegativeSqrt(number), err)
}

func Test_Sqrt_WitNonNumberValue(t *testing.T) {
	testValue := "This is an invalid value"

	_, err := math.Sqrt(testValue)

	assert.ErrorIs(t, target, err)
}
