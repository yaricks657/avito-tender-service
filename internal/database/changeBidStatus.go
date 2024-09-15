package database

import (
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// запрос на смену статуса тендера в БД
func (strg *Storage) ChangeBidStatus() (models.Bid, int, error) {
	db := strg.Mng.Db

	err := strg.checkRequiredFieldsForChangeStatusBid()
	if err != nil {
		return models.Bid{}, http.StatusBadRequest, err
	}

	query := `
	UPDATE bids
	SET status = $1, 
		version = version
	WHERE id = $2
	RETURNING id, name, status, author_type, author_id, version, created_at;
`

	var bid models.Bid
	err = db.QueryRow(query, strg.Status, strg.BidId).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		return models.Bid{}, http.StatusBadRequest, err
	}

	return bid, 200, nil
}

// проверка наличия обязательных полей
func (strg *Storage) checkRequiredFieldsForChangeStatusBid() error {
	if strg.Username == "" {
		return fmt.Errorf("no username")
	}
	if strg.BidId == "" {
		return fmt.Errorf("no bidId")
	}
	if strg.Status == "" {
		return fmt.Errorf("no status")
	}
	return nil
}
