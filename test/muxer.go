package main

import (
  "dev.io/cloud/utils"
  "net/http"
  "testing"
  "context"
  "time"
  "log"
  "fmt"
  "io"
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

  return srv
}

func stop(srv *http.Server) {
  ctx, cancel := context.WithTimeout(context.Background(),
                                     1 * time.Second)
  defer cancel()

  srv.Shutdown(ctx)
  log.Println("shutting down")
}

func TestStartStopServer(t *testing.T) {
  re := utils.NewApiServer()

  // build a simple api
  re.Version("v1").
    Endpoint("echo").
      Handle("GET",
        func(w http.ResponseWriter, r *http.Request) {
          re.Ok(w)("hello")
        }).
      Mock("/echo")

  // start a new server
  srv := start(re, 1080)

  // wait 10 seconds
  time.Sleep(10 * time.Second)

  // do http request
  if resp, err := http.Get("http://127.0.0.1:1080/echo"); err != nil {
    t.Error("request got error %s", err.Error())
  } else if body, err := io.ReadAll(resp.Body); err != nil {
    t.Error("parsing body got error %s", err.Error())
  } else if fmt.Sprintf("%s", body) != `{"code": 200, "data": "hello"}`{
    t.Error("can't fetch correct data: ", fmt.Sprintf("%s", body))
  }

  if resp, err := http.Get("http://127.0.0.1:1080/v1/echo"); err != nil {
    t.Error("request got error %s", err.Error())
  } else if body, err := io.ReadAll(resp.Body); err != nil {
    t.Error("parsing body got error %s", err.Error())
  } else if fmt.Sprintf("%s", body) != `{"code": 200, "data": "hello"}`{
    t.Error("can't fetch correct data: ", fmt.Sprintf("%s", body))
  }

  if resp, err := http.Get("http://127.0.0.1:1080/echo1"); err != nil {
    t.Error("request got error %s", err.Error())
  } else if body, err := io.ReadAll(resp.Body); err != nil {
    t.Error("parsing body got error %s", err.Error())
  } else if fmt.Sprintf("%s", body) == `{"code": 200, "data": "hello"}`{
    t.Error("can't fetch correct data")
  }

  // stop server grateful
  stop(srv)
}

func main() {}
