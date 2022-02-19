package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		return err
	}
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
