package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateCitiesInvalidHttpMethod(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/cities", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandleCreateCitiesInvalidFormData(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/cities", bytes.NewBuffer([]byte(`{"name":"test user"`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleCreateCitiesValidFormData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectExec("INSERT INTO cities").WillReturnResult(sqlmock.NewResult(1, 1))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/cities", bytes.NewBuffer([]byte(`{"name":"Berlin","latitude":52.520008,"longitude":13.404954}`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["id"])
	assert.EqualValues(t, "Berlin", m["name"])
	assert.EqualValues(t, 52.520008, m["latitude"])
	assert.EqualValues(t, 13.404954, m["longitude"])
}

func TestHandleCreateCitiesExistingName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectExec("INSERT INTO cities").WillReturnError(fmt.Errorf("Duplicate key for column name"))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("POST", "/cities", bytes.NewBuffer([]byte(`{"name":"Berlin","latitude":52.520008,"longitude":13.404954}`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, "Duplicate key for column name", m["error"])
}


func TestHandleUpdateCitiesInvalidFormData(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("PATCH", "/cities/1", bytes.NewBuffer([]byte(`{"name":"test user"`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleUpdateCitiesValidFormData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectExec("UPDATE cities SET (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("PATCH", "/cities/1", bytes.NewBuffer([]byte(`{"name":"Berlin","latitude":52.520008,"longitude":13.404954}`)))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["id"])
	assert.EqualValues(t, "Berlin", m["name"])
	assert.EqualValues(t, 52.520008, m["latitude"])
	assert.EqualValues(t, 13.404954, m["longitude"])
}

func TestHandleDeleteCitiesInvalidHttpMethod(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("GET", "/cities/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandleDeleteCitiesNotExistingCityID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	mock.ExpectQuery("^SELECT (.+) FROM cities (.+)").WillReturnError(fmt.Errorf("no rows in result set"))
	mock.ExpectExec("DELETE FROM cities").WillReturnError(fmt.Errorf("no rows in result set"))
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("DELETE", "/cities/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, "no rows in result set", m["error"])
}

func TestHandleDeleteCitiesInvalidCityID(t *testing.T) {
	a := App{}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("DELETE", "/cities/me", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleDeleteCitiesValidCityID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	a := App{}
	a.DB = db
	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude"}).
		AddRow("Berlin", 52.520008, 13.404954)
	mock.ExpectQuery("^SELECT (.+) FROM cities (.+)").WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM cities (.+) ").WillReturnResult(sqlmock.NewResult(1, 1))

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	req, _ := http.NewRequest("DELETE", "/cities/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.EqualValues(t, 1, m["id"])
	assert.EqualValues(t, "Berlin", m["name"])
	assert.EqualValues(t, 52.520008, m["latitude"])
	assert.EqualValues(t, 13.404954, m["longitude"])
}
