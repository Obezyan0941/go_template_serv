package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	db_manager "test_backend/internal/db"
	"test_backend/internal/token"

	"github.com/go-pg/pg/v10"
)

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

	fmt.Fprint(w, "Task is running in background\n")
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

func (s *Server) LogInUserHandler(w http.ResponseWriter, r *http.Request) {
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

	userData, err := db_manager.GetUserByName(s.db, req.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("user not found: %d. Error: %v", userData.Id, err), http.StatusUnauthorized)
		return
	}

	accessToken, _, err := s.jwt_maker.CreateToken(userData.Id, userData.Name, 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	res := LoginUserResponse{
		AccessToken: accessToken,
		User: UserResponse{
			Name: userData.Name,
			ID:   userData.Id,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	println(accessToken)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) AuthorizedAction(w http.ResponseWriter, r *http.Request) {
	// read the authorization header
	// verify the token
	claims, err := verifyClaimsFromAuthHeader(r, s.jwt_maker)
	if err != nil {
		http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
		return
	}

	userid := claims.ID
	userData, err := db_manager.GetUserDataByID(int(userid), s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("user not found: %d. Error: %v", userid, err), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}

func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := fields[1]
	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
