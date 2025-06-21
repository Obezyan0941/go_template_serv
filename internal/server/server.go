package server

import (
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

	db_manager.CreateSchema(db, (*db_manager.User)(nil))

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
	s.router.Post("/adduser", s.AddUserHandler)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
