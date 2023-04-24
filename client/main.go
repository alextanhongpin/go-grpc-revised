package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/alextanhongpin/go-grpc-revised/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	address = "localhost:50051"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "The name of the person to greet")
	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		st := status.Convert(err)
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				fmt.Println("Oops! Your request was rejected by the server")
				for _, violation := range t.GetFieldViolations() {
					fmt.Printf("The %q field was wrong:\n", violation.GetField())
					fmt.Printf("\t%s\n", violation.GetDescription())
				}
			case *errdetails.ErrorInfo:
				fmt.Println("domain:", t.GetDomain())
				fmt.Println("reason:", t.GetReason())
				fmt.Println("metadata:", t.GetMetadata())
			case *errdetails.RetryInfo:
				fmt.Println("retry delay:", t.GetRetryDelay())
			}
		}

		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
