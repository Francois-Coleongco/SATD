package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type UserCreds struct {
	username      string
	password_hash string
}

var (
	pool     *pgxpool.Pool
	hostname = flag.String("hostname", "localhost", "please provide a valid address")
	port     = flag.Int("port", 7777, "please provide a valid port")

	// .env initialized
	secretKey     []byte
	elasticApiKey string
	esClient      *elasticsearch.Client

	//sql

	dbHost     = "localhost" // i think this should always be localhost
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

func load_creds() error {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("WHERE THE HECK IS THE .env FILEEE")
		return err
	}

	secretKey = []byte(os.Getenv("SECRET_KEY_FOR_JWT"))

	return nil
}

func initialize_db() (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Println("pool couldn't be created")
		return pool, err
	}

	defer pool.Close()

	return pool, nil

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

	tokenString := r.Header.Get("token")

	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "missing auth token")
		return
	}

	tokenString = tokenString[len("Bearer "):]

	err := verifyToken(tokenString)

	if err != nil {
		log.Printf("received invalid token for /dashboard, error thrown: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		APIKey: elasticApiKey,
	})

	esClient.
		w.Header().Set("Content-Type", "text/event-stream")
}

func main() {
	load_creds()
	var err error
	pool, err = initialize_db()

	if err != nil {
		log.Fatalf("couldn't initialize database %s\n", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

	http.ListenAndServeTLS(*hostname+fmt.Sprint(*port), "cert.pem", "key.pem", nil)

}
