package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rushsteve1/mangadex-opds/server"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func main() {
	// This sets the shared.GlobalOptions variable
	shared.ReadOptionsFromEnv()

	srv := http.Server{
		Addr:              shared.GlobalOptions.Bind,
		ReadHeaderTimeout: time.Second * 5,
		Handler:           server.Router(),
	}

	slog.Info(
		"starting server",
		"addr", srv.Addr,
		"host", shared.GlobalOptions.Host.String(),
		"log level", shared.GlobalOptions.LogLevel.String(),
	)

	err := srv.ListenAndServe()
	slog.Error(err.Error())
}
