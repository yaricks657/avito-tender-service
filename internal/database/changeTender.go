package database

import (
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на смену тендера в БД
func (strg *Storage) ChangeTender(t *models.Tender) (models.Tender, int, error) {
	db := strg.Mng.Db

	err := strg.checkRequiredFieldsForChangeTender()
	if err != nil {
		return models.Tender{}, http.StatusBadRequest, err
	}

	// копируем текущую версию в таблицу истории
	historyQuery := `
	  INSERT INTO tender_versions (tender_id, name, description, service_type, status, version, created_by, organization_id, created_at)
	  SELECT id, name, description, service_type, status, version, created_by, organization_id, created_at
	  FROM tenders
	  WHERE id = $1;
  `
	_, err = db.Exec(historyQuery, strg.TenderId)
	if err != nil {
		return models.Tender{}, http.StatusBadRequest, err
	}

	fmt.Println(&t)
	// Обновляем запись в основной таблице и увеличиваем версию
	updateQuery := `
	  UPDATE tenders
	  SET name = $2, description = $3, service_type = $4, status = $5, version = version + 1
	  WHERE id = $1
	  RETURNING id, name, description, service_type, status, version, created_at;
  `
	var updated models.Tender
	err = db.QueryRow(updateQuery, strg.TenderId, t.Name, t.Description, t.ServiceType, t.Status).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.ServiceType,
		&updated.Status,
		&updated.Version,
		&updated.CreatedAt,
	)

	if err != nil {
		return models.Tender{}, http.StatusBadRequest, err
	}

	return updated, 200, nil

}

// проверка наличия обязательных полей при редактировании тендера
func (strg *Storage) checkRequiredFieldsForChangeTender() error {
	if strg.Username == "" {
		return fmt.Errorf("no username")
	}
	if strg.TenderId == "" {
		return fmt.Errorf("no tenderId")
	}
	return nil
}
