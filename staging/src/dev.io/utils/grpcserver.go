package utils

import (
  "google.golang.org/grpc"
  "net"
)

type Inventory interface {
  Version() string
  New(conn *grpc.ClientConn) error
}

type Implementer interface {
  Listen() (net.Listener, error)
  Version() string
  New(serv *grpc.Server) error
}

/*! \brief Connect inventory to implementer
 *
 *  This function is used to connect our inventory to remote implementer
 * before doing anything
 *
 *  \return error: if everything ok, we will receive nil object otherwide we
 *                 will receive error which indicate issue during connecting
 */
func (self Inventory) Connect(conn *grpc.ClientConn) error {
  return self.New(conn)
}

/*! \brief Serve an implementer to resolve requests
 *
 *  This function is used to start on-board our implementer to serve requests
 * from clients
 *
 *  \return error: if everything ok, we will receive nil object otherwide we
 *                 will receive error which indicate issue during connecting
 */
func (self Implementer) Serve() error {
  if lis, err := self.Listen(); err != nil {
    return err
  } else {
    serv := grpc.NewServer()

    if err := self.New(serv); err != nil {
      return err
    } else if err := serv.Serve(lis); err != nil {
      return err
    } else {
      return nil
    }
  }
}
