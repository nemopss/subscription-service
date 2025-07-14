package models

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	ServiceName string     `json:"service_name" gorm:"not null"`
	Price       int        `json:"price" gorm:"not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"not null"`
	StartDate   time.Time  `json:"start_date" gorm:"not null"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}
