package utils

import (
  "google.golang.org/grpc"
  "errors"
  "fmt"
  "net"
)

type Invent interface {
  // @NOTE: this method is used to get the current client's version
  Version() string

  // @NOTE: this method is used to get socket of this connection
  Socket() int

  // @NOTE: this method is used to create new Invent with specific protoc
  New(conn *grpc.ClientConn) error

  // @NOTE: this event is raised when a connection is atempted with specific
  // protocol
  OnConnecting(protocol string) error

  // @NOTE: this event is raised when a connection is established
  OnConnected(sock int) error

  // @NOTE: this event is raised when connection is broken
  OnBroken(sock int) error

  // @NOTE: this event is raised when connection is closed
  OnDisconnecting()
}

type Implement interface {
  // @NOTE: this method is used to get the current client's version
  Version() string

  // @NOTE: this method is used to init a new listener for this Implement
  Listen(protocol string) (net.Listener, error)

  // @NOTE: this method is used to create new Implement with specific protoc
  New(serv *grpc.Server) error
}

type iGRpcConnection struct {
  connection *grpc.ClientConn
  protocol string
  index int
}

type iGRpcImplement struct {
  // @NOTE: implementer stores the actual implementer which is used to raise
  // events during serving
  implementer Implement

  // @NOTE: protocol stores the current working protocol for this implementer
  protocol string

  // @NOTE: serving stores the grpc server object
  serving *grpc.Server
}

type iGRpcConnectivityBundle struct {
  // @NOTE: newClientInitializer defines a function which is used to generate
  // a type of GRpc connection between client and server and this could be 
  // used along side with specific type of Implementers
  newClientInitializer func(string) (*grpc.ClientConn, error)

  // @NOTE: listenerInitializer defines a function which is used to generate
  // a new listener object which is essential to create a new server
  listenerInitializer func() (net.Listener, error)

  // @NOTE: inventors is a container which stores every inventor of this
  // specific protocol
  inventors []Invent
}

type GRpcContext struct {
  // @NOTE: protocols is a mapping which redirect protocol name to protocol
  // object and let developer to access grpc resource and so on
  protocols map[string]*iGRpcConnectivityBundle

  // @NOTE: connections is an array which stores detail information about
  // each connectivity between client and server
  connections []*iGRpcConnection

  // @NOTE: implemnters is a container which stores every implementers of this
  // specific protocol
  implementers []*iGRpcImplement
}

/*! \brief Connect inventory to implementer
 *
 *  This function is used to connect our inventory to remote implementer
 * before doing anything
 *
 *  \return error: if everything ok, we will receive nil object otherwide we
 *                 will receive error which indicate issue during connecting
 */
func (self *GRpcContext) Connect(invent Invent) error {
  cnt := len(self.protocols)

  if self.protocols == nil {
    initGRpcProtocols(self)
  }

  for name, context := range self.protocols {
    cnt -= 1

    if err := invent.OnConnecting(name); err != nil {
      if cnt > 0 {
        // @TODO: again we must write log here
        continue
      }
      
      return err
    }

    if conn, err := context.newClientInitializer(name); err != nil {
      // @TODO: how to write log here, we would need to log this one for
      // further investigation
    } else if err := invent.New(conn); err != nil {
      if cnt > 0 {
        // @TODO: again we must write log here
        continue
      }
      
      return err
    } else if err := invent.OnConnected(len(self.connections)); err != nil {
      conn.Close()

      if cnt > 0 {
        // @TODO: again we must write log here
        continue
      }
      
      return err
    } else {
      // The connection has been established and we must store this one to
      // our cache to be used later

      self.connections = append(self.connections, &iGRpcConnection{
        connection: conn,
        protocol: name,
        index: len(context.inventors),
      })
      context.inventors = append(context.inventors, invent)
      return nil
    }
  }

  return errors.New("can't establish a new connection recently")
}

/*! \brief Disconnect a connection 
 *
 *  This function is used to disconnect our inventory to remote implementer
 * to close connection gently
 *
 *  \return error: if everything ok, we will receive nil object otherwide we
 *                 will receive error which indicate issue during connecting
 */
func (self *GRpcContext) Disconnect(invent Invent) error {
  sock := invent.Socket()

  if 0 <= sock && sock < len(self.connections) {
    invent.OnDisconnecting()
    self.connections[sock].connection.Close()

    // @NOTE: the protocol name can't cause any corruption here, but if it
    // crash, we could see the console log here which indicate an issue 
    // buffer-overload elsewhere
    protocol := self.protocols[self.connections[sock].protocol]
    index := self.connections[sock].index

    copy(protocol.inventors[index:], protocol.inventors[index + 1:])
    copy(self.connections[sock:], self.connections[sock + 1:])
    return nil
  }

  return errors.New("disconnect an disconnected invent")
}

