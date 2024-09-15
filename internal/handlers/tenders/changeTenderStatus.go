package tenders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
)

// Обработчик для добавления задачи в БД (/api/tenders/new)
func ChangeTenderStatusHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// проверка на нужный метод
	if r.Method != http.MethodPut {
		manager.Mng.Log.LogWarn("Некорректный метод запроса")
		handlers.SendErrorResponse(w, "Некорректный метод запроса", http.StatusMethodNotAllowed)
		return
	}

	// получить параметры из запроса
	err, errMessage := getParams(r, &dbStruct)
	if errMessage != "" {
		handlers.SendErrorResponse(w, errMessage, 400)
		manager.Mng.Log.LogError(errMessage, err)
		return
	}

	// проверка на корректный статус для смены
	ok := handlers.IsStatusAllowed(dbStruct.Status)
	if !ok {
		handlers.SendErrorResponse(w, "Некорректный статус", 400)
		manager.Mng.Log.LogError("Некорректный статус", err)
		return
	}

	// проверка доступа на редактирование
	tender, statusCode, err := dbStruct.GetTenderWithoutParams()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}
	err, statusCode = dbStruct.IsUserInOrganization(dbStruct.Username, tender.OrganizationId)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при проверке прав доступа", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка проверки прав доступа: %s", err), statusCode)
		return
	}

	// отправка запроса на изменение статуса в БД
	tenderResponse, statusCode, err := dbStruct.ChangeTenderStatus()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при записи тендера в БД", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка при записи тендера в БД: %s", err), statusCode)
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
