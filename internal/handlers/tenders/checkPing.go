package tenders

import (
	"fmt"
	"net/http"
)

// обработчик для проверки работоспособности сервера (api/ping)
func CheckPing(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем Content-Type как text/plain
	w.Header().Set("Content-Type", "text/plain")
	// Устанавливаем статус 200
	w.WriteHeader(http.StatusOK)
	// Отправляем "ok" в ответ
	fmt.Fprintln(w, "ok")
}
