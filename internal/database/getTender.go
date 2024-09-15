package database

import (
	"database/sql"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на получение тендеров из бд
func (strg *Storage) GetTender() (models.Tender, int, error) {
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
			return models.Tender{}, http.StatusInternalServerError, err
		}
		if !userExists {
			return models.Tender{}, http.StatusUnauthorized, fmt.Errorf("no such user") // Пользователь не существует
		}
	}

	// Базовый запрос
	query := `SELECT id, name, description, service_type, status, version,organization_id, created_at 
          FROM tenders`
	queryParams := []interface{}{}

	// Добавляем фильтрацию по ID, если указан
	if strg.TenderId != "" {
		query += ` WHERE id = $1`
		queryParams = append(queryParams, strg.TenderId)
	} else {
		return models.Tender{}, http.StatusNotFound, fmt.Errorf("no id") // Запись не найдена, возвращаем пустую структуру
	}

	// Добавляем фильтрацию по пользователю, если указан
	if strg.Username != "" {
		if len(queryParams) > 0 {
			query += ` AND created_by = $2`
			queryParams = append(queryParams, strg.Username)
		} else {
			query += ` WHERE created_by = $1`
			queryParams = append(queryParams, strg.Username)
		}
	}

	// Выполняем запрос
	row := db.QueryRow(query, queryParams...)

	// Обрабатываем результат
	var t models.Tender
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status, &t.Version, &t.OrganizationId, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tender{}, http.StatusNotFound, err // Запись не найдена, возвращаем пустую структуру
		}
		return models.Tender{}, http.StatusBadRequest, err
	}

	return t, 0, nil
}
