package main

import (
    "log"
    "net/http"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    cfg := LoadConfig()
    dbPath := NormalizeDatabasePath(cfg.DatabaseURL)

    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }

    if err := db.AutoMigrate(&EventConfig{}, &Registrant{}); err != nil {
        log.Fatalf("auto-migrate failed: %v", err)
    }

    r := gin.Default()

    // Basic CORS allow all (matches previous FastAPI setup)
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    })

    // Serve static files from the repo's static directory
    staticDir := filepath.Join("..", "static")
    r.Static("/static", staticDir)
    r.GET("/", func(c *gin.Context) { c.File(filepath.Join(staticDir, "public.html")) })
    r.GET("/desk", func(c *gin.Context) { c.File(filepath.Join(staticDir, "desk.html")) })
    r.GET("/admin", func(c *gin.Context) { c.File(filepath.Join(staticDir, "admin.html")) })

    api := r.Group("/api")

    // Public routes
    api.GET("/public/status", getPublicStatusHandler(db))

    // Desk routes require desk key
    desk := api.Group("/desk")
    desk.Use(requireHeaderMiddleware("X-Desk-Key", cfg.DeskAPIKey))
    desk.POST("/register", postRegisterHandler(db))
    desk.GET("/lookup", getLookupHandler(db))

    // Admin routes
    admin := api.Group("/admin")
    admin.Use(requireHeaderMiddleware("X-Admin-Key", cfg.AdminAPIKey))
    admin.GET("/stats", adminGetStatsHandler(db))
    admin.PUT("/config", adminPutConfigHandler(db))
    admin.GET("/status", adminGetStatsHandler(db))

    log.Println("Starting Go backend on :8080 (serve static from ../static)")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
