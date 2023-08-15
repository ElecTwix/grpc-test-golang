package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/ElecTwix/grpctest/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	greeter.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	fmt.Println(req.Name)
	return &greeter.HelloReply{
		Message:      "Hello, " + req.Name,
		NewSomething: &greeter.Nested{Wtf: "hello"},
	}, nil
}

func runServer() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	srv := grpc.NewServer()
	greeter.RegisterGreeterServer(srv, &server{})
	reflection.Register(srv)

	fmt.Println("Server listening on port 50051")
	if err := srv.Serve(listen); err != nil {
		fmt.Println("Failed to serve:", err)
	}
}

func runClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	client := greeter.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &greeter.HelloRequest{Name: "Wtf"})
	if err != nil {
		fmt.Println("Failed to call SayHello:", err)
		return
	}

	fmt.Println("Response:", resp)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./program <server|client>")
		return
	}

	switch os.Args[1] {
	case "server":
		runServer()
	case "client":
		runClient()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
