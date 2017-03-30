package main

// Функция Crawl(url string) []string
// * На вход передается web адрес корня сайта (/)
// * на выход список всех страниц, на которые есть ссылки на всех страницах этого домена.

// Неизвестные урлы учитывать не надо.
// Внешние урлы (хосты, отличные от текущего) к результату не относятся.
// Результат отсортирован в порядке появления на сайте.

import (
	"io/ioutil"
	"net/http"
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
	index := 0
	start := ""
	str := ""

	lenabaseref := len("base href=")
	for {
		if index < len(bodyBytes)-lenabaseref {
			baseref := string(bodyBytes[index : index+lenabaseref])
			if baseref == "base href=" {
				index += lenabaseref + 1
				str = ""
				for link := index; link < len(bodyBytes); link++ {
					if string(bodyBytes[index]) != "/" {
						break
					}
					if string(bodyBytes[link+3]) == ">" {
						index += link
						break
					}
					str = str + string(bodyBytes[link])
				}
				if str != "" {
					if str != "/index.html" {
						start = str
					}
					break
				}
			}
			index++
		} else {
			index = 0
			break
		}
	}

	slisLink := make(map[string]bool)
	lenahref := len("a href=")
	for {
		if index < len(bodyBytes)-lenahref {
			ahref := string(bodyBytes[index : index+lenahref])
			if ahref == "a href=" {
				index += lenahref + 1
				str = ""

				for link := index; link < len(bodyBytes); link++ {
					if string(bodyBytes[link+1]) == ">" {
						break
					}
					if string(bodyBytes[index:index+4]) == "http" {
						break
					}
					str = str + string(bodyBytes[link])
				}
				if str != "" {
					slisLink[start+str] = true
				}
			}
			index++
		} else {
			break
		}
	}
	return slisLink
}
func while(input map[string]bool, host string, final map[string]int) {
	for aref := range input {
		exist := make(map[string]bool)
		var flag bool

		client := http.Client{}
		resp1, _ := client.Get(host + aref)
		bodyBytes1, ok := ioutil.ReadAll(resp1.Body)
		if ok != nil {
			break
		}
		defer resp1.Body.Close()

		if len(bodyBytes1) > 6 {
			if string(bodyBytes1[:3]) == "404" {
				break
			}
			if string(bodyBytes1[:6]) == "<html>" {

				if aref != "/page1.html" {
					lee := len(final)
					final[aref] = lee
				}
				sl := Parser(bodyBytes1)

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
