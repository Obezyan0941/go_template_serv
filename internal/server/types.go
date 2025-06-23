package server

import "time"

type AddUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type LoginUserResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}
