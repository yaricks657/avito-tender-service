package database

import (
	"database/sql"
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	_ "github.com/lib/pq"
)

// структура для БД и фильтры для запросов
type Storage struct {
	Mng *manager.Manager

	Limit        int32
	Offset       int32
	Service_type []string
	Username     string
	TenderId     string
	BidId        string
	Status       string
	Version      int32
}

// Создание подключения к БД
func CreateDB(mng *manager.Manager) (*sql.DB, error) {
	// Параметры подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require target_session_attrs=read-write",
		mng.Cnf.PostgresHost, mng.Cnf.PostgresPort, mng.Cnf.PostgresUsername, mng.Cnf.PostgresPassword, mng.Cnf.PostgresDatabase)

	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	mng.Db = db

	// SQL запрос для создания таблицы tenders, если она еще не существует
	tendersTable := `
		CREATE TABLE IF NOT EXISTS tenders (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			service_type VARCHAR(100),
			status VARCHAR(50) NOT NULL,
			created_by VARCHAR(100) NOT NULL,
			organization_id UUID NOT NULL,
			version INT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	// SQL запрос для создания таблицы tender_versions, если она еще не существует
	tenderVersionsTable := `
		CREATE TABLE IF NOT EXISTS tender_versions (
			id UUID DEFAULT uuid_generate_v4(),
			tender_id UUID NOT NULL REFERENCES tenders(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			service_type VARCHAR(100),
			status VARCHAR(50) NOT NULL,
			created_by VARCHAR(100) NOT NULL,
			organization_id UUID NOT NULL,
			version INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (tender_id, version)
		);`

	// Выполняем запросы создания таблиц
	_, err = db.Exec(tendersTable)
	if err != nil {
		return db, fmt.Errorf("error creating tenders table: %v", err)
	}

	_, err = db.Exec(tenderVersionsTable)
	if err != nil {
		return db, fmt.Errorf("error creating tender_versions table: %v", err)
	}

	// SQL для создания таблицы bids
	bidsTable := `
		CREATE TABLE IF NOT EXISTS bids (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			tender_id TEXT,
			author_type VARCHAR(50) NOT NULL,
			author_id TEXT,
			status VARCHAR(50) NOT NULL,
			version INT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	// SQL для создания таблицы bid_versions
	bidVersionsTable := `
		CREATE TABLE IF NOT EXISTS bid_versions (
			id UUID DEFAULT uuid_generate_v4(),
			bid_id UUID NOT NULL REFERENCES bids(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			tender_id TEXT,
			author_type VARCHAR(50) NOT NULL,
			author_id TEXT,
			status VARCHAR(50) NOT NULL,
			version INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (bid_id, version)
		);`

	// Выполняем SQL запросы
	_, err = db.Exec(bidsTable)
	if err != nil {
		return db, fmt.Errorf("error creating bids table:: %v", err)
	}

	_, err = db.Exec(bidVersionsTable)
	if err != nil {
		return db, fmt.Errorf("error creating bid_versions table: %v", err)
	}

	return db, nil
}
