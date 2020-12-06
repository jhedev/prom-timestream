package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/jhedev/prom-timestream/server"
)

func main() {
	var (
		logger = log.New().WithField("component", "main")

		addr = flag.String("addr", ":4000", "")
	)
	flag.Parse()

	srv, err := server.New(nil)
	if err != nil {
		logger.Fatalf("error while creating server")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/write", srv.Write)
	mux.HandleFunc("/read", srv.Read)

	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	server.ListenAndServe()
}
