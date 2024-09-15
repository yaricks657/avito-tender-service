package tenders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
)

// Обработчик для получить тендеры по фильтру или без
func GeMyTendersHandler(w http.ResponseWriter, r *http.Request) {
	dbStruct := database.Storage{
		Mng: &manager.Mng,
	}

	// получить параметры из запроса
	err, errMessage := getParams(r, &dbStruct)
	if errMessage != "" {
		handlers.SendErrorResponse(w, errMessage, 400)
		manager.Mng.Log.LogError(errMessage, err)
		return
	}

	if dbStruct.Username == "" {
		manager.Mng.Log.LogError("пользователь не указан: ", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("пользователь не указан: %s", err), 401)
		return
	}

	// запрос в бд
	tenders, statusCode, err := dbStruct.GetTenders()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}

	jsonResponse, err := json.Marshal(tenders)
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при маршалинге данных: ", err)
		handlers.SendErrorResponse(w, fmt.Sprintf("Ошибка при маршалинге данных: %s", err), 500)
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
