package handler

import (
	"fmt"
	"net/http"
)

// Handler returns a string
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It works!")
}
