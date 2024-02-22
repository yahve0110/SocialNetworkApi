package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/post/")
	fmt.Fprintf(w, "Post ID: %s", userID)
}
