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

// Статтусы тендеров
const (
	statusPublished = "PUBLISHED"
	statusCreated   = "CREATED"
	statusClosed    = "CLOSED"
)

// отправка на клиент ответа об ошибке
func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
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

// Валидация полей запроса
func checkRequiredFields(t *models.RequestAddTender) error {
	// проверка ServiceType на валидность
	if t.ServiceType == "" || !isServiceTypeAllowed(t.ServiceType) {
		return fmt.Errorf("not allowed serviceType")
	}
	// проверка на

	return nil
}
