package main

import (
  pb "dev.io/cloud/protoc"
  srv "dev.io/cloud/utils"
  grpc "google.golang.org/grpc"

  "testing"
  "net"
)

type sample struct {
  listener net.Listener
  server *grpc.Server
}

type (self *sample) Version() string {
  return "v1"
}

type (self *sample) New(srv *grpc.Server) error {
  self.server = srv

  pb.RegisterGatewayServiceServer(srv), self)
  return nil
}

type (self *sample) Listen() (net.Listener, error) {
  if lis, err := net.Listen("tcp", "localhost:50051"); err != nil {
    return nil, err
  } else {
    self.listener = lis
  }

  return self.listener, nil
}

func TestStartStopServer(t *testing.T) {
  smp := &sample{}

  go func() {
    smp.Serve()
  }()
}
