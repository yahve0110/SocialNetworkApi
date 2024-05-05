package notifications

import (
	"log"
	"net/http"
	database "social/internal/db"
)

// DeleteNotification обрабатывает HTTP-запрос для удаления уведомления
func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор уведомления из запроса
	notificationID := r.URL.Query().Get("notification_id")

	// Получаем подключение к базе данных
	dbConnection := database.DB

	// Удаляем уведомление из базы данных по его идентификатору
	result, err := dbConnection.Exec("DELETE FROM notifications WHERE notification_id = $1", notificationID)
	if err != nil {
		log.Printf("Error deleting notification: %v", err)
		http.Error(w, "Error deleting notification", http.StatusInternalServerError)
		return
	}

	// Проверяем количество удаленных строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Error deleting notification", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("Notification with ID %s not found", notificationID)
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	// Если все прошло успешно, возвращаем статус OK и сообщение об успешном удалении
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification deleted successfully"))
}
