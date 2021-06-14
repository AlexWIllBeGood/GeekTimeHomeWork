package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	g, ctx = errgroup.WithContext(context.Background())
)

func main() {
	g.Go(RunHttpServer)
	g.Go(RunSignalHandleServer)
	log.Println(g.Wait().Error())
}

func RunHttpServer() error {
	server := http.Server{}
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	addr := "0.0.0.0:8080"
	mutex := http.ServeMux{}
	mutex.HandleFunc("/hi", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("say hi"))
	})
	server.Addr = addr
	server.Handler = &mutex

	return server.ListenAndServe()
}

func RunSignalHandleServer() error {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			return errors.New("RunSignalHandleServer stopped!")
		default:
			fmt.Println("other signal", s)
		}
	}
	return nil
}
