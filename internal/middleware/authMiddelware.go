package middleware

import (
	"context"
	"fmt"
	"net/http"
	"social/internal/db"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	db := database.DB
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session ID from the cookie
		cookie, err := r.Cookie("sessionID")
		fmt.Println(cookie)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Get the user ID based on the session ID from the database
		var userID string
		err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ? AND expiration_time > CURRENT_TIMESTAMP", cookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Pass the user ID to the next handler
		r = SetUserIDInContext(r, userID)

		fmt.Println(userID)
		// Call the next handler
		next(w, r)
	}
}

// SetUserIDInContext sets the user ID in the request context
func SetUserIDInContext(r *http.Request, userID string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "userID", userID))
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value("userID").(string); ok {
		return userID
	}
	return ""
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logging middleware: ", r.Method, r.URL.Path)
		next(w, r)
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers dynamically based on the request's Origin header
		fmt.Println("CORS")
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Allow only specific methods for actual requests
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}