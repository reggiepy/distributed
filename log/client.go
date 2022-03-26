package log

import (
	"bytes"
	"distributed/registry"
	"fmt"
	stlog "log"
	"net/http"
	"unsafe"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type clientLogger struct {
	url string
}

func (cl *clientLogger) Write(p []byte) (int, error) {
	// [91 71 114 97 100 105 110 103 83 101 114 118 105 99 101 93 32 45 32 103 101 116 32 111 110 101 10][GradingService] - get one
	//[91 71 114 97 100 105 110 103 83 101 114 118 105 99 101 93 32 45 32 103 101 116 32 97 108 108 10][GradingService] - get all
	// 最后一位 10 为换行符，切掉就不会多一个换行符了
	//b := bytes.NewBuffer([]byte(p[:len(p)-1]))
	b := bytes.NewBuffer([]byte(bytes.TrimSpace(p)))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log response: %v", res.Status)
	}
	return len(p), nil
}
