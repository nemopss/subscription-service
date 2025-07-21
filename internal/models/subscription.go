package models

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ID          int       `json:"id"                 gorm:"primaryKey"`
	ServiceName string    `json:"service_name"       gorm:"not null"`
	Price       int       `json:"price"              gorm:"not null"`
	UserID      uuid.UUID `json:"user_id"            gorm:"not null"`
	StartDate   string    `json:"start_date"         gorm:"not null"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type CreateSubscription struct {
	ServiceName string    `json:"service_name"       gorm:"not null"`
	Price       int       `json:"price"              gorm:"not null"`
	UserID      uuid.UUID `json:"user_id"            gorm:"not null"`
	StartDate   string    `json:"start_date"         gorm:"not null"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error"`
}

type TotalCostResponse struct {
	Total int `json:"total" example:"1000"`
}
