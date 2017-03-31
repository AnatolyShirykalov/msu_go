package main

// Функция Crawl(url string) []string
// * На вход передается web адрес корня сайта (/)
// * на выход список всех страниц, на которые есть ссылки на всех страницах этого домена.

// Неизвестные урлы учитывать не надо.
// Внешние урлы (хосты, отличные от текущего) к результату не относятся.
// Результат отсортирован в порядке появления на сайте.

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Crawl(host string) []string {
	final := make(map[string]int)
	input := make(map[string]bool)

	while(input, host, final)

	res := make([]string, len(final))
	for elem, ind := range final {
		res[ind] = elem
	}

	return res
}

func Parser(bodyBytes []byte) map[string]bool {
	start := ""
	start2 := 0
	lenabaseref := "base href="

	bodyStr := string(bodyBytes)
	start1 := strings.Index(bodyStr, lenabaseref)
	if start1 > -1 {
		start1 += len(lenabaseref)
		start2 = strings.Index(bodyStr, " />")
		start = string(bodyBytes[start1+1 : start2-1])
	}
	if start == "/index.html" {
		start = ""
	}

	lenahref := "a href="
	slisLink := make(map[string]bool)

	bodyStr = string(bodyBytes[start2:])
	splitBodyStr := strings.Split(bodyStr, lenahref)

	if len(splitBodyStr) > 1 {
		for _, elemsplit := range splitBodyStr[1:] {
			start2 = strings.Index(elemsplit, ">")
			if start2 > -1 {
				ref := start + string(elemsplit[1:start2-1])
				if strings.Index(ref, "/") != 0 {
					break
				}
				slisLink[ref] = true
			} else {
				panic("a href not found >")
			}
		}
	}
	return slisLink
}
func while(input map[string]bool, host string, final map[string]int) {
	var flag bool
	for aref := range input {
		exist := make(map[string]bool)

		client := http.Client{}
		resp, _ := client.Get(host + aref)
		bodyBytes, ok := ioutil.ReadAll(resp.Body)
		if ok != nil {
			panic(ok)
		}
		resp.Body.Close()

		if resp.StatusCode == 200 {
			final[aref] = len(final)
			sl := Parser(bodyBytes)
			for ref := range sl {
				if _, ok := final[ref]; !ok {
					exist[ref] = true
					flag = true
				}

			}
			if flag {
				while(exist, host, final)
				break
			}
		} else {
			fmt.Println("status ", string(bodyBytes), "ref ", aref, "\n")
		}
	}
	if len(input) == 0 {
		client := http.Client{}
		resp, _ := client.Get(host)
		input[resp.Request.URL.RequestURI()] = true
		while(input, host, final)
	}
	return
}
