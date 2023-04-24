package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/alextanhongpin/go-grpc-revised/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())

	name := in.GetName()
	// Reference: https://jbrandhorst.com/post/grpc-errors/
	if name == "" {
		// If you don't need metadata.
		//err := status.Error(codes.InvalidArgument, "name is required")
		st := status.New(codes.InvalidArgument, "name is required")
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: "name is required",
		}
		br := new(errdetails.BadRequest)
		br.FieldViolations = append(br.FieldViolations, v)
		st, err := st.WithDetails(br)
		if err != nil {
			// If this errored, it will always error here.
			panic(fmt.Errorf("unexpected error attaching metadata: %v", err))
		}

		return nil, st.Err()
	}

	if name == "special" {
		st := status.New(codes.InvalidArgument, "name is not allowed")
		ei := &errdetails.ErrorInfo{
			Reason: "denylist_name",
			Domain: "greetersvc.myapis.com",
			Metadata: map[string]string{
				"Name": name,
			},
		}

		st, err := st.WithDetails(ei)
		if err != nil {
			// If this errored, it will always error here.
			panic(fmt.Errorf("unexpected error attaching metadata: %v", err))
		}
		return nil, st.Err()
	}

	if name == "ddos" {
		st := status.New(codes.ResourceExhausted, "too many requests")
		ei := &errdetails.RetryInfo{
			RetryDelay: durationpb.New(30 * time.Second),
		}

		st, err := st.WithDetails(ei)
		if err != nil {
			// If this errored, it will always error here.
			panic(fmt.Errorf("unexpected error attaching metadata: %v", err))
		}

		return nil, st.Err()
	}

	if name == "unk" {
		// This will be treated as internal server error.
		return nil, errors.New("unknown error")
	}

	return &pb.HelloReply{
		Message: "Hello " + in.GetName(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("listening at port *:%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
