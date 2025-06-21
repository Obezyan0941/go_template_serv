package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db_manager "test_backend/internal/db"

	"github.com/go-pg/pg/v10"
)

type AddUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
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

func (s *Server) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var req AddUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if req.Name == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	// checks if user exists
	userData, err := db_manager.GetUserByName(s.db, req.Name)
	if err == nil {
		errMsg := fmt.Sprintf("User with this name already exists: %s", userData.Name)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	} else if err != pg.ErrNoRows {
		errMsg := fmt.Sprintf("Error retreiving userData: %v", err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	// hashing password
	req.Password, err = db_manager.HashPassword(req.Password)
	if err != nil {
		errMsg := fmt.Sprintf("Could not hash password: %v", err)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	new_user := &db_manager.User{Name: req.Name, Password: req.Password}
	a, err := s.db.Model(new_user).Insert()
	if err != nil {
		errMsg := fmt.Sprintf("Could not add user in db:\t%v", a)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	respondWithJSON(w, http.StatusCreated, req.Name)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
