package main

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

const MAX_REGISTRANTS = 10000

func ensureEventConfig(db *gorm.DB) (EventConfig, error) {
    var cfg EventConfig
    if err := db.Order("id").First(&cfg).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            cfg = EventConfig{RegistrationOpen: true}
            if err := db.Create(&cfg).Error; err != nil {
                return cfg, err
            }
            return cfg, nil
        }
        return cfg, err
    }
    return cfg, nil
}

func getPublicStatusHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        cfg, err := ensureEventConfig(db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var total int64
        db.Model(&Registrant{}).Count(&total)

        // counts by category
        rows := map[string]int64{}
        type CatCount struct{
            Category string
            Count int64
        }
        var results []CatCount
        db.Model(&Registrant{}).Select("category, count(*) as count").Group("category").Scan(&results)
        counts := map[string]int{}
        for _, r := range results {
            counts[r.Category] = int(r.Count)
        }

        c.JSON(http.StatusOK, gin.H{
            "countdown_target_utc": cfg.CountdownTargetUTC,
            "registration_open": cfg.RegistrationOpen,
            "public_message": cfg.PublicMessage,
            "total_registered": total,
            "count_by_category": counts,
        })
    }
}

type RegisterIn struct{
    FullName string `json:"full_name" binding:"required,min=1,max=200"`
    Email *string `json:"email" binding:"omitempty,email"`
    Phone *string `json:"phone" binding:"omitempty,max=64"`
    Company *string `json:"company" binding:"omitempty,max=255"`
    Category string `json:"category" binding:"required"`
    Notes *string `json:"notes" binding:"omitempty"`
}

func postRegisterHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var in RegisterIn
        if err := c.ShouldBindJSON(&in); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        cfg, err := ensureEventConfig(db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if !cfg.RegistrationOpen {
            c.JSON(http.StatusForbidden, gin.H{"detail": "Registration is closed"})
            return
        }
        var total int64
        db.Model(&Registrant{}).Count(&total)
        if int(total) >= MAX_REGISTRANTS {
            c.JSON(http.StatusConflict, gin.H{"detail": "Maximum capacity reached"})
            return
        }

        r := Registrant{
            FullName: strings.TrimSpace(in.FullName),
            Email: in.Email,
            Phone: in.Phone,
            Company: in.Company,
            Category: in.Category,
            Notes: in.Notes,
            CreatedAt: time.Now(),
        }
        if err := db.Create(&r).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, r)
    }
}

func getLookupHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        q := c.Query("q")
        if len(strings.TrimSpace(q)) < 2 {
            c.JSON(http.StatusOK, []Registrant{})
            return
        }
        like := "%" + strings.ToLower(q) + "%"
        var regs []Registrant
        db.Where("lower(full_name) LIKE ?", like).Order("full_name asc").Limit(50).Find(&regs)
        c.JSON(http.StatusOK, regs)
    }
}

func requireHeaderMiddleware(keyName, expected string) gin.HandlerFunc {
    return func(c *gin.Context) {
        v := c.GetHeader(keyName)
        if v == "" || v != expected {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Invalid key"})
            return
        }
        c.Next()
    }
}

func adminGetStatsHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // reuse public status
        cfg, err := ensureEventConfig(db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var total int64
        db.Model(&Registrant{}).Count(&total)
        type CatCount struct{Category string; Count int64}
        var results []CatCount
        db.Model(&Registrant{}).Select("category, count(*) as count").Group("category").Scan(&results)
        counts := map[string]int{}
        for _, r := range results { counts[r.Category] = int(r.Count) }

        c.JSON(http.StatusOK, gin.H{
            "total_registered": total,
            "total_checked_in": 0,
            "count_by_category": counts,
            "countdown_target_utc": cfg.CountdownTargetUTC,
            "registration_open": cfg.RegistrationOpen,
            "public_message": cfg.PublicMessage,
        })
    }
}

type AdminConfigUpdate struct{
    CountdownTargetUTC *time.Time `json:"countdown_target_utc"`
    RegistrationOpen *bool `json:"registration_open"`
    PublicMessage *string `json:"public_message"`
}

func adminPutConfigHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var upd AdminConfigUpdate
        if err := c.ShouldBindJSON(&upd); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        cfg, err := ensureEventConfig(db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if upd.CountdownTargetUTC != nil { cfg.CountdownTargetUTC = upd.CountdownTargetUTC }
        if upd.RegistrationOpen != nil { cfg.RegistrationOpen = *upd.RegistrationOpen }
        if upd.PublicMessage != nil { cfg.PublicMessage = upd.PublicMessage }
        cfg.UpdatedAt = time.Now()
        if err := db.Save(&cfg).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        // return stats after update
        c.JSON(http.StatusOK, gin.H{"ok": true})
    }
}
