package models

import "time"

type Visit struct {
	ID        int32     `gorm:"autoIncrement:true"`
	IP        string    `gorm:"not null"`
	Time      time.Time `gorm:"not null"`
	UserAgent *string
	ShortURL  string `gorm:"not null"`
}
