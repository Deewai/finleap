package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// )

// var a App

// func ensureTablesExists() {
// 	if _, err := a.DB.Exec(cityCreationQuery); err != nil {
// 		log.Fatal(err)
// 	}
// 	if _, err := a.DB.Exec(temperatureCreationQuery); err != nil {
// 		log.Fatal(err)
// 	}
// 	if _, err := a.DB.Exec(webhookCreationQuery); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func clearTable() {
// 	a.DB.Exec("DELETE FROM cities")
// 	a.DB.Exec("ALTER TABLE cities AUTO_INCREMENT = 1")
// 	a.DB.Exec("DELETE FROM temperatures")
// 	a.DB.Exec("ALTER TABLE temperatures AUTO_INCREMENT = 1")
// 	a.DB.Exec("DELETE FROM webhooks")
// 	a.DB.Exec("ALTER TABLE webhooks AUTO_INCREMENT = 1")
// }

// func TestMain(m *testing.M) {
// 	a = App{}
// 	a.Initialize("root", "morerin09", "weather-monster")
// 	ensureTablesExists()
// 	code := m.Run()
// 	clearTable()
// 	os.Exit(code)
// }

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	a.Router.ServeHTTP(rr, req)

// 	return rr
// }

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

// func TestCreateCity(t *testing.T) {
// 	clearTable()
// 	payload := []byte(`{"name":"Berlin","latitude":52.520008,"longitude":13.404954}`)
// 	req, _ := http.NewRequest("POST", "/cities", bytes.NewBuffer(payload))
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusCreated, response.Code)
// 	var m map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &m)
// 	if m["name"] != "Berlin" {
// 		t.Errorf("Expected city name to be 'Berlin'. Got '%v'", m["name"])
// 	}
// 	if m["latitude"] != 52.520008 {
// 		t.Errorf("Expected city latitude to be '52.520008'. Got '%v'", m["latitude"])
// 	}
// 	if m["longitude"] != 13.404954 {
// 		t.Errorf("Expected city longitude to be '13.404954'. Got '%v'", m["latitude"])
// 	}
// 	if m["id"] != 1.0 {
// 		t.Errorf("Expected city ID to be '1'. Got '%v'", m["id"])
// 	}
// }

// func TestUpdateCity(t *testing.T) {
// 	clearTable()
// 	payload := []byte(`{"name":"Berlin","latitude":52.520008,"longitude":13.404954}`)
// 	req, _ := http.NewRequest("POST", "/cities", bytes.NewBuffer(payload))
// 	response := executeRequest(req)
// 	var n map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &n)
// 	payload = []byte(`{"name":"Potsdam","latitude":52.520008,"longitude":13.404954}`)
// 	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/cities/%d", n["id"]), bytes.NewBuffer(payload))
// 	response = executeRequest(req)
// 	checkResponseCode(t, http.StatusCreated, response.Code)
// 	var m map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &m)
// 	if m["name"] != "Potsdam" {
// 		t.Errorf("Expected city name to be 'Potsdam'. Got '%v'", m["name"])
// 	}
// 	if m["latitude"] != 52.520008 {
// 		t.Errorf("Expected city latitude to be '52.520008'. Got '%v'", m["latitude"])
// 	}
// 	if m["longitude"] != 13.404954 {
// 		t.Errorf("Expected city longitude to be '13.404954'. Got '%v'", m["latitude"])
// 	}
// 	if m["id"] != n["id"] {
// 		t.Errorf("Expected city ID to be '%d'. Got '%v'", n["id"], m["id"])
// 	}
// }

// func TestDeleteCity(t *testing.T) {
// 	clearTable()
// 	payload := []byte(`{"name":"Berlin","latitude":52.520008,"longitude":13.404954}`)
// 	req, _ := http.NewRequest("POST", "/cities", bytes.NewBuffer(payload))
// 	response := executeRequest(req)
// 	var n map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &n)
// 	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/cities/%d", n["id"]), bytes.NewBuffer([]byte{}))
// 	response = executeRequest(req)
// 	checkResponseCode(t, http.StatusCreated, response.Code)
// 	var m map[string]interface{}
// 	json.Unmarshal(response.Body.Bytes(), &m)
// 	if m["name"] != "Berlin" {
// 		t.Errorf("Expected city name to be 'Berlin'. Got '%v'", m["name"])
// 	}
// 	if m["latitude"] != 52.520008 {
// 		t.Errorf("Expected city latitude to be '52.520008'. Got '%v'", m["latitude"])
// 	}
// 	if m["longitude"] != 13.404954 {
// 		t.Errorf("Expected city longitude to be '13.404954'. Got '%v'", m["latitude"])
// 	}
// 	if m["id"] != n["id"] {
// 		t.Errorf("Expected city ID to be '%d'. Got '%v'", n["id"], m["id"])
// 	}
// }
