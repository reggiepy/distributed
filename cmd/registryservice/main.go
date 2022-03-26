package main

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	go registry.SetupRegistryService()
	http.Handle("/services", &registry.RegistryService{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()
	go func() {
		fmt.Println("Registry service started. Press any key to stop")
		var s string
		_, _ = fmt.Scanln(&s)
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		select {
		case <-signalChan:
			err := srv.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}
			cancel()
		}
	}()
	<-ctx.Done()
	fmt.Println("Registry service shutdown completed.")
}
