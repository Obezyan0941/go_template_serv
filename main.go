package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var counter int = 0
var mu sync.Mutex
var limiter = rate.NewLimiter(5, 10)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is a main page!")
}

// rate limiting middleware
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Слишком много запросов, попробуйте позже", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprint(w, "Hello, World!")
}

func asyncHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	counter++
	mu.Unlock()

	go func(c int) {
		time.Sleep(5 * time.Second)
		fmt.Println("Async task is over!\t", c)
	}(counter)

	fmt.Fprint(w, "Task is running in background")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/async", rateLimitMiddleware(asyncHandler))

	fmt.Println("Server is running ib http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
