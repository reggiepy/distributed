package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	ServerPort  = ":3000"
	ServicesURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mux           *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mux.Lock()
	r.registrations = append(r.registrations, reg)
	r.mux.Unlock()
	// 再 registrations 中查找 require service 兵 发送
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			patchEntry{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return err
}

func (r *registry) notify(fullPatch patch) {
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequireServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, add := range fullPatch.Added {
					if reqService == add.Name {
						p.Added = append(p.Added, add)
						sendUpdate = true
					}
				}
				for _, remove := range fullPatch.Removed {
					if reqService == remove.Name {
						p.Removed = append(p.Removed, remove)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequireServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	for index := range r.registrations {
		if r.registrations[index].ServiceURL == url {
			r.notify(patch{
				Removed: []patchEntry{
					patchEntry{
						Name: r.registrations[index].ServiceName,
						URL:  r.registrations[index].ServiceURL,
					},
				},
			})
			r.mux.Lock()
			r.registrations = append(r.registrations[:index], r.registrations[index+1:]...)
			r.mux.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service url %s not found", url)
}

func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				for attemps := 0; attemps < 3; attemps++ {
					res, err := http.Get(reg.HeartBeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("heartbeat check passed for %v", reg.ServiceName)
						if !success {
							r.add(reg)
						}
						break
					}
					log.Printf("heartbeat check failed for %v", reg.ServiceName)
					if success {
						success = false
						r.remove(reg.ServiceURL)
					}
				}
				time.Sleep(1 * time.Second)
			}(reg)
		}
		wg.Wait()
		time.Sleep(freq)
	}
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}

var reg = registry{
	registrations: make([]Registration, 0),
	mux:           new(sync.RWMutex),
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
