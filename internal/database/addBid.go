package database

import (
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"
)

// добавить тендер в БД
func (strg *Storage) AddBid(b *models.Bid) (models.Bid, error) {
	db := strg.Mng.Db

	// Запрос для добавления нового предложения
	query := `
    INSERT INTO bids (name, description, tender_id, author_type, author_id, status, created_at, version)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id, name, author_type, author_id, status, created_at,version
    `

	var bid models.Bid
	// Выполняем запрос на добавление и сразу возвращаем добавленные значения
	err := db.QueryRow(
		query,
		b.Name,        // $1
		b.Description, // $2
		b.TenderID,    // $3
		b.AuthorType,  // $4
		b.AuthorID,    // $5
		b.Status,      // $6
		b.CreatedAt,   // $7
		b.Version,
	).Scan(
		&bid.ID,
		&bid.Name,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Status,
		&bid.CreatedAt,
		&bid.Version,
	)

	if err != nil {
		return models.Bid{}, err
	}

	// Сохранить версию предложения
	versionQuery := `
    INSERT INTO bid_versions (bid_id, name, description, tender_id, author_type, author_id, status, created_at,version)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)
    `
	fmt.Println(bid)
	_, err = db.Exec(
		versionQuery,
		bid.ID,
		bid.Name,
		b.Description,
		b.TenderID,
		bid.AuthorType,
		bid.AuthorID,
		bid.Status,
		bid.CreatedAt,
		bid.Version,
	)

	if err != nil {
		return models.Bid{}, err
	}

	return bid, nil
}
