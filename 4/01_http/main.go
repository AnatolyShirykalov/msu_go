package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	http.Handler
}

func (h Handler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Создаем http клиент. В стуктуру можно передать таймаут, куки и прочую инфу о запросе
	c := http.Client{}
	resp, err := c.Get("http://artii.herokuapp.com/make?text=" + path)
	if err != nil {
		log.Println(err)
	}
	// нужно закрыть тело, когда прочитаем что нужно
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("11", body, "22")

	// статус - ОК
	// http.StatusOK
	w.WriteHeader(200)
	w.Write(body)
}

// http текствый запрос
// keep-alive - не будем закрывать соединение после того как нам пришел ответ
// user-agent - с какого устроиства сделали запрос
// accept-encoding - кодировка, которую поддерживает браузер

func main() {
	// на сервер приходит какой-то путь ресурса
	// надо определить по какому пути слушать порт

	// w - структура, которая вернеться пользователю
	// r - запрос, который нам придет

	handler := Handler{}

	http.Handle("/", handler)

	http.ListenAndServe(":8081", nil)
}
