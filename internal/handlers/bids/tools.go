package bids

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

var allowedStatuses = [3]string{"Published", "Created", "Closed"}

// Статтусы тендеров
const (
	statusPublished = "Published"
	statusCreated   = "Created"
	statusClosed    = "Closed"
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

// Валидация полей запроса Post
func checkRequiredFields(t *models.RequestAddBid) error {

	// проверка наличия Name
	if t.Name == "" {
		return fmt.Errorf("need name")
	}
	// проверка наличия Description
	if t.Description == "" {
		return fmt.Errorf("need description")
	}

	// проверка наличия tenderId
	if t.TenderID == "" {
		return fmt.Errorf("need tenderId")
	}
	// проверка наличия AuthorId
	if t.AuthorID == "" {
		return fmt.Errorf("need AuthorId")
	}
	return nil
}

// Функция для проверки корректного статуса
func isStatusAllowed(status string) bool {
	for _, t := range allowedStatuses {
		if t == status {
			return true
		}
	}
	return false
}

// получить параметры из запроса для bids
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
	} else {
		strg.Limit = int32(limit)
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
	} else {
		strg.Offset = int32(offset)
	}

	// парсим user_name
	if user := r.URL.Query().Get("username"); user != "" {
		strg.Username = user
	}

	// парсим статус
	if status := r.URL.Query().Get("status"); status != "" {
		strg.Status = status
	}

	// Получаем полный путь запроса
	path := r.URL.Path
	start := strings.Index(path, "/bids/")
	if start == -1 {
		strg.TenderId = "" // Если не найден, записываем пустую строку
	} else {
		actions := []string{"/list"} // все возможные окончания для путей
		var end int
		for _, act := range actions {
			end = strings.Index(path, act)
			if end != -1 {
				break
			}
		}
		if end == -1 {
			strg.TenderId = "" // Если не найдено окончание, записываем пустую строку
		} else {
			tenderID := path[start+len("/bids/") : end]
			strg.TenderId = tenderID
		}
	}

	// Парсим версию тендера
	urlPath := r.URL.Path
	re := regexp.MustCompile(`/rollback/([0-9]+)$`)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) == 2 {
		version := matches[1] // захватываем версию
		if version != "" {
			num, err := strconv.Atoi(version)
			if err != nil {
				fmt.Println("Ошибка:", err)
				return err, fmt.Sprintln("Некорректная версия")
			}
			strg.Version = int32(num)
		}
	}

	// Парсим bidId
	startBid := strings.Index(path, "/bids/")
	if startBid == -1 {
		strg.BidId = "" // Если не найден, записываем пустую строку
	} else {
		actionsBid := []string{"/status", "/edit", "/rollback"}
		var endBid int
		for _, act := range actionsBid {
			endBid = strings.Index(path, act)
			if endBid != -1 {
				break
			}
		}
		if endBid == -1 {
			strg.BidId = "" // Если не найдено окончание, записываем пустую строку
		} else {
			bidID := path[startBid+len("/bids/") : endBid]
			strg.BidId = bidID
		}
	}

	return nil, ""
}
