package main

import (
	"context"
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

// rate limiting throttling middleware
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Создаем контекст с таймаутом ожидания
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		
		// Ждем доступного токена (блокируем выполнение)
		if err := limiter.Wait(ctx); err != nil {
			// Если таймаут ожидания истек
			http.Error(w, "Превышено время ожидания. Попробуйте позже.", http.StatusTooManyRequests)
			return
		}

		// Добавляем заголовки с информацией о лимитах
		w.Header().Set("X-RateLimit-Limit", fmt.Sprint(limiter.Limit()))
		w.Header().Set("X-RateLimit-Burst", fmt.Sprint(limiter.Burst()))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprint(int(limiter.Tokens())))

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
