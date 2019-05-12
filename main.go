package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"imup/server"
)

func main() {
	listen := flag.String("listen", "localhost:3000", "HTTP server listen address and port")
	flag.Parse()

	server := server.New(*listen)

	go func() {
		sigINT := make(chan os.Signal, 2)
		signal.Notify(sigINT, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigINT:
			if err := server.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
}
