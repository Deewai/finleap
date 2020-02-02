package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Deewai/finleap/model"
	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	enableMocks = false
	mocks       = make(map[string]*mock)
)

type mock struct {
	url        string
	httpMethod string
	response   *http.Response
	err        error
}

func StartMockups() {
	enableMocks = true
}

func DisableMockups() {
	enableMocks = false
}

func FlushMockups() {
	mocks = make(map[string]*mock)
}

func AddMockups(m mock) {
	mocks[m.url] = &m
}

type App struct {
	Router   *mux.Router
	DB       *sql.DB
	Webhooks struct {
		lock     sync.Mutex
		Webhooks []*model.Webhook
	}
	webhookChan    chan webhookAction
	newTemperature chan model.Temperature
}

type Error struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (a *App) Initialize(host, port, user, password, dbname string) {
	db, err := model.NewConn("mysql", host, port, user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
	a.webhookChan = make(chan webhookAction)
	a.newTemperature = make(chan model.Temperature)
	go a.webhookStoreRoutine()
	go a.webhookRoutine()
	a.restoreWebhooks()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/cities", a.handleCreateCities).Methods("POST")
	a.Router.HandleFunc("/cities/{id}", a.handleUpdateCities).Methods("PATCH")
	a.Router.HandleFunc("/cities/{id}", a.handleDeleteCities).Methods("DELETE")
	a.Router.HandleFunc("/temperatures", a.handleCreateTemperature).Methods("POST")
	a.Router.HandleFunc("/forecasts/{city_id}", a.handleForecast).Methods("GET")
	a.Router.HandleFunc("/webhooks", a.handleCreateWebhook).Methods("POST")
	a.Router.HandleFunc("/webhooks/{id}", a.handleDeleteWebhook).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Printf("http server started on %s", addr)
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func respondWithError(w http.ResponseWriter, e Error) {
	respondWithJSON(w, e.Code, e)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) sendRequest(url string, payload []byte) (*http.Response, error) {
	if enableMocks {
		mock := mocks[url]
		if mock == nil {
			return nil, errors.New("No mockup found for given request")
		}
		return mock.response, mock.err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return resp, err
	}
	return resp, nil
}
