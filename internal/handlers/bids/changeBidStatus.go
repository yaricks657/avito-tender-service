package bids

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
)

// Обработчик для добавления задачи в БД (/api/tenders/new)
func ChangeBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// проверка на нужный метод
	if r.Method != http.MethodPut {
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

	// проверка на корректный статус для смены
	ok := isStatusAllowed(dbStruct.Status)
	if !ok {
		sendErrorResponse(w, "Некорректный статус", 400)
		manager.Mng.Log.LogError("Некорректный статус", err)
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
	fmt.Println(dbStruct)
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

	// отправка запроса на изменение статуса в БД
	bidResponse, statusCode, err := dbStruct.ChangeBidStatus()
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
