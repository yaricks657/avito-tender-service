package bids

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// Обработчик для добавления задачи в БД (/api/tenders/new)
func AddBidHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// проверка на нужный метод
	if r.Method != http.MethodPost {
		manager.Mng.Log.LogWarn("Некорректный метод запроса")
		sendErrorResponse(w, "Некорректный метод запроса", http.StatusMethodNotAllowed)
		return
	}

	// чтение тела запроса в слайс байт
	body, err := io.ReadAll(r.Body)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при чтении тела запроса", err)
		sendErrorResponse(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Распаковка ответа от клиента
	var bid models.RequestAddBid
	if err = json.Unmarshal(body, &bid); err != nil {
		manager.Mng.Log.LogError("Bad request", err)
		sendErrorResponse(w, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}

	// проверка на наличие обязательных полей и их корректность
	err = checkRequiredFields(&bid)
	if err != nil {
		manager.Mng.Log.LogError("Отсутствуют обязательные поля или заполнены неверно", err)
		sendErrorResponse(w, fmt.Sprintf("Отсутствуют обязательные поля или заполнены неверно: %s", err), http.StatusBadRequest)
		return
	}

	// проверка прав доступа и наличия тендера
	dbStruct.TenderId = bid.TenderID
	tender, statusCode, err := dbStruct.GetTenderWithoutParams()
	if err != nil {
		manager.Mng.Log.LogError("ошибка при поиске bid", err)
		sendErrorResponse(w, fmt.Sprintf("ошибка при поиске bid: %s", err), statusCode)
		return
	}
	err, statusCode = dbStruct.IsUserInOrganizationOrCreator(bid.AuthorID, "", tender.OrganizationId, tender.CreatedBy)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при проверке прав доступа", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка проверки прав доступа: %s", err), statusCode)
		return
	}

	// структура для БД
	bidrDB := models.Bid{
		Name:        bid.Name,
		Description: bid.Description,
		TenderID:    bid.TenderID,
		Status:      statusCreated,
		AuthorType:  bid.AuthorType,
		AuthorID:    bid.AuthorID,
		Version:     1,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	// отправка запроса в БД на добавление тендера
	tenderResponse, err := dbStruct.AddBid(&bidrDB)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при записи тендера в БД", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при записи тендера в БД: %s", err), http.StatusBadRequest)
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