/*! \brief Serve an implementer to resolve requests
 *
 *  This function is used to start on-board our implementer to serve requests
 * from clients
 *
 *  \return error: if everything ok, we will receive nil object otherwide we
 *                 will receive error which indicate issue during connecting
 */
func (self *GRpcContext) Serve(imp Implement) error {
  cnt := 0

  for name, context := range self.protocols {
    cnt += 1

    if listener, err := imp.Listen(name); err != nil {
      return err
    } else {
      if listener != nil && context.listenerInitializer != nil {
        listener = context.listenerInitializer();
      } else if cnt < len(self.protocols) {
        continue
      } else {
        return errors.New("Can't serve this implement")
      }

      serving := grpc.NewServer()
      index := len(self.implementers)

      if err = imp.New(serving); err != nil {
        return err
      }

      self.implementers = append(self.implementers, &iGRpcImplement{
        implementer: imp,
        protocol: name,
        serving: serving,
      })

      defer func() {
        copy(self.implementers[index:], self.implementers[index + 1:])
      }()

      if err = serving.Serve(listener); err != nil {
        return err
      } else {
        return nil
      }
    }
  }

  return errors.New("can't serve this Implement")
}

/*! \brief 
 *
 *  This function is used to create an new GRpcContext which is used to store
 * everything 
 *
 *  \return net.Listener: if everything ok, we will receive a new GRpcContext
 *                        pointer
 */
func (self *GRpcContext) MakeListener(protocol string) (net.Listener, error) {
  if context, ok := self.protocols[protocol]; ! ok {
    return nil, errors.New(fmt.Sprintf("don't support %s", protocol))
  } else if context.listenerInitializer == nil {
    return nil, errors.New(fmt.Sprintf("%s's listener initializer is nil"))
  } else {
    listener, err := context.listenerInitializer()
    return listener, err
  }
}

/*! \brief 
 *
 *  This function is used to create an new GRpcContext which is used to store
 * everything 
 *
 *  \return *GRpcContext: if everything ok, we will receive a new GRpcContext
 *                        pointer
 */
func NewGRpcContext() *GRpcContext {
  return &GRpcContext{}
}

/*! \brief Init grpc's protocols
 *
 *  This function is used to init grpc's protocols, base on expected user
 *
 */
func initGRpcProtocols(ctx *GRpcContext) {
  // @NOTE: by default we will support tcp and ipc because it's default
  // grpc protocols

  ctx.protocols = make(map[string]*iGRpcConnectivityBundle)

  initGRpcTcpProtocol(ctx)
  initGRpcIpcProtocol(ctx)
  initGRpcTipcProtocol(ctx)
  initGRpcSctpProtocol(ctx)
  initGRpcQuicProtocol(ctx)
}

/*! \brief Init tcp protocol
 *
 *  This function is used to init grpc's tcp protocol
 *
 */
func initGRpcTcpProtocol(ctx *GRpcContext) {
  listenerInitializer := func() (net.Listener, error) {
    if lis, err := net.Listen("tcp", "localhost:50051"); err != nil {
      return nil, err
    } else {
      return lis, nil
    }
  }
  
  newClientInitializer := func(address string) (*grpc.ClientConn, error) {
    return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
  }

  ctx.protocols["tcp"] = &iGRpcConnectivityBundle{
    newClientInitializer: newClientInitializer,
    listenerInitializer: listenerInitializer,
    inventors: make([]Invent, 0),
  }
}

/*! \brief Init ipc protocol
 *
 *  This function is used to init grpc's ipc protocol
 *
 */
func initGRpcIpcProtocol(ctx *GRpcContext) {
}

/*! \brief Init tipc protocol
 *
 *  This function is used to init grpc's tipc protocol
 *
 */
func initGRpcTipcProtocol(ctx *GRpcContext) {
}

/*! \brief Init sctp protocol
 *
 *  This function is used to init grpc's sctp protocol
 *
 */
func initGRpcSctpProtocol(ctx *GRpcContext) {
}

/*! \brief Init quic protocol
 *
 *  This function is used to init grpc's quic protocol
 *
 */
func initGRpcQuicProtocol(ctx *GRpcContext) {
}
