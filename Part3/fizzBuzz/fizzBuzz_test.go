package fizzBuzz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const fizz = "Fizz"
const buzz = "Buzz"
const fizzAndBuzz = "Fizz Buzz"

func Test_FizzBuzz(t *testing.T) {
	testCases := []struct {
		name           string
		input          int
		expectedResult string
	}{
		{
			name:           "Test FizzBuzz function with value 1",
			input:          1,
			expectedResult: "1",
		},
		{
			name:           "Test FizzBuzz function with value 3",
			input:          3,
			expectedResult: fizz,
		},
		{
			name:           "Test FizzBuzz function with value 5",
			input:          5,
			expectedResult: buzz,
		},
		{
			name:           "Test FizzBuzz function with value multiple of 3",
			input:          18,
			expectedResult: fizz,
		},
		{
			name:           "Test FizzBuzz function with value multiple of 5",
			input:          20,
			expectedResult: buzz,
		},
		{
			name:           "Test FizzBuzz function with value multiple of both, 3 and 5",
			input:          30,
			expectedResult: fizzAndBuzz,
		},
		{
			name:           "Test FizzBuzz function with value -6",
			input:          -6,
			expectedResult: fizz,
		},
		{
			name:           "Test FizzBuzz function with value -8",
			input:          -8,
			expectedResult: "-8",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedResult, FizzBuzz(tt.input))
		})
	}
}
