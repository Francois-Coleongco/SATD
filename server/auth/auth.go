package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type UserCreds struct {
	username      string
	password_hash string
}

var (
	secretKey []byte

	//sql

	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = ""         // no password for default postgres
	dbName     = "postgres" // default db name
)

func createToken(username string) (string, error) {
	// should be good using HS256. symmetric doesnt matter since it's the same program signing and verifying, symmetric key never leaves server

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		log.Printf("invalid token, error thrown: %s\n", err)
		return fmt.Errorf("invalid token")
	}
	return nil
}

func forwardToServerEndPoint() {
	//
}

func load_creds() error {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("WHERE THE HECK IS THE .env FILEEE")
		return err
	}

	secretKey = []byte(os.Getenv("SECRET_KEY_FOR_JWT"))

	return nil
}

func initialize_db() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	if err != nil {
		log.Println("postgres couldn't be connected")
	}

	defer db.Close()

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u UserCreds

	json.NewDecoder(r.Body).Decode(&u)

	token, err := createToken(u.username)

	if err != nil {
		log.Printf("couldn't create token for user, error thrown %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})

	if err != nil {
		log.Printf("couldn't encode token into json, error thrown: %s\n", err)
	}

}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	load_creds()
	conn, err := initialize_db()

	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

}
