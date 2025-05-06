package grpcserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	fmt.Printf("Request: method=%s \n\n", info.FullMethod)
	resp, err := handler(ctx, req)

	fmt.Printf("\n\nResponse: method=%s duration=%s error=%v\n", info.FullMethod, time.Since(start), err)
	return resp, err
}

func MarcoUnaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("Hallo Marco")
	return invoker(ctx, method, req, reply, cc, opts...)
}

func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	fmt.Println("Hallo Marco, wir befinden uns im StreamClientInterceptor, ist Gustav schon durch die Leitung gekommen?")
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"token", "gustav",
	)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := streamer(mdCtx, desc, cc, method, opts...)
	if err != nil {
		fmt.Println("Hier ist ein Fehler aufgetreten im Stream. Keine Ahnung warum und wieso. Wahrscheinlich ist Gustav durch die Leitung gekommen...")
	}
	return stream, nil
}

type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Receive a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}
func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Send a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}
