package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	ctx    context.Context
	server *http.Server
}

func New(ctx context.Context, address string, handler http.Handler) *HttpServer {
	return &HttpServer{
		ctx: ctx,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", address),
			Handler: handler,
		},
	}
}

func (s *HttpServer) WithContext(ctx context.Context) *HttpServer {
	s.ctx = ctx
	return s
}

func (s *HttpServer) Run() {
	// Server context with cancellation
	srvCtx, srvCancel := context.WithCancel(s.ctx)

	// Listen for syscall signals for the process to interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go s.handleInterruption(sig, srvCtx, srvCancel)

	log.Printf("server is ready for connections, listening on port: %v\n", s.server.Addr)

	// Run server
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	// Wait for server context to be stopped
	<-srvCtx.Done()
}

func (s *HttpServer) handleInterruption(
	sig <-chan os.Signal,
	srvCtx context.Context,
	srvCancel context.CancelFunc,
) {
	func() {
		<-sig

		shutdownCtx, cancelCtx := context.WithTimeout(srvCtx, 30*time.Second)
		defer cancelCtx()

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out... forcing exit.")
			}
		}()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		log.Println("graceful shutdown...")

		srvCancel()
	}()
}
