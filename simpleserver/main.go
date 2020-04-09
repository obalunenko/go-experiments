package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var data map[string]*storage

type storage struct {
	name   string
	values []int
}

type service struct {
	data      map[string]*storage
	dataMutex sync.Mutex
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methods := map[string]bool{
		"/data": true,
		"/":     true,
	}

	if !methods[r.URL.Path] {
		http.Error(w, fmt.Sprintf("[%s] not found", r.RequestURI), http.StatusNotFound)
		return
	}
	switch r.URL.Path {
	case "/data":
		switch r.Method {
		case http.MethodGet:
			s.getData(w, r)
		case http.MethodPut:
			s.setData(w, r)
		case http.MethodPost:
			s.udpateData(w, r)
		default:
			http.Error(w, fmt.Sprintf("[%s] unsupported method", r.Method), http.StatusMethodNotAllowed)
		}

	case "/":
		makeResponse(w, "Welcome to main")
		return
	}

}

func main() {
	data = make(map[string]*storage)

	svc := service{
		data:      make(map[string]*storage),
		dataMutex: sync.Mutex{},
	}

	mux := http.DefaultServeMux
	mux.HandleFunc("/", indexHandler)
	mux.Handle("/data", &svc)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type requestBody struct {
	Key  string `json:"key"`
	Data []int  `json:"data"`
}

func (s *service) setData(w http.ResponseWriter, r *http.Request) {
	var rb requestBody

	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.dataMutex.Lock()
	data[rb.Key].values = append(data[rb.Key].values, rb.Data...)
	s.dataMutex.Unlock()

	makeResponse(w, "Success")
}

func (s *service) udpateData(w http.ResponseWriter, r *http.Request) {
	var rb requestBody
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.dataMutex.Lock()
	data[rb.Key].values = append(data[rb.Key].values, rb.Data...)
	s.dataMutex.Unlock()

	makeResponse(w, "Success")
}

func (s *service) getData(w http.ResponseWriter, r *http.Request) {
	s.dataMutex.Lock()
	responseData := data[r.URL.Query().Get("key")]
	s.dataMutex.Unlock()

	if len(responseData.values) == 0 {
		makeResponse(w, "no data")
		return
	}

	makeResponse(w, responseData)
}

func makeResponse(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	makeResponse(w, "Hellp there!")
}
