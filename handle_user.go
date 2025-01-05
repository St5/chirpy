package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/St5/goboot-srv/internal/auth"
	"github.com/St5/goboot-srv/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type UserToken struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

/**
 * Handle create user
 */
func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//Decode request
	decode := json.NewDecoder(r.Body)
	req := request{}
	err := decode.Decode(&req)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	//Validate password
	pswrd, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	//Create user
	userParams := database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: pswrd,
	}
	userDb, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	user := User{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
		IsChirpyRed: userDb.IsChirpyRed.Bool,
	}
	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		expiresInSeconds int    `json:"expires_in_seconds"`
	}

	//Decode request
	decode := json.NewDecoder(r.Body)
	req := request{}
	err := decode.Decode(&req)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	//Get user
	userDb, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	//Check password
	err = auth.CheckPasswordHash(req.Password, userDb.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	// An hour by default
	expiresInSeconds := 3600

	if req.expiresInSeconds > 0 && req.expiresInSeconds < 3600 {
		expiresInSeconds = req.expiresInSeconds
	}

	//Create token
	token, err := auth.MakeJWT(userDb.ID, cfg.tokenSecret, time.Duration(expiresInSeconds)*time.Second)

	if err != nil {
		respondWithError(w, 500, "Token error")
		return
	}

	//Create refresh token
	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, 500, "Token error")
		return
	}

	record, err := cfg.db.CreateToken(r.Context(), database.CreateTokenParams{
		UserID:    userDb.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		respondWithError(w, 500, "Token error")
		return
	}

	//Create user response
	user := UserToken{
		ID:           userDb.ID,
		CreatedAt:    userDb.CreatedAt,
		UpdatedAt:    userDb.UpdatedAt,
		Email:        userDb.Email,
		IsChirpyRed: userDb.IsChirpyRed.Bool,
		Token:        token,
		RefreshToken: record.Token,
	}

	//Respond
	respondWithJSON(w, 200, user)
}

/**
 * Handle refresh token
 */
func (cfg *apiConfig) handRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	// Create Access Token
	accessToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)

	if err != nil {
		respondWithError(w, 500, "Token error")
		return
	}

	//Respond
	respondWithJSON(w, 200, map[string]string{"token": accessToken})
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	RefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	_, err = cfg.db.GetUserFromRefreshToken(r.Context(), RefreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), RefreshToken)
	if err != nil {
		respondWithError(w, 500, "Token error")
		return
	}

	respondWithJSON(w, 204, nil)

}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//Decode request
	decode := json.NewDecoder(r.Body)
	req := request{}
	err := decode.Decode(&req)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	//Validate password
	pswrd, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	userParam := database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: pswrd,
	}

	userDb, err := cfg.db.UpdateUser(r.Context(), userParam)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	user := User{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
		IsChirpyRed: userDb.IsChirpyRed.Bool,
	}
	respondWithJSON(w, 200, user)
	
}

