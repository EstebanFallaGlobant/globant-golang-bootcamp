package main

import (
	"fmt"
	"os"
	"strconv"

	sqrroot "github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/loop-functions-lib/sqr-root"
)

func main() {
	var number float64 = 4

	// Checks if the program was executed sending a value as parameter
	if len(os.Args) > 1 {
		temp, err := strconv.Atoi(os.Args[1])

		if err == nil {
			number = float64(temp)
		} else {
			fmt.Printf("%t\nUsing default value: %d\n", err, int(number))
		}
	}

	fmt.Printf("Printing the result of the square root of '%g' in %d steps, also the result of every step\n", number, sqrroot.TotalSteps)

	steps, sqrt := sqrroot.SqrtFullRun(number)

	fmt.Printf("The square root of %g is %.2f\nThe steps for this result are: \n", number, sqrt)
	for i, step := range steps {
		fmt.Printf("%d: %g,\n", i+1, step)
	}
	fmt.Println()

	fmt.Printf("Printing the result of the square root of '%g' in %d steps, or less\n", number, sqrroot.TotalSteps)

	totalSteps := 0
	sqrt, totalSteps = sqrroot.Sqrt(number)

	fmt.Printf("The square root of %d is %.2f, the result took %d steps\n", int(number), sqrt, totalSteps)
}
