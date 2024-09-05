package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type storage struct {
	Name   string `json:"name"`
	Values []int  `json:"values"`
}

type repository interface {
	setData(key string, data []int)
	getData(key string) *storage
}

type smap struct {
	data sync.Map
}

func (s *smap) setData(key string, data []int) {
	st, ok := s.data.Load(key)
	if !ok || st == nil {
		st = &storage{
			Name:   key,
			Values: make([]int, 0),
		}
	}

	st.(*storage).Values = append(st.(*storage).Values, data...)
	s.data.Store(key, st)
}

func (s *smap) getData(key string) *storage {
	st, ok := s.data.Load(key)
	if !ok || st == nil {
		return nil
	}

	return st.(*storage)
}

type service struct {
	data      repository
	dataMutex sync.RWMutex
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
		case http.MethodPost:
			s.setData(w, r)
		case http.MethodPut:
			s.updateData(w, r)
		default:
			http.Error(w, fmt.Sprintf("[%s] unsupported method", r.Method), http.StatusMethodNotAllowed)
		}

	case "/":
		makeResponse(w, "Welcome to main")
		return
	}
}

func main() {
	svc := service{
		data:      &smap{},
		dataMutex: sync.RWMutex{},
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

	s.data.setData(rb.Key, rb.Data)

	makeResponse(w, "Success", http.StatusCreated)
}

func (s *service) updateData(w http.ResponseWriter, r *http.Request) {
	var rb requestBody
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	st := s.data.getData(rb.Key)
	if st == nil {
		makeResponse(w, fmt.Sprintf("no data for key [%s]", rb.Key), http.StatusNotFound)
		return
	}

	s.data.setData(rb.Key, rb.Data)

	makeResponse(w, "Success", http.StatusAccepted)
}

func (s *service) getData(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Query().Get("key")

	responseData := s.data.getData(k)

	if responseData == nil {
		makeResponse(w, fmt.Sprintf("no data for key [%s]", k), http.StatusNotFound)
		return
	}

	makeResponse(w, responseData, http.StatusOK)
}

func makeResponse(w http.ResponseWriter, v any, status ...int) {
	w.Header().Set("Content-Type", "application/json")

	if len(status) > 0 {
		w.WriteHeader(status[0])
	}

	if err := json.NewEncoder(w).Encode(map[string]any{
		"payload": v,
	}); err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	makeResponse(w, "Hello there!", http.StatusOK)
}
