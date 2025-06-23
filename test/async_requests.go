package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"strings"
)

func SendPostRequest(url string, wg *sync.WaitGroup, i int, requestBody map[string]string, token string, requestType string) error {
	defer wg.Done()
	if 	strings.ToUpper(requestType) != "POST" && strings.ToUpper(requestType) != "GET" {
		return fmt.Errorf("requestType should either be POST or GET")
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(requestType, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
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

	fmt.Println(strconv.Itoa(i) + ". body:\t" + string(body) + "status: " + resp.Status)
	return nil
}

func MakeRequest(url string, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making requests %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Printf("Status: %d, body: %s", resp.StatusCode, string(body))
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
	url := "http://localhost:8012/authaction"
	requestCount := 1

	var wg sync.WaitGroup
	var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTEsIm5hbWUiOiJiaWJhIiwic3ViIjoiYmliYSIsImV4cCI6MTc1MDcwNzAzNywiaWF0IjoxNzUwNzA2MTM3LCJqdGkiOiI3NGE1OTAyYS0yMDQyLTQxZmMtOTkwNS00ZGQ4ZDczYzNiYWQifQ.xW8twWfk1vd3gCx3R5ayX3z68PiMz1QS3vUiuEWgih4"

	requestBody := map[string]string{
		"name":     "biba",
		"password": "boba",
	}

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go SendPostRequest(url, &wg, i, requestBody, token, "GET")
	}

	wg.Wait()
	fmt.Println("Все запросы завершены.")
}

func main() {
	RunAsyncPostRequests()
}
