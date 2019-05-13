package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"imup/server"
)

// https://github.com/swaggo/swag#general-api-info
// @title Image Upload API
// @description upload image, then make thumbnail 100x100, returns image id
// @version 1.0

func main() {
	listen := flag.String("listen", "localhost:3000", "HTTP server listen address and port")
	storage, _ := filepath.Abs(".")
	flag.StringVar(&storage, "storage", storage, "image upload folder")
	thumbcmd := flag.String("thumbcmd", "/usr/local/bin/thumb100.sh", "thumbnailer ($1 - in, $2 - out)")
	swagurl := flag.String("swagurl", "http://localhost:3000/swagger", "swagger visible url")
	flag.Parse()

	server := server.New(*listen, storage, *thumbcmd, *swagurl)

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
