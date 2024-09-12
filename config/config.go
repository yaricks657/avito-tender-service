package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// структура для загрузки переменных окружения
type Environment struct {
	ServerAddress    string `env:"SERVER_ADDRESS" env-default:"8080"`
	PostgresConn     string `env:"POSTGRES_CONN" env-default:"postgres://cnrprod1725726029-team-79287:cnrprod1725726029-team-79287@rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net:6432/cnrprod1725726029-team-79287"`
	PostgresJdbcUrl  string `env:"POSTGRES_JDBC_URL" env-default:"jdbc:postgresql://rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net:6432/cnrprod1725726029-team-79287"`
	PostgresUsername string `env:"POSTGRES_USERNAME" env-default:"cnrprod1725726029-team-79287"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-default:"cnrprod1725726029-team-79287"`
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"6432"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" env-default:"cnrprod1725726029-team-79287"`
}

// получить переменные окружения
func GetEnv() (Environment, error) {
	var env Environment

	if err := cleanenv.ReadEnv(&env); err != nil {
		return env, err
	}

	return env, nil
}
