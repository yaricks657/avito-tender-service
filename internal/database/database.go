package database

import (
	"database/sql"
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	_ "github.com/lib/pq"
)

// структура для БД и филт=ьтра для запросов
type Storage struct {
	Mng *manager.Manager

	Limit        int32
	Offset       int32
	Service_type []string
	Username     string
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

	return db, nil
}
