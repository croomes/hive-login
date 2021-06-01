package index

import (
	"net/http"
)

type Handler struct {
	Name string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	renderIndex(w)
}
