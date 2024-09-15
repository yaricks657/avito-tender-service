package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

func (strg *Storage) RollbackBidVersion() (models.Bid, int, error) {
	db := strg.Mng.Db

	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return models.Bid{}, http.StatusInternalServerError, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Обработчик для отката в случае ошибки
	defer func() {
		if err != nil {
			tx.Rollback() // Откатываем транзакцию
		}
	}()

	// Сохраняем текущую версию в архив
	err = archiveCurrentBidVersion(tx, strg.BidId)
	if err != nil {
		return models.Bid{}, http.StatusBadRequest, err
	}

	// Получаем запись целевой версии
	query := `
        SELECT name, description, tender_id, author_type, author_id, status, version, created_at
        FROM bid_versions
        WHERE bid_id = $1 AND version = $2
    `
	row := tx.QueryRow(query, strg.BidId, strg.Version)
	var bid models.Bid
	err = row.Scan(&bid.Name, &bid.Description, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Status, &bid.Version, &bid.CreatedAt)
	if err != nil {
		return bid, http.StatusNotFound, err
	}

	// Обновляем текущую версию на целевую
	updateQuery := `
        UPDATE bids
        SET name = $1, description = $2, tender_id = $3, author_type = $4, version = version + 1, author_id = $5, status = $6, created_at = $7
        WHERE id = $8 AND version = (SELECT MAX(version) FROM bids WHERE id = $8)
    `
	_, err = tx.Exec(updateQuery, bid.Name, bid.Description, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.Status, bid.CreatedAt, strg.BidId)
	if err != nil {
		return bid, http.StatusInternalServerError, err
	}

	// Если всё прошло успешно, фиксируем изменения
	err = tx.Commit()
	if err != nil {
		return bid, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return bid, http.StatusOK, nil
}

// проверка наличия обязательных полей для отката версии
func (strg *Storage) checkRequiredFieldsForRollbackBid() error {
	if strg.Username == "" {
		return fmt.Errorf("no username")
	}
	if strg.BidId == "" {
		return fmt.Errorf("no tenderId")
	}
	if strg.Version == 0 {
		return fmt.Errorf("no version")
	}
	return nil
}

// Изменение функции archiveCurrentTenderVersion для использования транзакции
func archiveCurrentBidVersion(tx *sql.Tx, bidId string) error {
	query := `
    SELECT name, description, tender_id, author_type, author_id, status, version, created_at
    FROM bids
    WHERE id = $1 AND version = (SELECT MAX(version) FROM bids WHERE id = $1)
`
	row := tx.QueryRow(query, bidId)

	var name, description, tenderId, authorType, authorId, status string
	var version int
	var createdAt time.Time

	// Используйте правильный порядок и количество переменных при Scan
	err := row.Scan(&name, &description, &tenderId, &authorType, &authorId, &status, &version, &createdAt)
	if err != nil {
		return err
	}

	// Запрос для вставки данных в таблицу версий
	archiveQuery := `
    INSERT INTO bid_versions (bid_id, name, description, tender_id, author_type, author_id, status, version, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	_, err = tx.Exec(archiveQuery, bidId, name, description, tenderId, authorType, authorId, status, version, createdAt)
	if err != nil {
		return err
	}

	return nil
}
