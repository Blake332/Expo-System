package main

import (
    "os"
    "strings"
)

type Config struct {
    DeskAPIKey  string
    AdminAPIKey string
    DatabaseURL string
}

func LoadConfig() Config {
    cfg := Config{
        DeskAPIKey:  getEnv("DESK_API_KEY", "dev-desk-key"),
        AdminAPIKey: getEnv("ADMIN_API_KEY", "dev-admin-key"),
        DatabaseURL: getEnv("DATABASE_URL", "sqlite:///./expo_registration.db"),
    }
    return cfg
}

func getEnv(k, d string) string {
    v := os.Getenv(k)
    if v == "" {
        return d
    }
    return v
}

// NormalizeDatabasePath extracts a file path when DATABASE_URL is a sqlite URL
// e.g. sqlite+aiosqlite:///./expo_registration.db -> ./expo_registration.db
func NormalizeDatabasePath(databaseURL string) string {
    if strings.Contains(databaseURL, "sqlite") {
        // find last ":///" and return remainder
        if idx := strings.LastIndex(databaseURL, "///"); idx != -1 && idx+3 < len(databaseURL) {
            return databaseURL[idx+3:]
        }
        // fallback: return as-is
        return databaseURL
    }
    return databaseURL
}
