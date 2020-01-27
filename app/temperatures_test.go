package app

import (
	"bytes"
	"fmt"
	"encoding/json"
	"errors"
	"finleap/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSendTemperatureMissingFields(t *testing.T) {
	a := App{}
	err := a.sendTemperature(model.Temperature{})
	assert.NotNil(t, err)
	assert.EqualValues(t, errors.New("Missing fields in temperature object"), err)
}

func TestSendTemperatureCorrectFieldsNoMockUrl(t *testing.T) {
	a := App{}
	FlushMockups()
	a.Webhooks.Webhooks = []*model.Webhook{
		&model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "https://my.service.com/high-temperature",
		},
	}
	err := a.sendTemperature(model.Temperature{
		CityID:    1,
		Min:       10,
		Max:       20,
		Timestamp: 10000,
	})
	assert.NotNil(t, err)
	assert.EqualValues(t, errors.New("No mockup found for given request"), err)
}

func TestSendTemperatureCorrectFieldsNoWebhook(t *testing.T) {
	a := App{}
	err := a.sendTemperature(model.Temperature{
		CityID:    1,
		Min:       10,
		Max:       20,
		Timestamp: 10000,
	})
	assert.Nil(t, err)
}

func TestSendTemperatureCorrectFieldsInvalidUrl(t *testing.T) {
	a := App{}
	FlushMockups()
	a.Webhooks.Webhooks = []*model.Webhook{
		&model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "https://my.service.com/high-temperature",
		},
	}
	AddMockups(mock{
		url:        "https://my.service.com/high-temperature",
		httpMethod: http.MethodPost,
		err:        errors.New("invalid response"),
	})
	err := a.sendTemperature(model.Temperature{
		CityID:    1,
		Min:       10,
		Max:       20,
		Timestamp: 10000,
	})
	assert.NotNil(t, err)
	assert.EqualValues(t, errors.New("invalid response"), err)
}

func TestSendTemperatureCorrectFieldsValidUrl(t *testing.T) {
	a := App{}
	FlushMockups()
	a.Webhooks.Webhooks = []*model.Webhook{
		&model.Webhook{
			ID:          1,
			CityID:      1,
			CallbackURL: "https://my.service.com/high-temperature",
		},
	}
	AddMockups(mock{
		url:        "https://my.service.com/high-temperature",
		httpMethod: http.MethodPost,
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`"message":"success"`)),
		},
	})
	err := a.sendTemperature(model.Temperature{
		CityID:    1,
		Min:       10,
		Max:       20,
		Timestamp: 10000,
	})
	assert.Nil(t, err)
}

func TestHandleCreateTemperatureInvalidHttpMethod(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/temperatures", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandleCreateTemperatureWithInvalidFormData(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/temperatures", bytes.NewBuffer([]byte(`{"name":"test user"`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleCreateTemperatureWithDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectExec("INSERT INTO temperatures").WillReturnError(fmt.Errorf("a database error"))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/temperatures", bytes.NewBuffer([]byte(`{"city_id":1,"max":40,"min":10}`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandleCreateTemperatureWithValidFormData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	a.Webhooks.Webhooks = nil
	a.newTemperature = make(chan model.Temperature)
	go a.webhookRoutine()
	mock.ExpectExec("INSERT INTO temperatures").WillReturnResult(sqlmock.NewResult(1, 1))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/temperatures", bytes.NewBuffer([]byte(`{"city_id":1,"max":40,"min":10}`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["id"])
	assert.EqualValues(t, 1, m["city_id"])
	assert.EqualValues(t, 40, m["max"])
	assert.EqualValues(t, 10, m["min"])
}

func TestHandleForecastWithInValidCityID(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/forecasts/me", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleForecastWithNotExistingCityID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectQuery("^SELECT (.+) FROM temperatures (.+)").WillReturnError(fmt.Errorf("City data doesn't exist"))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/forecasts/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, "City data doesn't exist", m["error"])
}

func TestHandleForecastWithValidCityID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	rows := sqlmock.NewRows([]string{"id", "city_id", "max", "min",}).
		AddRow(1, 1, 30, 10).
		AddRow(1, 1, 20, 5)
	mock.ExpectQuery("^SELECT (.+) FROM temperatures (.+)").WillReturnRows(rows)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/forecasts/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["city_id"])
	assert.EqualValues(t, 25, m["max"])
	assert.EqualValues(t, 7.5, m["min"])
	assert.EqualValues(t, 2, m["sample"])
}

