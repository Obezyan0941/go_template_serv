package main

import (
	"fmt"
	"io"
	"bytes"
	"strconv"
	"net/http"
	"encoding/json"
	"sync"
)

func SendPostRequest(url string, wg *sync.WaitGroup, i int, name string, password string) error {
	defer wg.Done()
	requestBody := map[string]string{
		"name":     name,
		"password": password,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Println(strconv.Itoa(i) + ":\t" +string(body))
	return nil
}

func MakeRequest(url string, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка при запросе к %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения тела ответа: %v\n", err)
		return
	}
	fmt.Printf("Ответ от %s: статус %d, длина тела %d байт\n", url, resp.StatusCode, len(body))
}

func RunAsyncRequests() {
	url := "http://localhost:8012/async" // URL для запроса
	requestCount := 200                  // Количество одновременных запросов

	var wg sync.WaitGroup

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go MakeRequest(url, &wg)
	}

	wg.Wait() // Ждём завершения всех горутин
	fmt.Println("Все запросы завершены.")
}

func RunAsyncPostRequests() {
	url := "http://localhost:8012/adduser"
	requestCount := 20

	var wg sync.WaitGroup

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go SendPostRequest(url, &wg, i, "biba" + strconv.Itoa(i), "boba" + strconv.Itoa(i+13)  + strconv.Itoa(i * 26))
	}

	wg.Wait()
	fmt.Println("Все запросы завершены.")	
}


func main() {
	RunAsyncPostRequests()
}
