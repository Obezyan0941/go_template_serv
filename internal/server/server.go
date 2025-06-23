package server

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"test_backend/internal/config"
	db_manager "test_backend/internal/db"
	jwt_manager "test_backend/internal/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-pg/pg/v10"
)

type Server struct {
	router    *chi.Mux
	mu        sync.Mutex
	counter   int
	db        *pg.DB
	jwt_maker *jwt_manager.JWTMaker
}

func NewServer() *Server {
	db_init_data, err := config.LoadDBConfig()
	db, err := db_manager.NewPostgresConnection(*db_init_data)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	db_manager.CreateSchema(db, (*db_manager.User)(nil))

	s := &Server{
		router:    chi.NewRouter(),
		counter:   0,
		db:        db,
		jwt_maker: jwt_manager.NewJWTMaker(init_jwt_key()),
	}
	s.configureRouter()
	return s
}

func init_jwt_key() string {
	var jwt_secret string = os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("No JWT secret key found in ENV")
	}
	if len(jwt_secret) < 32 {
		log.Fatal("JWT secret key should be at least 32 bytes long")
	}
	return jwt_secret
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
	s.router.Post("/login", s.LogInUserHandler)
	s.router.Get("/authaction", s.AuthorizedAction)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
