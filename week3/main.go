package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func server(addr string, handler http.Handler, stop <-chan os.Signal) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func serverApp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello App! %s", time.Now())
}

func serverDebug(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello debug! %s", time.Now())
}

type appHandle struct{}

func (app *appHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello App! %s", time.Now())
}

type debugHandle struct{}

func (debug *debugHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Debug! %s", time.Now())
}

func main() {
	var g errgroup.Group

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGQUIT)

	g.Go(func() error {
		return server(":8080", &appHandle{}, stop)
	})
	g.Go(func() error {
		return server(":8081", &debugHandle{}, stop)
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
