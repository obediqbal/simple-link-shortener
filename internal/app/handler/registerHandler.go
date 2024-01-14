package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"simplelinkshortener/internal/pkg/auth"
	"simplelinkshortener/internal/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	exists, err := isUserExists(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	if len(password) < 6 {
		http.Error(w, "Password too short (min 6 chars)", http.StatusBadRequest)
		return
	}

	if err := register(username, password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := auth.GenerateJWT(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, tokenString)
}

func register(username string, password string) error {
	db := database.New()
	hashedPassword, err := hashPassword(password)

	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := stmt.Exec(username, hashedPassword)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Printf("Registered %d user\n", rowsAffected)
	return nil
}

func isUserExists(username string) (bool, error) {
	db := database.New()

	stmt, err := db.Prepare("SELECT username FROM users WHERE username = $1")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()

	var data string
	row := stmt.QueryRow(username)
	if err := row.Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			log.Println(err)
			return false, err
		}
	}
	return true, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashString := string(hash)

	return hashString, nil
}
