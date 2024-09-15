package database

import (
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на смену статуса тендера в БД
func (strg *Storage) ChangeTenderStatus() (models.Tender, int, error) {
	db := strg.Mng.Db

	err := strg.checkRequiredFieldsForChangeStatus()
	if err != nil {
		return models.Tender{}, http.StatusBadRequest, err
	}

	query := `
	UPDATE tenders
	SET status = $1, 
		version = version
	WHERE id = $2
	RETURNING id, name, description, service_type, status, version, created_at;
`

	var tender models.Tender
	err = db.QueryRow(query, strg.Status, strg.TenderId).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&tender.CreatedAt,
	)

	if err != nil {
		return models.Tender{}, http.StatusBadRequest, err
	}

	return tender, 200, nil
}

// проверка наличия обязательных полей
func (strg *Storage) checkRequiredFieldsForChangeStatus() error {
	if strg.Username == "" {
		return fmt.Errorf("no username")
	}
	if strg.TenderId == "" {
		return fmt.Errorf("no tenderId")
	}
	if strg.Status == "" {
		return fmt.Errorf("no status")
	}
	return nil
}
