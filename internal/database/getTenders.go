package database

import (
	"net/http"
	"strconv"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на получение тендеров из бд
func (strg *Storage) GetTenders() ([]models.Tender, int, error) {
	db := strg.Mng.Db

	// Базовый запрос
	query := `SELECT id, name, description, service_type, status, version, created_at 
		 FROM tenders`
	queryParams := []interface{}{}

	// Добавляем фильтрацию, если есть
	if len(strg.Service_type) > 0 {
		placeholders := make([]string, len(strg.Service_type))
		for i := range strg.Service_type {
			placeholders[i] = "$" + strconv.Itoa(i+1)
			queryParams = append(queryParams, strg.Service_type[i])
		}
		query += ` WHERE service_type IN (` + strings.Join(placeholders, ", ") + `)`
	}

	// Добавляем пагинацию
	//query += ` ORDER BY name LIMIT $` + strconv.Itoa(len(queryParams)+1) + ` OFFSET $` + strconv.Itoa(len(queryParams)+2)
	//queryParams = append(queryParams, strg.Limit, strg.Offset)

	// запрос в БД
	rows, err := db.Query(query, queryParams...)
	if err != nil {
		strg.Mng.Log.LogError("Error fetching tenders", err)
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	// Collect tenders
	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.Version, &tender.CreatedAt); err != nil {
			strg.Mng.Log.LogError("Error scanning tenders", err)
			return nil, http.StatusInternalServerError, err
		}
		tenders = append(tenders, tender)
	}

	return tenders, 0, nil
}
