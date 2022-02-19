package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, serviceName string, host, port string, registryHandlerFunc func()) (context.Context, error) {
	registryHandlerFunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName string, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()
	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", serviceName)
		var s string
		_, _ = fmt.Scanln(&s)
		_ = srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
