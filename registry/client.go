package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	heartbeatURL, err := url.Parse(r.HeartBeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serviceURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	//注册处理patchEntry的方法，再下面服务注册成功后，会发送依赖服务请求来更新prov
	http.Handle(serviceURL.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	// 注册服务到 registry server 中
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to register service. Register service %v responded with %v",
			r.ServiceName, res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

// 每个服务都会注册这个方法，在服务注册成功后，会发送请求将依赖增加或删除更新到prov
func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("update received %v\n", p)
	prov.Update(p)
}

func ShutdownService(url string) error {
	buf := bytes.NewBuffer([]byte(url))
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to shutdown service. Register service %v responded with %v",
			url, res.StatusCode)
	}
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}
	for _, patchEntry := range pat.Removed {
		if providerURLS, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLS {
				if providerURLS[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerURLS[:i], providerURLS[i+1:]...)
				}
			}
		}
	}
}

func (p providers) get(name ServiceName) (string, error) {
	services, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("service name not found %v", name)
	}
	idx := int(rand.Float32() * float32(len(services)))
	return services[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
