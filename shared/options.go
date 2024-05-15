package shared

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

var GlobalOptions Options

type Options struct {
	Bind      string
	Host      url.URL
	Language  string
	Query     url.Values
	DataSaver bool
	MDUploads bool
	DevApi    bool
	ExpVars   bool
}

var defaultHost = url.URL{
	Scheme: "http",
	Host:   "localhost:4444",
}

func ReadOptionsFromEnv() Options {
	h := env("HOST", "http://localhost:4444")
	u, err := url.Parse(h)
	if err != nil {
		slog.Error("error reading host variable", "error", err)
	}
	if u == nil {
		u = &defaultHost
	}

	return Options{
		Bind:      env("ADDRESS", defaultHost.Host),
		Host:      *u,
		Language:  env("LANGUAGE", "en"),
		Query:     env("QUERY", url.Values{}),
		DataSaver: env("DATA_SAVER", false),
		MDUploads: env("MD_UPLOADS", false),
		DevApi:    env("DEV_API", false),
		ExpVars:   env("EXP_VARS", true),
	}
}

func TestOptions() Options {
	return Options{
		Bind:      defaultHost.Host,
		Host:      defaultHost,
		Language:  "en",
		DataSaver: true,
		DevApi:    true,
		ExpVars:   true,
	}
}

func env[T any](key string, def T) (out T) {
	e := os.Getenv(key)
	if len(e) == 0 {
		return def
	}

	// TODO this is kinda annoying but also useful
	err := json.Unmarshal([]byte(e), &out)
	if err != nil {
		slog.Error("config error", "key", key, "error", err)
		os.Exit(1)
		return out
	}

	return out
}

func LoadDotEnv() {
	envFile, err := os.ReadFile(".env")
	if err != nil {
		slog.Warn("could not load .env file, ignoring")
	}

	lines := strings.Split(string(envFile), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			slog.Warn("not enough parts in env string", "line number", i+1)
			continue
		}

		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
}
