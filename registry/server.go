package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort  = ":3000"
	ServicesURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mux           *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.registrations = append(r.registrations, reg)
	return nil
}

func (r *registry) remove(url string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	for index, reg := range r.registrations {
		if reg.ServiceURL == url {
			r.registrations = append(r.registrations[:index], r.registrations[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("service url %s not found", url)
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
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing services: URL %v", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
