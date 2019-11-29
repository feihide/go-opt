package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	ex "example/ex/proto/ex"
)

type Ex struct{}

func (e *Ex) Handle(ctx context.Context, msg *ex.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *ex.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
