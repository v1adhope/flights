package httpsrv

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*http.Server
	shutdownTimeout time.Duration
}

func New(h http.Handler, opts ...Option) *Server {
	cfg := config(opts...)

	return &Server{
		Server: &http.Server{
			Handler:      h,
			Addr:         cfg.Socket,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
		},
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			log.Fatalf("httpsrv: httpsrv: listenAndServe: %s", err)
		}
	}()

	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("httpsrv: httpsrv: gracefulShutdown: shutdown: %s", err)
	}

	select {
	case <-ctx.Done():
		log.Printf("timeout of %s seconds", s.shutdownTimeout.String())
	}

	log.Print("server exiting")
}
