package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rushsteve1/mangadex-opds/server"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func main() {
	var err error
	shared.GlobalOptions, err = shared.ReadOptionsFromEnv()
	if err != nil {
		slog.Error("error parsing options", "error", err.Error())
	}

	srv := http.Server{
		Addr:              fmt.Sprintf(":%d", 4444),
		ReadHeaderTimeout: time.Second * 5,
		Handler:           server.Router(),
	}

	slog.Info("starting server", "addr", srv.Addr)

	err = srv.ListenAndServe()
	slog.Error(err.Error())
}
