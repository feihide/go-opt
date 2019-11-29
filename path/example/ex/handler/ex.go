package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	ex "example/ex/proto/ex"
)

type Ex struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Ex) Call(ctx context.Context, req *ex.Request, rsp *ex.Response) error {
	log.Log("Received Ex.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Ex) Stream(ctx context.Context, req *ex.StreamingRequest, stream ex.Ex_StreamStream) error {
	log.Logf("Received Ex.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&ex.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Ex) PingPong(ctx context.Context, stream ex.Ex_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&ex.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
