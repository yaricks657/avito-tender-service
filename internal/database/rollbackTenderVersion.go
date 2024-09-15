package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

func (strg *Storage) RollbackTenderVersion() (models.Tender, int, error) {
	db := strg.Mng.Db

	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return models.Tender{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Обработчик для отката в случае ошибки
	defer func() {
		if err != nil {
			tx.Rollback() // Откатываем транзакцию
		}
	}()

	// Сохраняем текущую версию в архив
	err = archiveCurrentTenderVersion(tx, strg.TenderId)
	if err != nil {
		return models.Tender{}, http.StatusBadRequest, err
	}

	// Получаем запись целевой версии
	query := `
        SELECT name, description, service_type, status, version, created_at, created_by, organization_id
        FROM tender_versions
        WHERE tender_id = $1 AND version = $2
    `
	row := tx.QueryRow(query, strg.TenderId, strg.Version)
	var tender models.Tender
	err = row.Scan(&tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.Version, &tender.CreatedAt, &tender.CreatedBy, &tender.OrganizationId)
	if err != nil {
		return tender, http.StatusNotFound, err
	}

	// Обновляем текущую версию на целевую
	updateQuery := `
        UPDATE tenders
        SET name = $1, description = $2, service_type = $3, status = $4, version = version + 1, created_at = $5, created_by = $6, organization_id = $7
        WHERE id = $8 AND version = (SELECT MAX(version) FROM tenders WHERE id = $8)
    `
	_, err = tx.Exec(updateQuery, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.CreatedAt, tender.CreatedBy, tender.OrganizationId, strg.TenderId)
	if err != nil {
		return tender, http.StatusInternalServerError, err
	}

	// Если всё прошло успешно, фиксируем изменения
	err = tx.Commit()
	if err != nil {
		return tender, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tender, http.StatusOK, nil
}

// проверка наличия обязательных полей для отката версии
func (strg *Storage) checkRequiredFieldsForRollback() error {
	if strg.Username == "" {
		return fmt.Errorf("no username")
	}
	if strg.TenderId == "" {
		return fmt.Errorf("no tenderId")
	}
	if strg.Version == 0 {
		return fmt.Errorf("no version")
	}
	return nil
}

// Изменение функции archiveCurrentTenderVersion для использования транзакции
func archiveCurrentTenderVersion(tx *sql.Tx, tenderID string) error {
	query := `
        SELECT name, description, service_type, status, version, created_at, created_by, organization_id
        FROM tenders
        WHERE id = $1 AND version = (SELECT MAX(version) FROM tenders WHERE id = $1)
    `
	row := tx.QueryRow(query, tenderID)

	var name, description, serviceType, status, createdBy, organizationID string
	var version int
	var createdAt time.Time

	err := row.Scan(&name, &description, &serviceType, &status, &version, &createdAt, &createdBy, &organizationID)
	if err != nil {
		return err
	}

	archiveQuery := `
        INSERT INTO tender_versions (tender_id, name, description, service_type, status, version, created_by, organization_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err = tx.Exec(archiveQuery, tenderID, name, description, serviceType, status, version, createdBy, organizationID)
	if err != nil {
		return err
	}

	return nil
}
