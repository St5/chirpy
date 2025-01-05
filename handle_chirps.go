package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/St5/goboot-srv/internal/auth"
	"github.com/St5/goboot-srv/internal/database"
	"github.com/google/uuid"
)

type Chirpy struct {
	ID        uuid.UUID `json:"id"`
	CreateAt  string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

/**
 * Handle create chirp
 */
func (confg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type requstChirpy struct {
		Body   string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, confg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	//Decode request
	var chirpReq requstChirpy
	err = json.NewDecoder(r.Body).Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	//Validate chirp
	if len(chirpReq.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	newMsg := validateMsg(chirpReq.Body)

	//Create chirp
	chirpyDb, err := confg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   newMsg,
		UserID: userID,
	})

	//Conver to json convertable format
	chirpy := Chirpy{
		ID:        chirpyDb.ID,
		CreateAt:  chirpyDb.CreatedAt.String(),
		UpdatedAt: chirpyDb.UpdatedAt.String(),
		Body:      chirpyDb.Body,
		UserID:    chirpyDb.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirpy)

}

/**
 * Handle get one chirp by id
 */ 
func (confg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing chirpID")
		return
	}

	chirp, err := confg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, 404, "Chirpy doesn`t found")
		return
	}

	//Convert database.Chirpy to Chirpy model for json response with correct format fields
	chirpy := Chirpy{
		ID:        chirp.ID,
		CreateAt:  chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpy)
}

/**
 * Handle get all chirps
 */
func (confg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {

	authorId := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "asc"
	}

	sortBy = "created_at " + sortBy
	
	var chirps []database.Chirp
	var err error

	if authorId != "" {
		chirps, err = confg.db.GetChirpsByUserID(r.Context(), database.GetChirpsByUserIDParams{
			UserID:  uuid.MustParse(authorId),
			Column2: sortBy,
		})
	} else {
		if sortBy == "created_at asc" {
			chirps, err = confg.db.GetAllChirps(r.Context())
		} else {
			chirps, err = confg.db.GetAllChirpsDesc(r.Context())
		}
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	//Convert database.Chirpy to Chirpy model for json response with correct format fields
	chirpsResponse := make([]Chirpy, len(chirps))
	for i, chirp := range chirps {
		chirpsResponse[i] = Chirpy{
			ID:        chirp.ID,
			CreateAt:  chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirpsResponse)
}

/**
 * Handle delete chirp
 */
func (confg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	//Authenticate user
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, confg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	//Find a chirp
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing chirpID")
		return
	}

	chirp, err := confg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpID))	
	if err != nil {
		respondWithError(w, 404, "Chirpy doesn`t found")
		return
	}

	//Check if user is owner of chirp
	if chirp.UserID != userID {
		respondWithError(w, 403, "Forbidden")
		return
	}

	//Delete chirp
	err = confg.db.DeleteChirpByID(r.Context(), uuid.MustParse(chirpID))

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	respondWithJSON(w, 204, nil)
}

/**
 * Validate message and replace bad words with ****
 */
func validateMsg(msg string) string {
	listOfBadWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(msg, " ")
	newMsg := []string{}
	for _, word := range words {
		if containe(listOfBadWords, word) {
			newMsg = append(newMsg, "****")
			continue
		}
		newMsg = append(newMsg, word)
	}

	return strings.Join(newMsg, " ")
}

/**
 * Check if word is in list
 */
func containe(list []string, word string) bool {
	for _, elem := range list {
		if strings.ToLower(elem) == strings.ToLower(word) {
			return true
		}
	}
	return false
}
