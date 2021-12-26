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

	signal.Notify(stop, syscall.SIGQUIT, syscall.SIGTERM)

	appServer := &http.Server{Addr: ":8080", Handler: &appHandle{}}
	debugServer := &http.Server{Addr: ":8081", Handler: &debugHandle{}}

	g.Go(func() error {
		return appServer.ListenAndServe()
	})
	g.Go(func() error {
		return debugServer.ListenAndServe()
	})

	g.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()

			case <-stop:
				appServer.Shutdown(errCtx)
				debugServer.Shutdown(errCtx)
				cancel()
			}
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group error:", err)
	}
	fmt.Println("All Done")
}
