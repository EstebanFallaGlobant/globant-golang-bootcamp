package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var target ErrNegativeSqrt = ErrNegativeSqrt(0)

func Test_Sqrt_WithValidNumber(t *testing.T) {
	testCases := []struct {
		input          interface{}
		expectedResult float64
		name           string
	}{
		{
			name:           "Testing with int",
			input:          int(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with int8",
			input:          int8(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with int16",
			input:          int16(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with int32",
			input:          int32(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with int64",
			input:          int64(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with uint",
			input:          uint(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with uint8",
			input:          uint8(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with uint16",
			input:          uint16(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with uint32",
			input:          uint32(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with uint64",
			input:          uint64(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with float32",
			input:          float32(16),
			expectedResult: float64(4),
		},
		{
			name:           "Testing with float64",
			input:          float64(16),
			expectedResult: float64(4),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := Sqrt(tt.input)
			assert.EqualValues(t, tt.expectedResult, value)
		})
	}
}

func Test_Sqrt_WithNegativeNumber(t *testing.T) {
	number := -4

	_, err := Sqrt(number)

	assert.ErrorIs(t, ErrNegativeSqrt(number), err)
}

func Test_Sqrt_WitNonNumberValue(t *testing.T) {
	testValue := "This is an invalid value"

	_, err := Sqrt(testValue)

	assert.ErrorIs(t, target, err)
}

func Test_Abs_PositiveNumber(t *testing.T) {
	testValue, expectedValue := float64(100), float64(100)

	val := abs(testValue)

	assert.EqualValues(t, expectedValue, val)
}

func Test_Abs_NegativeNumber(t *testing.T) {
	testValue, expectedValue := float64(-100), float64(100)

	val := abs(testValue)

	assert.EqualValues(t, expectedValue, val)
}
