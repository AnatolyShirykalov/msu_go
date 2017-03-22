package main

import (
	"fmt"
)

func simpleType(i interface{}) string {
	switch i.(type) {

	case bool:
		return "bool"
	case int:
		return "int"
	case uint:
		return "uint"
	case int8:
		return "int8"
	case int32:
		return "int32"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case string:
		return "string"

	case []bool:
		return "[]bool"
	case []int:
		return "[]int"
	case []string:
		return "[]string"
	case []float32:
		return "[]float"
	default:
		return "none"
	}
}
func mapType(i interface{}) string {
	switch i.(type) {

	case map[bool]string:
		return "map[bool]string"
	case map[bool]int:
		return "map[bool]int"
	case map[bool]float32:
		return "map[bool]float32"
	case map[bool]bool:
		return "map[bool]bool"

	case map[float32]string:
		return "map[float32]string"
	case map[float32]int:
		return "map[float32]int"
	case map[float32]float32:
		return "map[float32]float32"
	case map[float32]bool:
		return "map[float32]bool"

	case map[string]string:
		return "map[string]string"
	case map[string]int:
		return "map[string]int"
	case map[string]float32:
		return "map[string]float32"
	case map[string]bool:
		return "map[string]bool"

	case map[int]string:
		return "map[int]string"
	case map[int]int:
		return "map[int]int"
	case map[int]float32:
		return "map[int]float32"
	case map[int]bool:
		return "map[int]bool"
	default:
		panic("type unknown")
	}
}
func showMeTheType(i interface{}) string {
	if simpleType(i) == "none" {
		return mapType(i)
	} else {
		return simpleType(i)
	}
}

func main() {
	var y string
	for _, tt := range []interface{}{1, uint(2), int8(4), float64(7), "", 'w', []int{}, map[string]bool{}} {
		y = fmt.Sprintf("%T", tt)
		fmt.Println(y, " - ", showMeTheType(tt))
	}
}
