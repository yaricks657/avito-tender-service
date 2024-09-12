package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
)

// Обработчик для получить тендеры по фильтру или без
func GetTendersHandler(w http.ResponseWriter, r *http.Request) {
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
	tenders, statusCode, err := dbStruct.GetTenders()
	if err != nil {
		manager.Mng.Log.LogError("Ошибка при обращении к БД: ", err)
		sendErrorResponse(w, fmt.Sprintf("Ошибка при обращении к БД: %s", err), statusCode)
		return
	}

	jsonResponse, err := json.Marshal(tenders)
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

// получить параметры из запроса
func getParams(r *http.Request, strg *database.Storage) (error, string) {
	// Стандартные значения пагинации. Нужно вынести в .env
	limit := 50
	offset := 0

	// парсим limit
	if l := r.URL.Query().Get("limit"); l != "" {
		parsedLimit, err := strconv.Atoi(l)
		if err != nil {
			return err, fmt.Sprintln("Некорректное значение параметра limit")
		}
		if parsedLimit > limit {
			return nil, fmt.Sprintln("Параметр limit превышает максимальное число")
		}
		strg.Limit = int32(parsedLimit)
	}

	// парсим offset
	if o := r.URL.Query().Get("offset"); o != "" {
		parsedOffset, err := strconv.Atoi(o)
		if err != nil {
			return err, fmt.Sprintln("Некорректное значение параметра offset")
		}
		if parsedOffset < offset {
			return nil, fmt.Sprintln("Параметр limit отрицательный")
		}
		strg.Offset = int32(parsedOffset)
	}

	// парсим service_type
	serviceTypes := r.URL.Query()["service_type"]
	if len(serviceTypes) > 0 {
		for _, service := range serviceTypes {
			if !isServiceTypeAllowed(service) {
				return nil, fmt.Sprintln("Некорректный параметр service_type")
			}
		}
		strg.Service_type = serviceTypes
	}

	// парсим user_name
	if user := r.URL.Query().Get("user_name"); user != "" {
		strg.Username = user
	}

	return nil, ""
}

// Функция для проверки, есть ли строка в массиве
func isServiceTypeAllowed(serviceType string) bool {
	for _, t := range allowedServiceTypes {
		if t == serviceType {
			return true
		}
	}
	return false
}
