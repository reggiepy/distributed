package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort  = ":3000"
	ServicesURL = "http://localhost:" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mux           *sync.Mutex
}

func (r registry) add(reg Registration) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.registrations = append(r.registrations, reg)
	return nil
}

var reg = registry{
	registrations: make([]Registration, 0),
	mux:           new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println("Error decoding", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding services: %v with URL : %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
