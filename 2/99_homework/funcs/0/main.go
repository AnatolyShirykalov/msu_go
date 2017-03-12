package main

import (
	"fmt"
)

func mod(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x

}

// TODO: Реализовать вычисление Квадратного корня
func Sqrt(x float64) float64 {
	if x < 0 {
		// panic("not valid value")
		return 0
	}

	var e, xk float64 = 1E-3, 0
	x0 := float64(1)

	for {
		xk = (x0*x0 + x) / (2 * x0)
		if mod(xk*xk-x) < 2*e*mod(xk) {
			break
		}

		x0 = xk
	}
	return xk
}

func main() {
	fmt.Println(Sqrt(2))
}
