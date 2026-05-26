package main

import "time"

type EventConfig struct {
    ID                 uint       `gorm:"primaryKey" json:"id"`
    CountdownTargetUTC *time.Time `json:"countdown_target_utc"`
    RegistrationOpen   bool       `json:"registration_open" gorm:"default:true"`
    PublicMessage      *string    `json:"public_message"`
    UpdatedAt          time.Time  `json:"updated_at"`
}

type Registrant struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    FullName  string     `json:"full_name" gorm:"size:200;not null;index:idx_full_name"`
    Email     *string    `json:"email" gorm:"size:255"`
    Phone     *string    `json:"phone" gorm:"size:64"`
    Company   *string    `json:"company" gorm:"size:255"`
    Category  string     `json:"category" gorm:"size:32;index"`
    Notes     *string    `json:"notes" gorm:"type:text"`
    CheckedIn bool       `json:"checked_in" gorm:"default:false"`
    CreatedAt time.Time  `json:"created_at"`
}
