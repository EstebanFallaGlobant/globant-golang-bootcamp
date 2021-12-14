package math

import "fmt"

const iterations = 100

//
func Sqrt(value interface{}) (float64, error) {
	var result, temp, num float64

	switch v := value.(type) {
	case int:
		num = float64(v)
	case int8:
		num = float64(v)
	case int16:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case uint:
		num = float64(v)
	case uint8:
		num = float64(v)
	case uint16:
		num = float64(v)
	case uint32:
		num = float64(v)
	case uint64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		goto RETURNERROR
	}

	if num < 0 {
		goto RETURNERROR
	}

	result = float64(1)

	for i := 0; i < iterations; i++ {
		result -= newtonAprox(result, num)

		if v, _ := Abs(temp - result); v == float64(0) {
			break
		}
		temp = result
	}

	return result, nil

RETURNERROR:
	return 0, ErrNegativeSqrt(num)
}

func Abs(value interface{}) (float64, error) {
	var num float64

	switch v := value.(type) {
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case float32:
		num = float64(v)
	case float64:
		num = v

	default:
		return 0, fmt.Errorf("value %t is not a valid number", value)
	}

	if num < 0 {
		return -num, nil
	}

	return num, nil
}

//Returns a number close to the square root of the value passed as parameter, using a previous aproximation and the powered value
func newtonAprox(aprox float64, value float64) float64 {

	return (aprox*aprox - value) / (2 * aprox)
}
