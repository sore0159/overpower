package main

import (
	"fmt"
	"net/http"
)

func (h *Handler) pageAuthCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "AUTH CREATE PAGE")
}
