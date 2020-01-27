package app

import (
	"encoding/json"
	"errors"
	"finleap/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (a *App) sendTemperature(temp model.Temperature) error {
	invalidError := errors.New("Missing fields in temperature object")
	if temp.CityID == 0 || temp.Timestamp == 0{
		return invalidError
	}
	receiverUrls := []string{}
	a.Webhooks.lock.Lock()
	for _, hook := range a.Webhooks.Webhooks {
		if hook.CityID == temp.CityID {
			receiverUrls = append(receiverUrls, hook.CallbackURL)
		}
	}
	a.Webhooks.lock.Unlock()
	if len(receiverUrls) == 0{
		return nil
	}
	requestBody, _ := json.Marshal(map[string]interface{}{
		"city_id":   temp.CityID,
		"max":       temp.Max,
		"min":       temp.Min,
		"Timestamp": temp.Timestamp,
	})
	for _, url := range receiverUrls {
		_, err := a.sendRequest(url, requestBody)
		if err != nil {
			return err
		}
	}
	return nil
}

//handler for "/temperatures" POST endpoint
func (a *App) handleCreateTemperature(w http.ResponseWriter, r *http.Request) {
	var temperature *model.Temperature
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&temperature); err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid resquest payload")})
		return
	}
	if temperature.CityID == 0 {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid city id %v", r.FormValue("city_id"))})
		return
	}
	defer r.Body.Close()
	temperature.Timestamp = time.Now().Unix()
	err := temperature.Create(a.DB)
	if err != nil {
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	a.newTemperature <- *temperature
	respondWithJSON(w, http.StatusCreated, temperature)
}

//handler for "/forecasts/:city_id" GET endpoint
func (a *App) handleForecast(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	CityID, err := strconv.Atoi(params["city_id"])
	if err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid city id %v", params["city_id"])})
		return
	}
	timestamp24HoursAgo := time.Now().AddDate(0, 0, -1).Unix()
	temperatures, err := model.GetTemperatures(a.DB, CityID, timestamp24HoursAgo)
	if err != nil {
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	var forecast model.Forecast
	var totalMin int
	var totalMax int
	total := len(temperatures)
	for _, temp := range temperatures {
		totalMax += temp.Max
		totalMin += temp.Min
	}
	forecast = model.Forecast{CityID: CityID, Max: float32(totalMax) / float32(total), Min: float32(totalMin) / float32(total), Sample: total}
	respondWithJSON(w, http.StatusOK, forecast)
}
