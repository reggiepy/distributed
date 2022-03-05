package log

import (
	"bytes"
	"distributed/registry"
	"fmt"
	stlog "log"
	"net/http"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl *clientLogger) Write(p []byte) (int, error) {
	// [91 71 114 97 100 105 110 103 83 101 114 118 105 99 101 93 32 45 32 103 101 116 32 111 110 101 10][GradingService] - get one
	//[91 71 114 97 100 105 110 103 83 101 114 118 105 99 101 93 32 45 32 103 101 116 32 97 108 108 10][GradingService] - get all
	// 最后一位 10 为换行符，切掉就不会多一个换行符了
	b := bytes.NewBuffer([]byte(p[:len(p)-1]))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log response: %v", res.Status)
	}
	return len(p), nil
}
