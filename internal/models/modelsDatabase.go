package models

import (
	"time"
)

// Структура для тендера
type Tender struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
	Status      string `json:"status"`
	CreatedBy   string `json:"creatorId,omitempty"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}

// Структура для организации
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Структуда сотрудника
type Employee struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
