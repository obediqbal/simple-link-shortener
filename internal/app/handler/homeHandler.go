package handler

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home...")
	fmt.Fprint(w, "url shortener home, belom dibikin")
}
