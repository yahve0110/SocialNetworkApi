package notifications

import (
	"log"
	"net/http"
	database "social/internal/db"
)

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	notificationID := r.URL.Query().Get("notification_id")

	dbConnection := database.DB

	result, err := dbConnection.Exec("DELETE FROM notifications WHERE notification_id = $1", notificationID)
	if err != nil {
		log.Printf("Error deleting notification: %v", err)
		http.Error(w, "Error deleting notification", http.StatusInternalServerError)
		return
	}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification deleted successfully"))
}
