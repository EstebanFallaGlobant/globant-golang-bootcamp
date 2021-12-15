package tests

import (
	"fmt"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part3/fizzBuzz"
	"github.com/stretchr/testify/suite"
)

type FizzBuzzTestSuite struct {
	suite.Suite
	initialValue int
}

const fizz = "Fizz"
const buzz = "Buzz"
const fizzAndBuzz = "Fizz Buzz"

func Test_FizzBuzz(t *testing.T) {
	suite.Run(t, new(FizzBuzzTestSuite))
}

func (suite *FizzBuzzTestSuite) SetupTest() {
	suite.initialValue = 1
}

func (suite *FizzBuzzTestSuite) Test_FizzBuzz_Value1() {
	suite.Equal(fmt.Sprint(1), fizzBuzz.FizzBuzz(suite.initialValue))
}

func (suite *FizzBuzzTestSuite) Test_FizzBuzz_Value3() {
	suite.initialValue = 3
	suite.Equal(fizz, fizzBuzz.FizzBuzz(suite.initialValue))
}

func (suite *FizzBuzzTestSuite) Test_FizzBuzz_Value5() {
	suite.initialValue = 5
	suite.Equal(buzz, fizzBuzz.FizzBuzz(suite.initialValue))
}

func (suite *FizzBuzzTestSuite) Test_FizzBuzz_ValueMultipleOf3() {
	suite.initialValue = 18
	suite.Equal(fizz, fizzBuzz.FizzBuzz(suite.initialValue))
}

func (suite *FizzBuzzTestSuite) Test_FizzBuzz_ValueMultipleOf5() {
	suite.initialValue = 35
	suite.Equal(buzz, fizzBuzz.FizzBuzz(suite.initialValue))
}

func (suite *FizzBuzzTestSuite) Test_FizzBuzz_ValueMultipleOf3And5() {
	suite.initialValue = 30
	suite.Equal(fizzAndBuzz, fizzBuzz.FizzBuzz(suite.initialValue))
}
