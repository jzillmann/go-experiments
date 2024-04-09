package main

import (
	"context"
	"fmt"
	"go-server/proto"
	"log/slog"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
  port=8080
)

type server struct {
	proto.ServiceServer
}

func (s *server) CurrentTime(ctx context.Context, in *proto.Empty) (*timestamppb.Timestamp, error) {
    return timestamppb.Now(), nil
}

func (s *server) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingReply, error) {
    return &proto.PingReply{Message: "Reply: "+ strings.ToUpper(in.Message)}, nil
}

func main() {
  lis,err :=  net.Listen("tcp", fmt.Sprintf(":%v",port))

  if err != nil {
		slog.Error("Failed to start the server", "port", port, "err",err)
		panic(err)
	}

  slog.Info("TCP Server started","port",port)

  
  var grpcServer = grpc.NewServer()
  proto.RegisterServiceServer(grpcServer, &server{})

  // routine
  go func() {
    slog.Info("gRPC server initialized")
    err = grpcServer.Serve(lis)
		if err != nil {
      slog.Error("Error serving gRPC", "s.Serve", err)
			panic(err)
		}
	}()

  // wait on routine
  select {}
}