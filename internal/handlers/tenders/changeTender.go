package tenders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// Обработчик для добавления задачи в БД (/api/tenders/{tenderId}/edit)
func ChangeTenderHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// проверка на нужный метод
	if r.Method != http.MethodPatch {
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

	// проверка доступа на редактирование
	tenderDB, statusCode, err := dbStruct.GetTenderWithoutParams()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}
	err, statusCode = dbStruct.IsUserInOrganization(dbStruct.Username, tenderDB.OrganizationId)
	if err != nil {
		manager.Mng.Log.LogError("ошибка при проверке прав доступа", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка проверки прав доступа: %s", err), statusCode)
		return
	}

	// структура для БД
	tenderDBNew := models.Tender{
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tenderDB.Status,
		OrganizationId: tenderDB.OrganizationId,
		Version:        tenderDB.Version,
		CreatedAt:      tenderDB.CreatedAt,
		CreatedBy:      tenderDB.CreatedBy,
	}

	// отправка запроса в БД на добавление тендера
	tenderResponse, statusCode, err := dbStruct.ChangeTender(&tenderDBNew)
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
