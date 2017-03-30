package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func main() {
	todos := []Todo{
		{"Выучить Go", false},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// здесь надо отдать статический файл, который будет общаться с API из браузера
		// открываем файл
		fileContents, err := ioutil.ReadFile("index.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// и выводим содержимое файла
		w.Write(fileContents)
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request ", r.URL.Path)
		defer r.Body.Close()

		// разные методы обрабатываются по-разному
		switch r.Method {
		// GET для получения данных
		case http.MethodGet:
			// преобразуем структуру в json
			productsJson, _ := json.Marshal(todos)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(productsJson)
		// POST для добавления чего-то нового
		case http.MethodPost:
			decoder := json.NewDecoder(r.Body)
			todo := Todo{}
			// преобразуем json запрос в структуру
			err := decoder.Decode(&todo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			todos = append(todos, todo)
		// PUT для обновления существующей информации
		case http.MethodPut:
			// обновленный ид передается в url
			id := r.URL.Path[len("/todos/"):]
			// 10 -система счисления и размер значения 0 - по умолчанию размер инта
			index, _ := strconv.ParseInt(id, 10, 0)
			todos[index].Done = true
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	http.ListenAndServe(":8081", nil)
}
