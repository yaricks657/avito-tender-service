package bids

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// Обработчик для добавления задачи в БД (/api/tenders/{tenderId}/edit)
func ChangeBidHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// проверка на нужный метод
	if r.Method != http.MethodPatch {
		manager.Mng.Log.LogWarn("Некорректный метод запроса")
		sendErrorResponse(w, "Некорректный метод запроса", http.StatusMethodNotAllowed)
		return
	}

	// получить параметры из запроса
	err, errMessage := getParams(r, &dbStruct)
	if errMessage != "" {
		sendErrorResponse(w, errMessage, 400)
		manager.Mng.Log.LogError(errMessage, err)
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
	var tenderCl models.RequestAddBid
	if err = json.Unmarshal(body, &tenderCl); err != nil {
		manager.Mng.Log.LogError("Bad request", err)
		sendErrorResponse(w, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}

	// Проверка авторизации
	bid, statusCode, err := dbStruct.GetBid()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}
	dbStruct.TenderId = bid.TenderID
	tender, statusCode, err := dbStruct.GetTenderWithoutParams()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}
	err, statusCode = dbStruct.IsUserInOrganizationOrCreator("", dbStruct.Username, tender.OrganizationId, tender.CreatedBy)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при проверке прав доступа", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка проверки прав доступа: %s", err), statusCode)
		return
	}

	// структура для БД
	tenderDBNew := models.Bid{
		Name:        tenderCl.Name,
		Description: tenderCl.Description,
		TenderID:    bid.TenderID,
		Status:      bid.Status,
		AuthorType:  bid.AuthorType,
		AuthorID:    bid.AuthorID,
		Version:     bid.Version,
		CreatedAt:   bid.CreatedAt,
	}
	// отправка запроса в БД на добавление тендера
	bidResponse, statusCode, err := dbStruct.ChangeBid(&tenderDBNew)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при записи тендера в БД", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при записи тендера в БД: %s", err), statusCode)
		return
	}

	// отправка успешного ответа клиенту
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bidResponse)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при декодировании успешного ответа", err)
	}

}
