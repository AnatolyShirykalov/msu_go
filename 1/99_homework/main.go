package main

import "fmt"

func ReturnInt() int {
	return 1
}
func ReturnFloat() float32 {
	return 1.1
}
func ReturnIntArray() [3]int {
	return [3]int{1, 3, 4}
}
func ReturnIntSlice() []int {
	return []int{1, 2, 3}
}
func IntSliceToString(sl []int) string {
	var r string
	for _, val := range sl {
		r += fmt.Sprintf("%d", val)
	}
	return r
}

func MergeSlices(sl1 []float32, sl2 []int32) []int {
	sl := make([]int, 0, len(sl1)+len(sl2))
	for _, val := range sl1 {
		sl = append(sl, int(val))
	}
	for _, val := range sl2 {
		sl = append(sl, int(val))
	}
	return sl
}
func GetMapValuesSortedByKey(input map[int]string) []string {
	i := 0
	sl := make([]string, 0, len(input))
	for {
		if i == len(input) {
			break
		}
		i++
		sl = append(sl, input[i])
	}
	return sl
}
