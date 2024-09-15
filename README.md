# Сервис обработчик тендеров и предложений (Tender Management)

    Этот проект представляет собой сервис для управления тендерами и предложениями компаний и использованием базы данных PostgreSQL.

    Структура приложения

    	1.	Основные компоненты:
	        •	Logger: Логгирует все действия и ошибки сервиса.
	        •	Config: Загружает переменные окружения из файла конфигурации.
	        •	Manager: Управляет основной логикой приложения.
	        •	Database: Работает с PostgreSQL для хранения данных.
	        •	HTTP сервер: Обрабатывает запросы API.
	    2.	API Endpoints:
	        •	GET /api/ping: Проверка доступности сервера.
            •	GET /api/tenders: Получение списка тендеров.
            •	POST /api/tenders/new: Создание нового тендера.
            •	GET /api/tenders/my: Получить тендеры пользователя.
            •	GET /api/tenders/{tenderId}/status: Получение текущего статуса тендера.
            •	PUT /api/tenders/{tenderId}/status: Изменение статуса тендера.
            •	PATCH /api/tenders/{tenderId}/edit: Редактирование тендера.
            •	PUT /api/tenders/{tenderId}/rollback/{version}: Откат версии тендера.
             •	POST /api/bids/new: Создание нового предложения.
            •	GET /api/bids/my: Получение списка ваших предложений.
            •	GET /api/bids/{tenderId}/list: Получение списка предложений для тендера.
            •	GET /api/bids/{bidId}/status: Получение текущего статуса предложения.
            •	PUT /api/bids/{bidId}/status: Изменение статуса предложения.
            •	PATCH /api/bids/{bidId}/edit: Редактирование параметров предложения.
            •	PUT /api/bids/{bidId}/rollback/{version}: Откат версии предложения.

## Docker

    Для создания и запуска образа:

        docker build -t tender-service:latest .     

        docker run --env-file ./.env -p 8080:8080 tender-service:latest

## Переменные окружения

    SERVER_ADDRESS — адрес и порт, который будет слушать HTTP сервер при запуске. Пример: 0.0.0.0:8080.
    POSTGRES_CONN — URL-строка для подключения к PostgreSQL в формате postgres://{username}:{password}@{host}:{5432}/{dbname}.
    POSTGRES_JDBC_URL — JDBC-строка для подключения к PostgreSQL в формате jdbc:postgresql://{host}:{port}/{dbname}.
    POSTGRES_USERNAME — имя пользователя для подключения к PostgreSQL.
    POSTGRES_PASSWORD — пароль для подключения к PostgreSQL.
    POSTGRES_HOST — хост для подключения к PostgreSQL (например, localhost).
    POSTGRES_PORT — порт для подключения к PostgreSQL (например, 5432).
    POSTGRES_DATABASE — имя базы данных PostgreSQL, которую будет использовать приложение.
