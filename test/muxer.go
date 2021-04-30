package main

import (
  "dev.io/cloud/utils"
  "net/http"
  "testing"
  "time"
  "log"
  "fmt"
)

func start() *http.Server {
  r := utils.NewApiServer()

  srv := &http.Server{
    Addr: "0.0.0.0:8080",
    WriteTimeout: time.Second * 15,
    ReadTimeout: time.Second * 15,
    IdleTimeout: time.Second * 60,
    Handler: r.GetMuxer(),
  }

  go func() {
    if err := srv.ListenAndServe(); err != nil {
      log.Println(err)
    }
  }()
}

func stop(srv *http.Server) {
  ctx, cancel := context.WithTimeout(context.Background(), wait)
  defer cancel()

  srv.Shutdown(ctx)
  log.Println("shutting down")
}

func TestStartStopServer(t *testing.T) {
  // start a new server
  start()

  // wait 10 seconds
  time.Sleep(10 * time.Second)

  // stop server grateful
  stop()
}
