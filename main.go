package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rushsteve1/mangadex-opds/server"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func main() {
	shared.LoadDotEnv()
	shared.GlobalOptions = shared.ReadOptionsFromEnv()

	slog.SetLogLoggerLevel(shared.GlobalOptions.LogLevel)

	srv := http.Server{
		Addr:              shared.GlobalOptions.Bind,
		ReadHeaderTimeout: time.Second * 5,
		Handler:           server.Router(),
	}

	slog.Info("starting server", "addr", srv.Addr, "log level", shared.GlobalOptions.LogLevel.String())

	err := srv.ListenAndServe()
	slog.Error(err.Error())
}
