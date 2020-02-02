package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Deewai/finleap/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	// "time"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAddWebhookInvalidWebhook(t *testing.T) {
	a := App{}
	a.Webhooks.Webhooks = nil
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	a.webhookChan <- webhookAction{
		action:  "add",
		webhook: &model.Webhook{},
	}
	assert.Equal(t, 0, len(a.Webhooks.Webhooks))
}

func TestAddWebhookValidWebhook(t *testing.T) {
	a := App{}
	a.Webhooks.Webhooks = nil
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	a.webhookChan <- webhookAction{
		action: "add",
		webhook: &model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "http://google.com",
		},
	}
	assert.Equal(t, 1, len(a.Webhooks.Webhooks))
}

func TestDeleteWebhookInvalidWebhookID(t *testing.T) {
	a := App{}
	a.Webhooks.Webhooks = []*model.Webhook{
		&model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "http://google.com",
		},
	}
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	a.webhookChan <- webhookAction{"delete", &model.Webhook{
		ID: 0,
	}}
	assert.Equal(t, 1, len(a.Webhooks.Webhooks))
}

func TestDeleteWebhookValidWebhookID(t *testing.T) {
	a := App{}
	a.Webhooks.Webhooks = []*model.Webhook{
		&model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "http://google.com",
		},
	}
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	a.webhookChan <- webhookAction{"delete", &model.Webhook{
		ID: 1,
	}}
	assert.Equal(t, 0, len(a.Webhooks.Webhooks))
}

func TestRestoreWebhooksDatabaseError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	a := App{}
	a.Webhooks.Webhooks = nil
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a.DB = db

	mock.ExpectQuery("^SELECT (.+) FROM webhooks$").WillReturnError(fmt.Errorf("Error fetching result from database"))
	a.restoreWebhooks()
	assert.True(t, strings.Contains(buf.String(), "Error fetching result from database"))
	time.Sleep(2 * time.Second)
	assert.Equal(t, 0, len(a.Webhooks.Webhooks))
}

func TestRestoreWebhooksWhenWebhooksDoesNotExistsInDB(t *testing.T) {
	a := App{}
	a.Webhooks.Webhooks = nil
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a.DB = db
	rows := sqlmock.NewRows([]string{"id", "city_id", "callback_url"})
	mock.ExpectQuery("^SELECT (.+) FROM webhooks$").WillReturnRows(rows)
	a.restoreWebhooks()
	time.Sleep(2 * time.Second)
	assert.Equal(t, 0, len(a.Webhooks.Webhooks))
}

func TestRestoreWebhooksWhenWebhooksExistsInDB(t *testing.T) {
	a := App{}
	a.Webhooks.Webhooks = nil
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a.DB = db
	rows := sqlmock.NewRows([]string{"id", "city_id", "callback_url"}).
		AddRow(1, 1, "http.google.com")

	mock.ExpectQuery("^SELECT (.+) FROM webhooks$").WillReturnRows(rows)
	a.restoreWebhooks()
	//wait for goroutine to add webhook
	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, len(a.Webhooks.Webhooks))
	assert.Equal(t, 1, a.Webhooks.Webhooks[0].ID)
	assert.Equal(t, 1, a.Webhooks.Webhooks[0].CityID)
	assert.Equal(t, "http.google.com", a.Webhooks.Webhooks[0].CallbackURL)
}

func TestWebhookRoutineSendTemperatureReturnsError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	a := App{}
	a.newTemperature = make(chan model.Temperature)
	go a.webhookRoutine()
	a.newTemperature <- model.Temperature{}
	time.Sleep(2 * time.Second)
	assert.True(t, strings.Contains(buf.String(), "Missing fields in temperature object"))
}

func TestHandleCreateWebhookInvalidHttpMethod(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/webhooks", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandleCreateWebhookInvalidFormData(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/webhooks", bytes.NewBuffer([]byte(`{"name":"test user"`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleCreateWebhookValidFormData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	a.Webhooks.Webhooks = nil
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	mock.ExpectExec("INSERT INTO webhooks").WillReturnResult(sqlmock.NewResult(1, 1))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/webhooks", bytes.NewBuffer([]byte(`{"city_id":1,"callback_url":"http://google.com"}`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["id"])
	assert.EqualValues(t, 1, m["city_id"])
	assert.Equal(t, "http://google.com", m["callback_url"])
}

func TestHandleDeleteWebhookInvalidHttpMethod(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/webhooks/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandleDeleteWebhookNotExistingWebhookID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectQuery("^SELECT (.+) FROM webhooks (.+)").WillReturnError(fmt.Errorf("no rows in result set"))
	mock.ExpectExec("DELETE FROM webhooks").WillReturnError(fmt.Errorf("no rows in result set"))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("DELETE", "/webhooks/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHandleDeleteWebhookInvalidWebhookID(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("DELETE", "/webhooks/me", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleDeleteWebhookValidWebhookID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	a.Webhooks.Webhooks = []*model.Webhook{
		&model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "http://google.com",
		},
	}
	a.webhookChan = make(chan webhookAction)
	go a.webhookStoreRoutine()
	rows := sqlmock.NewRows([]string{"id", "city_id", "callback_url"}).
		AddRow(1, 1, "http://google.com")
	mock.ExpectQuery("^SELECT (.+) FROM webhooks (.+)").WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM webhooks (.+) ").WillReturnResult(sqlmock.NewResult(1, 1))

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("DELETE", "/webhooks/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["id"])
	assert.EqualValues(t, 1, m["city_id"])
	assert.Equal(t, "http://google.com", m["callback_url"])
}
