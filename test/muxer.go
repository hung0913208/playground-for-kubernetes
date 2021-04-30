package main

import (
  "dev.io/cloud/utils"
  "net/http"
  "testing"
  "time"
  "log"
  "fmt"
)

func start(r *utils.ApiServer, port int) *http.Server {
  srv := &http.Server{
    Addr: fmt.Sprintf("0.0.0.0:%d", port),
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
  r := utils.NewApiServer()

  // build a simple api
  r.version("v1").
    endpoint("echo").
      handle("GET",
      func(w *http.ResponseWriter, r *http.Request) {
        r.ok(w)("hello")
      }).
      mock("/echo")

  // start a new server
  start(r, 1080)

  // wait 10 seconds
  time.Sleep(10 * time.Second)

  // do http request
  if _, err := http.Get("http://127.0.0.1:1080/echo"); err != nil {
    t.Error("request got error %s", err.Error())
  }

  if _, err := http.Get("http://127.0.0.1:1080/echo1"); err == nil {
    t.Error("request can't produce error")
  }

  // stop server grateful
  stop()
}
