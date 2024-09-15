package tenders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// Обработчик для добавления задачи в БД (/api/tenders/new)
func AddTendersHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// проверка на нужный метод
	if r.Method != http.MethodPost {
		manager.Mng.Log.LogWarn("Некорректный метод запроса")
		handlers.SendErrorResponse(w, "Некорректный метод запроса", http.StatusMethodNotAllowed)
		return
	}

	// чтение тела запроса в слайс байт
	body, err := io.ReadAll(r.Body)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при чтении тела запроса", err)
		handlers.SendErrorResponse(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Распаковка ответа от клиента
	var tender models.RequestAddTender
	if err = json.Unmarshal(body, &tender); err != nil {
		manager.Mng.Log.LogError("Bad request", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}

	// проверка на наличие обязательных полей и их корректность
	err = handlers.CheckRequiredFields(&tender)
	if err != nil {
		manager.Mng.Log.LogError("Отсутствуют обязательные поля или заполнены неверно", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Отсутствуют обязательные поля или заполнены неверно: %s", err), http.StatusBadRequest)
		return
	}

	// проверка наличия доступа на редактирование компании
	err, statusCode := dbStruct.IsUserInOrganization(tender.CreatorUsername, tender.OrganizationId)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при проверке прав доступа", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка проверки прав доступа: %s", err), statusCode)
		return
	}

	// структура для БД
	tenderDB := models.Tender{
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         handlers.StatusCreated,
		OrganizationId: tender.OrganizationId,
		Version:        1,
		CreatedAt:      time.Now().Format(time.RFC3339),
		CreatedBy:      tender.CreatorUsername,
	}

	// отправка запроса в БД на добавление тендера
	tenderResponse, err := dbStruct.AddTender(&tenderDB)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при записи тендера в БД", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка при записи тендера в БД: %s", err), http.StatusBadRequest)
		return
	}

	// отправка успешного ответа клиенту
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tenderResponse)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании успешного ответа", err)
	}
}
