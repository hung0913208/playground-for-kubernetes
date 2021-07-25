package main

import (
  pb "dev.io/cloud/protoc"
  srv "dev.io/cloud/utils"
  grpc "google.golang.org/grpc"

  "testing"
  "errors"
  "fmt"
  "net"
)

type sample struct {
  listener net.Listener
  server *grpc.Server
}

func (self *sample) Version() string {
  return "v1"
}

func (self *sample) New(srv *grpc.Server) error {
  self.server = srv

  pb.RegisterGatewayServiceServer(srv, self)
  return nil
}

func (self *sample) Listen(protocol string) (net.Listener, error) {
  if protocol != "tcp" {
    return nil, errors.New(fmt.Sprintf("don't support %s", protocol))
  }

  if lis, err := net.Listen("tcp", "localhost:50051"); err != nil {
    return nil, err
  } else {
    self.listener = lis
  }

  return self.listener, nil
}

func TestStartStopServer(t *testing.T) {
  grpc.EnableTracing = true

  ctx := srv.NewGRpcContext()
  smp := &sample{}

  go func() {
    ctx.Serve(smp)
  }()
}
