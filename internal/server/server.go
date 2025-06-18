package server

import (
	"fmt"
	"log"
	"time"
	"sync"
	"net/http"

	"test_backend/internal/db"
	"test_backend/internal/config"

	"github.com/go-pg/pg/v10"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router 		*chi.Mux
	mu	   		sync.Mutex
	counter		int
	db     		*pg.DB
}

func NewServer() *Server {
	db_init_data, err := config.LoadDBConfig()
	db, err := db_manager.NewPostgresConnection(*db_init_data)
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		router: chi.NewRouter(),
		counter: 0,
		db: db,
	}
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	// Middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(httprate.LimitByIP(10, 1*time.Minute))
	s.router.Use(middleware.Throttle(15))
	s.router.Use(middleware.Timeout(30 * time.Second))

	// Routes
	s.router.Get("/", s.handleHome)
	s.router.Get("/hello", s.helloHandler)
	s.router.Get("/async", s.asyncHandler)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Главная страница"))
}

func (s *Server) asyncHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.counter++
	s.mu.Unlock()

	go func(c int) {
		time.Sleep(5 * time.Second)
		fmt.Println("Async task is over!\t", c)
	}(s.counter)

	fmt.Fprint(w, "Task is running in background")
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprint(w, "Hello, World!")
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
