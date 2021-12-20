package fizzBuzz

import (
	"fmt"
)

func FizzBuzz(number int) string {
	var result string

	mult3, mult5 := number%3 == 0, number%5 == 0

	if mult3 && mult5 {
		result = "Fizz Buzz"
	} else if mult3 {
		result = "Fizz"
	} else if mult5 {
		result = "Buzz"
	} else {
		result = fmt.Sprint(number)
	}

	return result
}
