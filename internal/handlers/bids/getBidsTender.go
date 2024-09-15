package bids

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
)

// Обработчик для получить тендеры по фильтру или без
func GetBidsTenderHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// получить параметры из запроса
	err, errMessage := getParams(r, &dbStruct)
	if errMessage != "" {
		sendErrorResponse(w, errMessage, 400)
		manager.Mng.Log.LogError(errMessage, err)
		return
	}

	// запрос в бд
	tender, statusCode, err := dbStruct.GetTenderWithoutParams()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}
	// проверка наличия доступа на редактирование компании
	err, statusCode = dbStruct.IsUserInOrganizationOrCreator("", dbStruct.Username, tender.OrganizationId, tender.CreatedBy)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при проверке прав доступа", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка проверки прав доступа: %s", err), statusCode)
		return
	}

	dbStruct.Username = ""
	// запрос в БД на получение связанных предложений
	bids, statusCode, err := dbStruct.GetBids()
	if err != nil {
		manager.Mng.Log.LogError("ошибка при получении данных из БД", err)
		sendErrorResponse(w, fmt.Sprintf("ошибка при получении данных из БД %s", err), statusCode)
		return
	}

	jsonResponse, err := json.Marshal(bids)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при маршалинге данных: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при маршалинге данных: %s", err), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при отправке ответа: ", err)
		return
	}

}
