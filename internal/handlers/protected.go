package handlers

import (
	"fmt"
	"net/http"
)

func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the protected route! You are authorized.")
}
