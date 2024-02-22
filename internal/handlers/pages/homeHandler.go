package pagesHandlers

import (
	"fmt"
	"net/http"
)

// HandleHome is the handler for the root path
func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
