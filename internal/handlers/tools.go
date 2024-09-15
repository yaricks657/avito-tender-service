package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// массив разрешенных типов сервисов
var allowedServiceTypes = [3]string{"Construction", "Delivery", "Manufacture"}
var allowedStatuses = [3]string{"Published", "Created", "Closed"}

// Статтусы тендеров
const (
	StatusPublished = "Published"
	StatusCreated   = "Created"
	StatusClosed    = "Closed"
)

// отправка на клиент ответа об ошибке
func SendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := models.ErrorResponse{
		Reason: errorMsg,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании (sendErrorResponse)", err)
	}
}

// Валидация полей запроса Post
func CheckRequiredFields(t *models.RequestAddTender) error {
	// проверка ServiceType на валидность
	if t.ServiceType == "" || !IsServiceTypeAllowed(t.ServiceType) {
		return fmt.Errorf("not allowed serviceType")
	}
	// проверка наличия Name
	if t.Name == "" {
		return fmt.Errorf("need name")
	}
	// проверка наличия Description
	if t.Description == "" {
		return fmt.Errorf("need description")
	}
	return nil
}

// Функция для проверки, есть ли строка в массиве
func IsServiceTypeAllowed(serviceType string) bool {
	for _, t := range allowedServiceTypes {
		if t == serviceType {
			return true
		}
	}
	return false
}

// Функция для проверки корректного статуса
func IsStatusAllowed(status string) bool {
	for _, t := range allowedStatuses {
		if t == status {
			return true
		}
	}
	return false
}
