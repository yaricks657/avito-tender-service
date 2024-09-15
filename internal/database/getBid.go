package database

import (
	"database/sql"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на получение тендеров из бд
func (strg *Storage) GetBid() (models.Bid, int, error) {
	db := strg.Mng.Db
	// Сначала проверяем, существует ли пользователь, если указан
	if strg.Username != "" {
		// Запрос для проверки существования пользователя
		userExistsQuery := `SELECT EXISTS (
					SELECT 1 
					FROM employee 
					WHERE username = $1
				)`
		var userExists bool
		err := db.QueryRow(userExistsQuery, strg.Username).Scan(&userExists)
		if err != nil {
			strg.Mng.Log.LogError("Error checking user existence", err)
			return models.Bid{}, http.StatusInternalServerError, err
		}
		if !userExists {
			return models.Bid{}, http.StatusUnauthorized, fmt.Errorf("no such user") // Пользователь не существует
		}
	}

	// Базовый запрос
	query := `SELECT id, name, description, author_type,author_id, status, version, created_at, tender_id
	FROM bids`
	queryParams := []interface{}{}

	// Добавляем фильтрацию по ID, если указан
	if strg.BidId != "" {
		query += ` WHERE id = $1`
		queryParams = append(queryParams, strg.BidId)
	} else {
		return models.Bid{}, http.StatusNotFound, fmt.Errorf("no id") // Запись не найдена, возвращаем пустую структуру
	}

	// Выполняем запрос
	row := db.QueryRow(query, queryParams...)

	// Обрабатываем результат
	var bid models.Bid
	err := row.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.AuthorType, &bid.AuthorID, &bid.Status, &bid.Version, &bid.CreatedAt, &bid.TenderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Bid{}, http.StatusNotFound, err // Запись не найдена, возвращаем пустую структуру
		}
		return models.Bid{}, http.StatusBadRequest, err
	}

	return bid, 0, nil
}
