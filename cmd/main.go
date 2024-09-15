package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/config"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/database"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers/bids"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726029-team-79287/zadanie-6105/internal/handlers/tenders"
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
	router.Get("/api/ping", tenders.CheckPing)
	// получить тендеры
	router.Get("/api/tenders", tenders.GetTendersHandler)
	// создать тендер
	router.Post("/api/tenders/new", tenders.AddTendersHandler)
	// получить мои тендеры
	router.Get("/api/tenders/my", tenders.GeMyTendersHandler)
	// получить статус тендера
	router.Get("/api/tenders/{tenderId}/status", tenders.GetTenderStatusHandler)
	// смена статуса тендера
	router.Put("/api/tenders/{tenderId}/status", tenders.ChangeTenderStatusHandler)
	// редактирование тендера
	router.Patch("/api/tenders/{tenderId}/edit", tenders.ChangeTenderHandler)
	// откат версии тендера
	router.Put("/api/tenders/{tenderId}/rollback/{version}", tenders.RollbackTenderVersionHandler)
	// добавление новой задачи
	router.Post("/api/bids/new", bids.AddBidHandler)
	// мои предложения
	router.Get("/api/bids/my", bids.GetMyBidsHandler)
	// список предложений тендера
	router.Get("/api/bids/{tenderId}/list", bids.GetBidsTenderHandler)
	// получить статус bid
	router.Get("/api/bids/{bidId}/status", bids.GetBidStatusHandler)
	// изменить статус bid
	router.Put("/api/bids/{bidId}/status", bids.ChangeBidStatusHandler)
	// редактировать bid
	router.Patch("/api/bids/{bidId}/edit", bids.ChangeBidHandler)
	// откатить версию
	router.Put("/api/bids/{bidId}/rollback/{version}", bids.RollbackBidVersionHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), router); err != nil {
		logger.LogError("Ошибка запуска сервера (main)", err)
		os.Exit(1)
	}

	// закрываем БД
	defer Mng.Db.Close()

}
