package database

import (
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на смену тендера в БД
func (strg *Storage) ChangeBid(t *models.Bid) (models.Bid, int, error) {
	db := strg.Mng.Db

	err := strg.checkRequiredFieldsForChangeBid()
	if err != nil {
		return models.Bid{}, http.StatusBadRequest, err
	}

	// копируем текущую версию в таблицу истории
	historyQuery := `
	  INSERT INTO bid_versions (bid_id, name, description, tender_id, author_type, author_id, status, version, created_at)
	  SELECT id, name, description, tender_id, author_type, author_id, status, version, created_at
	  FROM bids
	  WHERE id = $1;
  `
	_, err = db.Exec(historyQuery, strg.BidId)
	if err != nil {
		return models.Bid{}, http.StatusBadRequest, err
	}

	fmt.Println(&t)
	// Обновляем запись в основной таблице и увеличиваем версию
	updateQuery := `
	  UPDATE bids
	  SET name = $2, description = $3, author_type = $4, status = $5, version = version + 1
	  WHERE id = $1
	  RETURNING id, name, status, author_type, author_id, version, created_at;
  `
	var updated models.Bid
	err = db.QueryRow(updateQuery, strg.BidId, t.Name, t.Description, t.AuthorType, t.Status).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Status,
		&updated.AuthorType,
		&updated.AuthorID,
		&updated.Version,
		&updated.CreatedAt,
	)

	if err != nil {
		return models.Bid{}, http.StatusBadRequest, err
	}

	return updated, 200, nil

}

// проверка наличия обязательных полей при редактировании тендера
func (strg *Storage) checkRequiredFieldsForChangeBid() error {
	if strg.Username == "" {
		return fmt.Errorf("no username")
	}
	if strg.BidId == "" {
		return fmt.Errorf("no tenderId")
	}
	return nil
}
