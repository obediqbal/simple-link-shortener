package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simplelinkshortener/internal/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if err := validate(username, password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Here's your token")
}

func validate(username string, password string) error {
	db := database.New()

	stmt, err := db.Prepare("SELECT password FROM users WHERE username = $1")
	if err != nil {
		log.Println(err)
		return err
	}

	var hashedPassword string
	row := stmt.QueryRow(username)
	if err := row.Scan(&hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Username or password is incorrect")
		}
		log.Println(err)
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("Username or password is incorrect")
	}

	return nil
}
