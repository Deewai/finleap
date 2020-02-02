package app

import (
	"encoding/json"
	"fmt"
	"github.com/Deewai/finleap/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//handler for "/cities" POST endpoint
func (a *App) handleCreateCities(w http.ResponseWriter, r *http.Request) {
	var city *model.City
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&city); err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid resquest payload")})
		return
	}
	if city.Name == "" {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: "Invalid name value"})
		return
	}
	if city.Latitude == 0 {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: "Invalid latitude value"})
		return
	}
	if city.Longitude == 0 {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: "Invalid longitude value"})
		return
	}
	defer r.Body.Close()
	err := city.Create(a.DB)
	if err != nil {
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	respondWithJSON(w, http.StatusCreated, city)
}

//handler for "/cities/:id" PATCH endpoint
func (a *App) handleUpdateCities(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid city id %v", params["id"])})
		return
	}
	var city *model.City
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&city); err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid resquest payload")})
		return
	}
	if city.Name == "" {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: "Invalid name value"})
		return
	}
	if city.Latitude == 0 {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: "Invalid latitude value"})
		return
	}
	if city.Longitude == 0 {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: "Invalid longitude value"})
		return
	}
	defer r.Body.Close()
	city.ID = id
	err = city.Update(a.DB)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, Error{Code: http.StatusNotFound, Error: err.Error()})
			return
		}
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	respondWithJSON(w, http.StatusCreated, city)
}

//handler for "/cities/:id" DELETE endpoint
func (a *App) handleDeleteCities(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid city id %v", params["id"])})
		return
	}
	city := &model.City{ID: id}
	err = city.Delete(a.DB)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, Error{Code: http.StatusNotFound, Error: err.Error()})
			return
		}
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	respondWithJSON(w, http.StatusCreated, city)
}
