package sqrroot

const TotalSteps = 100
const ErrorMargin = 0
const initialValue = 1

//Using the Newton's aproximation method returns the square root of a number, also returns all the steps of the aproximation
func SqrtFullRun(value float64) ([]float64, float64) {
	var steps []float64
	var z float64 = float64(initialValue)

	for i := 0; i < TotalSteps; i++ {
		z -= newtonAprox(z, value)
		steps = append(steps, z)
	}

	return steps, z
}

//Using the Newton's aproximation returns the square root of a number, also returns the number of steps taken to find the result
func Sqrt(value float64) (float64, int) {
	var z, temp float64 = float64(initialValue), 0
	i := 0
	for i < TotalSteps {
		i++

		z -= newtonAprox(z, value)

		if AbsVal(temp-z) == ErrorMargin {
			break
		}
		temp = z
	}

	return z, i
}

func AbsVal(value float64) float64 {
	if value < 0 {
		return value * -1
	} else {
		return value
	}
}

//Returns a number close to the square root of the value passed as parameter, using a previous aproximation and the powered value
func newtonAprox(aprox float64, value float64) float64 {

	return (aprox*aprox - float64(value)) / (2 * aprox)
}
