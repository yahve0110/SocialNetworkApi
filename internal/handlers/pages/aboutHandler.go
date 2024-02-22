package pagesHandlers

import (
	"fmt"
	"net/http"
)

// HandleAbout is the handler for the "/about" path
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About")
}
