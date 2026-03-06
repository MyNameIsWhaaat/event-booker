package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
	PGDSN    string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		PGDSN: firstNonEmpty(
			os.Getenv("PG_DSN"),
			os.Getenv("DB_URL"),
			"",
		),
	}

	if cfg.PGDSN == "" {
		log.Fatal("missing required env: PG_DSN (or DB_URL)")
	}

	return cfg
}

func getenv(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
