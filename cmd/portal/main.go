package main

import (
	"context"
	"distributed/log"
	"distributed/portal"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "7000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	reg := registry.Registration{
		ServiceName:      registry.Portal,
		ServiceURL:       serviceAddress,
		RequireServices:  []registry.ServiceName{registry.LogService, registry.GradingService},
		ServiceUpdateURL: serviceAddress + "/services",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		portal.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("log service found: %v\n", logProvider)
		log.SetClientLogger(logProvider, reg.ServiceName)
	} else {
		fmt.Printf("log service not found: %v\n", reg.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("shutting down Portal.")
}
