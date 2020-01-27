package app

import (
	"encoding/json"
	"errors"
	"finleap/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type webhookAction struct {
	action  string
	webhook *model.Webhook
}

func (a *App) webhookStoreRoutine() {
	for {
		webhook := <-a.webhookChan
		if webhook.action == "add" {
			err := a.addWebhook(webhook.webhook)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			err := a.deleteWebhook(webhook.webhook)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (a *App) webhookRoutine() {
	for {
		temp := <-a.newTemperature
		err := a.sendTemperature(temp)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (a *App) restoreWebhooks() {
	webhooks, err := model.GetWebhooks(a.DB)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, webhook := range webhooks {
		a.webhookChan <- webhookAction{action: "add", webhook: &webhook}
	}

}

func (a *App) addWebhook(webhook *model.Webhook) error {
	if webhook.ID == 0 || webhook.CityID == 0 || webhook.CallbackURL == "" {
		return errors.New("Invalid webhook")
	}
	a.Webhooks.lock.Lock()
	defer a.Webhooks.lock.Unlock()
	a.Webhooks.Webhooks = append(a.Webhooks.Webhooks, webhook)
	return nil
}

func (a *App) deleteWebhook(webhook *model.Webhook) error {
	a.Webhooks.lock.Lock()
	defer a.Webhooks.lock.Unlock()
	var index int
	var exists = false
	for id, hook := range a.Webhooks.Webhooks {
		if hook.ID == webhook.ID {
			index = id
			exists = true
		}
	}
	if !exists {
		return errors.New("Webhook not found")
	}
	copy(a.Webhooks.Webhooks[index:], a.Webhooks.Webhooks[index+1:])
	a.Webhooks.Webhooks[len(a.Webhooks.Webhooks)-1] = nil
	a.Webhooks.Webhooks = a.Webhooks.Webhooks[:len(a.Webhooks.Webhooks)-1]
	return nil
}

//handler for "/webhooks" POST endpoint
func (a *App) handleCreateWebhook(w http.ResponseWriter, r *http.Request) {
	var webhook *model.Webhook
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&webhook); err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid resquest payload")})
		return
	}
	if webhook.CityID == 0 {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid city id %v", r.FormValue("city_id"))})
		return
	}
	if webhook.CallbackURL == "" {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid callback_url value '%v'", r.FormValue("callback_url"))})
		return
	}
	defer r.Body.Close()
	err := webhook.Create(a.DB)
	if err != nil {
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	a.webhookChan <- webhookAction{action: "add", webhook: webhook}
	respondWithJSON(w, http.StatusCreated, webhook)
}

//handler for "/webhooks/:id" DELETE endpoint
func (a *App) handleDeleteWebhook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, Error{Code: http.StatusBadRequest, Error: fmt.Sprintf("Invalid webhook id %v", params["id"])})
		return
	}
	webhook := &model.Webhook{ID: id}
	err = webhook.Delete(a.DB)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, Error{Code: http.StatusNotFound, Error: err.Error()})
			return
		}
		respondWithError(w, Error{Code: http.StatusInternalServerError, Error: err.Error()})
		return
	}
	a.webhookChan <- webhookAction{action: "delete", webhook: webhook}
	respondWithJSON(w, http.StatusCreated, webhook)
}
