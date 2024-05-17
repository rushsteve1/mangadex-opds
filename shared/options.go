package shared

import (
	"cmp"
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
)

var GlobalOptions Options

type Options struct {
	Bind          string
	Host          url.URL
	Language      string
	Query         url.Values
	DataSaver     bool
	MDUploads     bool
	DevApi        bool
	ExpVars       bool
	GzipResponses bool
	LogLevel      slog.Level
}

var defaultBind = url.URL{
	Scheme: "http",
	Host:   "0.0.0.0:4444",
}

var Version string

func init() {
	info, _ := debug.ReadBuildInfo()

	var rev string
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			rev = s.Value
		}
	}

	if len(rev) > 8 {
		rev = rev[:8]
	}

	Version = cmp.Or(rev, info.Main.Version)
}

func ReadOptionsFromEnv() Options {
	h := env("HOST", "http://localhost:4444")
	u, err := url.Parse(h)
	if err != nil {
		slog.Error(
			"error reading HOST variable, using fallback",
			"error",
			err,
			"fallback",
			defaultBind.String(),
		)
	}
	if u == nil {
		u = &defaultBind
	}

	if u.Scheme == "" {
		slog.Warn("no scheme on HOST, assuming https")
		u.Scheme = "https"
	}

	return Options{
		Bind:          env("BIND", defaultBind.Host),
		Host:          *u,
		Language:      env("LANGUAGE", "en"),
		Query:         env("QUERY", url.Values{}),
		DataSaver:     env("DATA_SAVER", false),
		MDUploads:     env("MD_UPLOADS", false),
		DevApi:        env("DEV_API", false),
		ExpVars:       env("EXP_VARS", true),
		GzipResponses: env("GZIP_RESPONSES", true),
		LogLevel:      env("LOG_LEVEL", slog.LevelInfo),
	}
}

func TestOptions() Options {
	return Options{
		Bind:      defaultBind.Host,
		Host:      defaultBind,
		Language:  "en",
		DataSaver: true,
		DevApi:    true,
		ExpVars:   true,
		LogLevel:  slog.LevelDebug,
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
