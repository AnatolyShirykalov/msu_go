package main

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Тип HandlerFunc - это адаптер, позволяющий использовать обычные функции в качестве обработчиков HTTP
// HandlerFunc (f) является обработчиком, который вызывает f

// w - структура, которая вернеться пользователю
// r - запрос, который нам придет

func Case(root string, f func(string) []string) []string {
	// NewRequest возвращает новый входящий запрос сервера, подходящий для передачи в http.Handler для тестирования.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.RequestURI
		if req == "/" {
			// Перенаправить ответы(w) на запрос(r) с перенаправлением на URL-адрес,
			// который может быть путем по отношению к пути запроса.
			// 300 - перемещение riderect ( 301 - постоянно перемещенно, 302 - временно перемещенно)
			http.Redirect(w, r, "/page1.html", 301)
		} else {
			// ServeFile отвечает на запрос содержимым именованного файла или каталога.
			// Если предоставленный файл или имя каталога относительный путь,
			// он интерпретируется относительно текущего каталога и может восходить к родительским каталогам.
			// Если предоставленное имя построено на основе пользовательского ввода,
			// оно должно быть дезинфицировано до вызова ServeFile.
			// В качестве меры предосторожности ServeFile будет отклонять запросы,
			// где r.URL.Path содержит элемент пути «..».

			// В качестве особого случая ServeFile перенаправляет любой запрос, где r.URL.Path
			// заканчивается в "/index.html" по тому же пути, без окончательного "index.html".
			// Чтобы избежать таких перенаправлений, измените путь или используйте ServeContent.
			http.ServeFile(w, r, "testdata/"+root+"/"+req)
		}
	}))
	defer ts.Close()
	return f(ts.URL)
}

func TestCrawler(t *testing.T) {
	for ii, tt := range []struct {
		root     string
		expected []string
	}{
		{root: "hosta", expected: []string{"/page1.html", "/page2.html", "/page3.html", "/page4.html", "/page5.html"}},
		{root: "hostb", expected: []string{"/page1.html", "/page2.html"}},
		{root: "hostd", expected: []string{"/page1.html", "/subdir/page2.html", "/page3.html"}},
	} {
		res := Case(tt.root, Crawl)
		if !reflect.DeepEqual(res, tt.expected) {
			t.Errorf("Case [%d]: failed, expected %v, got %v", ii, tt.expected, res)
		}
	}
}
