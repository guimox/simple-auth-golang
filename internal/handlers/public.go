package handlers

import (
	"fmt"
	"net/http"
)

func Public(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a public route. Everyone can access this!")
}
