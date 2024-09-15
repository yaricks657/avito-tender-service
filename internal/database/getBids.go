package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на получение тендеров из бд
func (strg *Storage) GetBids() ([]models.Bid, int, error) {
	db := strg.Mng.Db
	userIdFetch := ""
	// Сначала проверяем, существует ли пользователь, если указан
	if strg.Username != "" {
		// Запрос для получения id пользователя
		userQuery := `SELECT id 
					  FROM employee 
					  WHERE username = $1`
		var userId string
		err := db.QueryRow(userQuery, strg.Username).Scan(&userId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, http.StatusUnauthorized, fmt.Errorf("no such user") // Пользователь не существует
			}
			strg.Mng.Log.LogError("Error retrieving user id", err)
			return nil, http.StatusInternalServerError, err
		}

		// Присваиваем id пользователя в структуру strg или используем дальше
		userIdFetch = userId
	}
	// Базовый запрос
	query := `SELECT id, name, description, author_type,author_id, status, version, created_at 
		 FROM bids`
	queryParams := []interface{}{}

	// Добавляем фильтрацию по пользователю, если указан
	if strg.Username != "" {
		if len(queryParams) > 0 {
			query += ` AND author_id = $` + strconv.Itoa(len(queryParams)+1)
		} else {
			query += ` WHERE author_id = $` + strconv.Itoa(len(queryParams)+1)
		}
		queryParams = append(queryParams, userIdFetch)
	}

	// Добавляем фильтрацию по тендеру, если указан
	if strg.TenderId != "" {
		if len(queryParams) > 0 {
			query += ` AND tender_id = $` + strconv.Itoa(len(queryParams)+1)
		} else {
			query += ` WHERE tender_id = $` + strconv.Itoa(len(queryParams)+1)
		}
		queryParams = append(queryParams, strg.TenderId)
	}

	// Добавляем пагинацию
	query += ` ORDER BY name LIMIT $` + strconv.Itoa(len(queryParams)+1) + ` OFFSET $` + strconv.Itoa(len(queryParams)+2)
	queryParams = append(queryParams, strg.Limit, strg.Offset)

	// запрос в БД
	rows, err := db.Query(query, queryParams...)
	if err != nil {
		strg.Mng.Log.LogError("Error fetching tenders", err)
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	// Collect tenders
	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.AuthorType, &bid.AuthorID, &bid.Status, &bid.Version, &bid.CreatedAt); err != nil {
			strg.Mng.Log.LogError("Error scanning tenders", err)
			return nil, http.StatusInternalServerError, err
		}
		bids = append(bids, bid)
	}

	return bids, 0, nil
}
