package config

import (
	"log"
	"strings"

	wconfig "github.com/wb-go/wbf/config"
)

type Config struct {
	HTTPAddr string
	PGDSN    string
}

func Load() Config {
	cfg := wconfig.New()
	_ = cfg.LoadEnvFiles(".env")
	cfg.EnableEnv("")

	cfg.SetDefault("HTTP_ADDR", ":8080")

	httpAddr := cfg.GetString("HTTP_ADDR")
	pgDSN := firstNonEmpty(
		cfg.GetString("PG_DSN"),
		cfg.GetString("DB_URL"),
	)

	if pgDSN == "" {
		log.Fatal("missing required env: PG_DSN (or DB_URL)")
	}

	return Config{
		HTTPAddr: strings.TrimSpace(httpAddr),
		PGDSN:    strings.TrimSpace(pgDSN),
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
