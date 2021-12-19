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

func server(ctx context.Context, addr string, handler http.Handler, stop <-chan os.Signal) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stop:
		fmt.Println("shutdown~~")
		s.Shutdown(context.Background())
		return ctx.Err()
	default:
		return s.ListenAndServe()
	}
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
	g, ctx := errgroup.WithContext(context.Background())

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGQUIT, syscall.SIGTERM)

	g.Go(func() error {
		return server(ctx, ":8080", &appHandle{}, stop)
	})
	g.Go(func() error {
		return server(ctx, ":8081", &debugHandle{}, stop)
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
