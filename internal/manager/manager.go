package manager

import (
	"database/sql"
	"fmt"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/config"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// методы для manager
type IManager interface {
	SetHandlers(r *chi.Mux)
}

// структура-контейнер для управления проектом
type Manager struct {
	//стартовый конфиг
	Cnf config.Environment

	// логирование zerolog
	Log *logger.Logger

	// бд
	Db *sql.DB
}

// структура для создания manager
type CreateConfig struct {
	Cnf config.Environment
	Log *logger.Logger
}

// Глобальная переменная с начальной конфигурацией
var Mng Manager

// старт manager
func New(config CreateConfig) (*Manager, error) {
	// добавить проверку требуемых полей методом для CreateConfig
	err := config.checkRequiredFields()
	if err != nil {
		return nil, err
	}

	Mng = Manager{
		Log: config.Log,
		Cnf: config.Cnf,
	}

	return &Mng, nil
}

// проверка наличия обязательных полей
func (cc *CreateConfig) checkRequiredFields() error {
	var missingFields []string

	if cc.Cnf.ServerAddress == "" {
		missingFields = append(missingFields, "SERVER_ADDRESS")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("отсутствуют обязательные поля: %s", strings.Join(missingFields, ", "))
	}
	return nil
}
