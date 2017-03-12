package main

import (
	"fmt"
	"strings"
)

// передается int и slice []int
type memoizeFunction func(int, ...int) interface{}

// TODO реализовать
// функции истинны - возвращают один и тот же результат для одинх и тех же аргументов
var fibonacci memoizeFunction
var romanForDecimal memoizeFunction

//TODO Write memoization function
// сохранение результатов выполнения функций для предотвращения повторных вычислений
func memoize(function memoizeFunction) memoizeFunction {
	temporal := make(map[string]interface{})

	return func(x int, anotherArg ...int) interface{} {
		key := string(x)
		for _, arg := range anotherArg {
			key += fmt.Sprintf(":%d", arg)
		}
		if value, ok := temporal[key]; ok {
			return value
		}
		//вызов рекурсивный
		value := function(x)
		temporal[key] = value
		return value
	}
}

// Будет вызвана после ипортирования всех пакетов и создания констант
// TODO обернуть функции fibonacci и roman в memoize
func init() {
	fibonacci = memoize(func(x int, _ ...int) interface{} {
		if x < 2 {
			return x
		}
		return fibonacci(x-1).(int) + fibonacci(x-2).(int)
	})
	decimals := []int{
		1000,
		900,
		500,
		400,
		100,
		90,
		50,
		40,
		10,
		9,
		5,
		4,
		1}
	romans := []string{
		"M",
		"CM",
		"D",
		"CD",
		"C",
		"XC",
		"L",
		"XL",
		"X",
		"IX",
		"V",
		"IV",
		"I"}

	romanForDecimal = memoize(func(x int, _ ...int) interface{} {
		if x < 0 {
			panic("romanForDecimal not work for number < 0")
		}
		var transDecRom string
		for arg, decimal := range decimals {
			if decimal > x {
				continue
			}
			remain := x / decimal
			x %= decimal
			//повтор строчки remain раз
			transDecRom += strings.Repeat(romans[arg], remain)
			// i := 0
			// for {
			// 	i++
			// 	if i > remain {
			// 		break
			// 	}
			// 	transDecRom += romans[arg]
			// }
		}
		return transDecRom
	})
}

func main() {
	fmt.Println("Fibonacci(45) =", fibonacci(45))
	for _, x := range []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 16, 17, 18, 19, 20, 25, 30, 40, 50, 60, 69, 70, 80,
		90, 99, 100, 200, 300, 400, 500, 600, 666, 700, 800, 900,
		1000, 1009, 1444, 1666, 1945, 1997, 1999, 2000, 2008, 2010,
		2012, 2500, 3000, 3999, 4000, 4539} {
		fmt.Printf("%4d = %s\n", x, romanForDecimal(x))
	}
}
