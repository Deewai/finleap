package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var app App

func TestMain(m *testing.M) {
	StartMockups()
	os.Exit(m.Run())
}

func TestApp_sendRequestErrorGotten(t *testing.T) {
	FlushMockups()
	AddMockups(mock{
		url:        "https://my.service.com/high-temperature",
		httpMethod: http.MethodPost,
		err:        errors.New("invalid response"),
	})
	a := &App{}
	payload, _ := json.Marshal(map[string]string{"test": "test"})
	resp, err := a.sendRequest("https://my.service.com/high-temperature", payload)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid response")
}

func TestApp_sendRequestSuccessfulRequest(t *testing.T) {
	FlushMockups()
	AddMockups(mock{
		url:        "https://my.service.com/high-temperature",
		httpMethod: http.MethodPost,
		err:        errors.New("invalid response"),
	})
	a := &App{}
	payload, _ := json.Marshal(map[string]string{"test": "test"})
	resp, err := a.sendRequest("https://my.service.com/high-temperature", payload)
	assert.Nil(t, resp)
	assert.NotNil(t, err)

}
