package model

import (
	"time"
)

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}

type Business struct {
	ID             int64     `json:"id"`
	OwnerID        int64     `json:"owner_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Description    *string   `json:"description,omitempty"`
	Address        string    `json:"address,omitempty"`
	Latitude       *float64  `json:"latitude,omitempty"`
	Longitude      *float64  `json:"longitude,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	LogoURL        *string   `json:"logo_url,omitempty"`
	AvgServiceTime int       `json:"avg_service_time"`
	MaxQueueSize   int       `json:"max_queue_size"`
	IsActive       bool      `json:"is_active"`
	Rating         float64   `json:"rating,omitempty"`
	TotalReviews   int       `json:"total_reviews,omitempty"`
	OpeningTime    *string   `json:"opening_time,omitempty"`
	ClosingTime    *string   `json:"closing_time,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

type Queue struct {
	ID                int64      `json:"id"`
	BusinessID        int64      `json:"business_id"`
	CustomerID        int64      `json:"customer_id"`
	QueueNumber       string     `json:"queue_number"`
	Status            string     `json:"status"`
	PartySize         int64      `json:"party_size"`
	Position          int64      `json:"position"`
	EstimatedWaitTime *int       `json:"estimated_wait_time,omitempty"`
	JoinedAt          *time.Time `json:"joined_at,omitempty"`
	CalledAt          *time.Time `json:"called_at,omitempty"`
	CompletedAt       *time.Time ` json:"completed_at,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time ` json:"updated_at,omitempty"`
}
