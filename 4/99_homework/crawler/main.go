package main

// Функция Crawl(url string) []string
// * На вход передается web адрес корня сайта (/)
// * на выход список всех страниц, на которые есть ссылки на всех страницах этого домена.

// Неизвестные урлы учитывать не надо.
// Внешние урлы (хосты, отличные от текущего) к результату не относятся.
// Результат отсортирован в порядке появления на сайте.

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Crawl(host string) []string {
	slisLink := make(map[string]bool)
	slisLink["/page1.html"] = true
	final := make(map[string]int)
	final["/page1.html"] = 0

	while(slisLink, host, final)
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
				// fmt.Println(ref)
				slisLink[ref] = true
			} else {
				panic("a href not found >")
			}
		}
	}
	return slisLink
}
func while(input map[string]bool, host string, final map[string]int) {
	for aref := range input {
		exist := make(map[string]bool)
		var flag bool

		client := http.Client{}
		resp, _ := client.Get(host + aref)
		bodyBytes, ok := ioutil.ReadAll(resp.Body)
		if ok != nil {
			panic(ok)
		}
		resp.Body.Close()
		if len(bodyBytes) > 6 {
			if string(bodyBytes[:3]) == "404" {
				break
			}
			if string(bodyBytes[:6]) == "<html>" {

				if aref != "/page1.html" {
					lee := len(final)
					final[aref] = lee
				}
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
			}
		}
	}
	return
}
