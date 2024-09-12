package models

// Структура для отправки ошибки на клиент
type ErrorResponse struct {
	Reason string `json:"reason"`
}

// Структура для распаковки запроса на новый тендер
type RequestAddTender struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}
