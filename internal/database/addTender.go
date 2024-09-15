package database

import "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/models"

// добавить тендер в БД
func (strg *Storage) AddTender(t *models.Tender) (models.Tender, error) {
	db := strg.Mng.Db

	query := `
	INSERT INTO tenders ( name, description, service_type, status, created_by, organization_id, version, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, name, description, service_type, status, version, created_at
`
	var tender models.Tender
	// Выполняем запрос на добавление и сразу возвращаем добавленные значения
	err := db.QueryRow(
		query,
		t.Name,        // $2
		t.Description, // $3
		t.ServiceType, // $4
		t.Status,      //$5
		t.CreatedBy,   // $6
		t.OrganizationId,
		t.Version,   // $7
		t.CreatedAt, // $8
	).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&tender.CreatedAt,
	)

	if err != nil {
		return models.Tender{}, err
	}

	// Сохранить версию тендера
	versionQuery := `
	   INSERT INTO tender_versions (tender_id, name, description, service_type, status, created_by,organization_id, version, created_at)
	   VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)
   `
	_, err = db.Exec(
		versionQuery,
		tender.ID,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.Status,
		t.CreatedBy,
		t.OrganizationId,
		tender.Version,
		tender.CreatedAt,
	)

	if err != nil {
		return models.Tender{}, err
	}

	return tender, nil
}
