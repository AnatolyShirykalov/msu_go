// Один из механизмов синхронизации - каналы
// Каналы, это объект через который можно обеспечить взаимодействие нескольких горутин
// В принимающей (или возвращающей) канал функции, можно указать направление работы с каналом
// Только для чтения - "<-chan" или только для записи "chan<-"
package main

import "fmt"

var c chan int

func main() {
	// Создаем канал
	c := make(chan string, 2)
	// стартуем пишущую горутину
	go greet(c)
	for i := 0; i < 5; i++ {
		// Читаем пару строк из канала
		fmt.Println("1ree")

		fmt.Println(<-c, "f")
		fmt.Println("2ree")
		fmt.Println("vjbkl", <-c)
		fmt.Println("3ree")

	}
	fmt.Println("ree")
	stuff := make(chan int, 7)
	for i := 0; i < 19; i = i + 3 {
		stuff <- i
	}
	close(stuff)
	fmt.Println("Res", process(stuff))
}

func greet(c chan<- string) {
	// Запускаем бесконечный цикл
	for {
		// и пишем в канал пару строк
		// Подпрограмма будет заблокирована до того, как кто-то захочет прочитать из канала
		fmt.Println("111ddd")

		c <- fmt.Sprintf("Владыка")
		fmt.Println("22ddd")
		c <- fmt.Sprintf("Штурмовик")
		fmt.Println("3ddd")
	}
}

func process(input <-chan int) (res int) {
	for r := range input {
		res += r
	}
	return
}
