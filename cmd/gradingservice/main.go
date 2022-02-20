package main

import (
	"context"
	"distributed/grades"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	reg := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(context.Background(), host, port, reg, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down GradingService.")
}
