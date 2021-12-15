package tests

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/loop-functions-lib/sqrroot"
	"github.com/stretchr/testify/assert"
)

func Test_SqrtFullRun_ValidNumber(t *testing.T) {
	initialNumber, expectedResult := float64(4), float64(2)

	_, result := sqrroot.SqrtFullRun(initialNumber)

	assert.Equal(t, expectedResult, result)
}

func TestSqrtFullRun_InvalidNumber(t *testing.T) {
	initialNumber := float64(-4)

	assert.Panics(t, func() {
		_, result := sqrroot.SqrtFullRun(initialNumber)

		t.Log(result)
	})
}
