package models

import (
	"time"
)

// Структура для добавления тендера
type Tender struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ServiceType    string `json:"serviceType"`
	Status         string `json:"status"`
	CreatedBy      string `json:"creatorSlug,omitempty"`
	OrganizationId string `json:"organizationID"`
	Version        int    `json:"version"`
	CreatedAt      string `json:"createdAt"`
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

// Структура для создания нового предложения
type Bid struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderID    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorID    string `json:"authorId"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}
