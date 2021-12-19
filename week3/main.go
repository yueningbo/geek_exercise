package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"sync/errgroup"
)

func serverStart(srv *http.Server) error {
	return srv.ListenAndServe()
}

type appHandle struct{}
type debugHandle struct{}

func (ah *appHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello App! %s", time.Now())
}

func (dh *debugHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Debug! %s", time.Now())
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	g, errCtx := errgroup.WithContext(ctx)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop)

	appServer := &http.Server{Addr: ":8080", Handler: &appHandle{}}
	debugServer := &http.Server{Addr: ":8081", Handler: &debugHandle{}}

	g.Go(func() error {
		return serverStart(appServer)
	})
	g.Go(func() error {
		return serverStart(debugServer)
	})

	g.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()

			case <-stop:
				cancel()
			}
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group error:", err)
	}
	fmt.Println("All Done")
}
