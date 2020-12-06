package main

import (
	"flag"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"

	"github.com/jhedev/prom-timestream/adapter"
	"github.com/jhedev/prom-timestream/server"
)

func main() {
	var (
		logger = log.New().WithField("component", "main")

		addr         = flag.String("addr", ":4000", "")
		databaseName = flag.String("database-name", "prom", "The database name to use in timestream")
		tableName    = flag.String("table-name", "metrics", "The table name to use in timestream")
	)
	flag.Parse()

	logger.Infof("setting up AWS session...")
	sess, err := session.NewSession()
	if err != nil {
		logger.Fatalf("error while creating aws session: %s", err)
	}

	logger.Infof("setting up timestream adapter...")
	adapt, err := adapter.New(*databaseName, *tableName, sess)
	if err != nil {
		logger.Fatalf("error while creating adapter: %s", err)
	}

	logger.Infof("setting up server...")
	srv, err := server.New(adapt)
	if err != nil {
		logger.Fatalf("error while creating server: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/write", srv.Write)
	mux.HandleFunc("/read", srv.Read)

	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	logger.Infof("listening on %s...", *addr)
	server.ListenAndServe()
}
