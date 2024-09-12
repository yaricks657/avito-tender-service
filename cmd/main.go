package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/config"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/manager"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	// создание логгера
	logFilePath := "./app.log"
	logger, err := logger.New(logFilePath)
	if err != nil {
		log.Fatal("Ошибка при создании логгера (main)", err)
		os.Exit(1)
	}

	// загрузка конфига
	config, err := config.GetEnv()
	if err != nil {
		logger.LogError("Ошибка при загрузке конфига (main)", err)
		os.Exit(1)
	}
	logger.LogInfo("Переменные окружения получены", fmt.Sprintln(config))

	// создание контейнера общего
	createConfig := manager.CreateConfig{
		Cnf: config,
		Log: logger,
	}
	// создаем менеджер для дальнейшей работы с ним
	Mng, err := manager.New(createConfig)
	if err != nil {
		logger.LogError("Ошибка при регистрации manager (main) ", err)
		os.Exit(1)
	}
	logger.LogInfo("manager зарегистрирован", fmt.Sprintln(&Mng))

	// подключение БД
	_, err = database.CreateDB(Mng)
	if err != nil {
		logger.LogError("Ошибка при подключении к БД (main) ", err)
		os.Exit(1)
	}
	logger.LogInfo("БД подключена успешно")

	// создание и запуск сервера
	router := chi.NewRouter()

	/* Установка обработчиков */
	// проверить работу сервера
	router.Get("/api/ping", handlers.CheckPing)
	// получить тендеры
	router.Get("/api/tenders", handlers.GetTendersHandler)
	// создать тендер
	router.Post("/api/tenders/new", handlers.AddTendersHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerAddress), router); err != nil {
		logger.LogError("Ошибка запуска сервера (main)", err)
		os.Exit(1)
	}
	//logger.LogInfo("Сервер запущен на порту", config.ServerPort)

	// закрываем БД
	defer Mng.Db.Close()

}
