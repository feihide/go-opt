package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"example/ex/handler"
	"example/ex/subscriber"

	ex "example/ex/proto/ex"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.ex"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	ex.RegisterExHandler(service.Server(), new(handler.Ex))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.ex", service.Server(), new(subscriber.Ex))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.ex", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
