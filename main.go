package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	template "github.com/mono0x/hero-example/template"
)

//go:generate hero -source=./template -pkgname=template

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		users := []string{
			"Cinnamon",
			"Kitty",
			"Mell",
		}
		template.RenderIndex(users, w)
	})

	server := http.Server{Handler: mux}

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		return err
	}

	go func() {
		if err := server.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)

	for {
		s := <-signalChan
		if s == syscall.SIGTERM {
			if err := server.Shutdown(context.Background()); err != nil {
				return err
			}
			return nil
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
