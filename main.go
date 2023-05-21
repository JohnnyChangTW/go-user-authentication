package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CreateAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
}

type VerifyAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
}

// Account represents an account entity.
type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	// const (
	// 	UserName string = "root"
	// 	Password string = "password"
	// 	Addr     string = "db"
	// 	Port     int    = 3306
	// 	Database string = "account_db"
	// )
	// conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	// db, err = sql.Open("mysql", conn)
	// db, err = sql.Open("mysql", "root:password@tcp(docker.for.mac.localhost:3306)/account_db")
	db, err = sql.Open("mysql", "root:password@tcp(db:3306)/account_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the accounts table if it doesn't exist
	createTable()

	router := gin.Default()
	fmt.Println("Server started")
	router.POST("/accounts", createAccountHandler)
	router.POST("/accounts/verify", verifyAccountHandler)

	log.Fatal(router.Run(":8000"))
}

// createTable creates the account table if it doesn't exist.
func createTable() {
	fmt.Println("Creating account table")
	query :=
		`CREATE TABLE IF NOT EXISTS accounts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(32) NOT NULL,
			password VARCHAR(32) NOT NULL)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create account table: %v", err)
	}
}

func createAccountHandler(c *gin.Context) {
	var request CreateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !isUsernameValid(request.Username) {
		respondWithError(c, http.StatusBadRequest, "Invalid username")
		return
	}

	exists, err := isUsernameExists(request.Username)
	if err != nil {
		fmt.Println(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to verify account")
		return
	}

	if exists {
		respondWithError(c, http.StatusBadRequest, "Username already exists")
		return
	}

	if !isPasswordValid(request.Password) {
		respondWithError(c, http.StatusBadRequest, "Invalid password")
		return
	}

	err = insertAccount(request.Username, request.Password)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create account")
		return
	}

	respondWithJSON(c, http.StatusOK, AccountResponse{Success: true})
}

func verifyAccountHandler(c *gin.Context) {
	var request VerifyAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !isUsernameValid(request.Username) {
		respondWithError(c, http.StatusBadRequest, "Invalid username")
		return
	}

	locked, err := isAccountLocked(request.Username)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to verify account")
		return
	}

	if locked {
		respondWithError(c, http.StatusTooManyRequests, "Account locked. Please try again after 1 minute")
		return
	}

	valid, err := isPasswordValidForUsername(request.Username, request.Password)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to verify account")
		return
	}

	if !valid {
		lockAccount(request.Username)
		respondWithError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	respondWithJSON(c, http.StatusOK, VerifyResponse{Success: true})
}

func isUsernameValid(username string) bool {
	// Implement the username validation logic
	// Return true if the username is valid, otherwise false
	// Example implementation:
	return len(username) >= 3 && len(username) <= 32
}

func isUsernameExists(username string) (bool, error) {
	// Implement the username existence check logic
	// Query the database to check if the username already exists
	// Return true if the username exists, otherwise false
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE username = ?", username).Scan(&count)
	fmt.Println(err)
	fmt.Println(count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func isPasswordValid(password string) bool {
	// Implement the password validation logic
	// Return true if the password is valid, otherwise false
	// Example implementation:
	if len(password) < 8 || len(password) > 32 {
		return false
	}
	hasUpperCase := false
	hasLowerCase := false
	hasNumber := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpperCase = true
		}
		if char >= 'a' && char <= 'z' {
			hasLowerCase = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}
	return hasUpperCase && hasLowerCase && hasNumber
}

func insertAccount(username, password string) error {
	// Implement the account insertion logic
	// Insert the account details into the database
	// Example implementation:
	_, err := db.Exec("INSERT INTO accounts (username, password) VALUES (?, ?)", username, password)
	return err
}

func isAccountLocked(username string) (bool, error) {
	// Implement the account lock check logic
	// Query the database to check if the account is locked
	// Return true if the account is locked, otherwise false
	// Example implementation:
	var locked bool
	err := db.QueryRow("SELECT locked FROM accounts WHERE username = ?", username).Scan(&locked)
	if err != nil {
		return false, err
	}
	return locked, nil
}

func lockAccount(username string) error {
	// Implement the account locking logic
	// Update the account status in the database to locked
	// Example implementation:
	_, err := db.Exec("UPDATE accounts SET locked = true WHERE username = ?", username)
	return err
}

func isPasswordValidForUsername(username, password string) (bool, error) {
	// Implement the password validation for a given username logic
	// Query the database to check if the password is valid for the given username
	// Return true if the password is valid, otherwise false
	// Example implementation:
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	response := AccountResponse{
		Success: false,
		Reason:  message,
	}
	c.JSON(statusCode, response)
}

func respondWithJSON(c *gin.Context, statusCode int, payload interface{}) {
	c.JSON(statusCode, payload)
}
