package main

import (
	"fmt"

	"golang.org/x/tour/pic"
)

var x, y int = 8, 8

func main() {
	pic.Show(Pic)
	matrix := Pic(x, y)

	for i := 0; i < y; i++ {
		line := ""
		for a := 0; a < x; a++ {
			line += fmt.Sprintf("%d", matrix[i][a])
		}
		fmt.Println(line)
	}
}

func Pic(dx, dy int) [][]uint8 {
	result := make([][]uint8, dy)

	for y := 0; y < dy; y++ {
		temp := make([]uint8, dx)

		for x := 0; x < dx; x++ {
			temp[x] = uint8(((x + 1) % (y + 1)) * 20)
		}

		result[y] = temp
	}

	return result
}
