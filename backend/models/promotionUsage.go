package models

import "time"

type PromotionUsage struct {
	ID          string    `gorm:"primaryKey"`
	UserID      string    `gorm:"not null"`
	PromotionID string    `gorm:"not null"`
	ConsumedAt  time.Time `gorm:"not null"`
}
