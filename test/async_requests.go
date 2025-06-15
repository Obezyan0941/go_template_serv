package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func makeRequest(url string, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка при запросе к %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Ответ от %s: статус %d, длина тела %d байт\n", url, resp.StatusCode, len(body))
}

func main() {
	url := "http://localhost:8080/async" // URL для запроса
	requestCount := 220                  // Количество одновременных запросов

	var wg sync.WaitGroup

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go makeRequest(url, &wg)
	}

	wg.Wait() // Ждём завершения всех горутин
	fmt.Println("Все запросы завершены.")
}
