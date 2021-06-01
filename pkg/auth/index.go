package auth

import (
	"fmt"
	"net/http"
)

type LoginHandler struct {
	Name string
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Write to the response write directly.
	fmt.Fprintf(w, "Hello %s", h.Name)
	// renderIndex(w)
}
