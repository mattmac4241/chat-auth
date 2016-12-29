package service

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // needed
)

var DB *sql.DB

type Database interface {
	addToken(token *Token) error
	addUser(user *User) (uint, error)
	getUserByUsername(username string) (User, error)
	getTokenByKey(key string) (Token, error)
	getTokenByUserID(userID uint) (Token, error)
	redisGetValue(key string) (string, error)
	redisSetValue(key, value string, seconds time.Duration) error
}

type dataHandler struct{}

func (d *dataHandler) addToken(token *Token) error {
	var lastInsertID int
	err := DB.QueryRow("INSERT INTO tokens (key, user_id, expires_at) VALUES($1, $2, $3) returning id;", token.Key, token.UserID, token.ExpiresAt).Scan(&lastInsertID)
	return err
}

func (d *dataHandler) addUser(user *User) (uint, error) {
	var lastInsertID uint
	err := DB.QueryRow("INSERT INTO users (username, password, email) VALUES($1, $2, $3) returning id;", user.Username, user.Password, user.Email).Scan(&lastInsertID)
	return lastInsertID, err
}

func (d *dataHandler) getUserByUsername(username string) (User, error) {
	var user User
	err := DB.QueryRow("SELECT ID, USERNAME, PASSWORD FROM USERS WHERE username=$1;", username).Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func (d *dataHandler) getTokenByKey(key string) (Token, error) {
	var token Token
	err := DB.QueryRow("SELECT KEY, CREATED_AT, EXPIRES_AT, USER_ID FROM TOKENS WHERE key=$1;", key).Scan(&token.Key, &token.CreatedAt, &token.ExpiresAt, &token.UserID)
	fmt.Println(err)
	return token, err
}

func (d *dataHandler) getTokenByUserID(userID uint) (Token, error) {
	var token Token
	err := DB.QueryRow("SELECT KEY, CREATED_AT, EXPIRES_AT, USER_ID FROM TOKENS WHERE user_id=$1 ORDER BY CREATED_AT DESC LIMIT 1;", userID).Scan(&token.Key, &token.CreatedAt, &token.ExpiresAt, &token.UserID)
	return token, err
}

func (d *dataHandler) redisGetValue(key string) (string, error) {
	return REDIS.Get(key).Result()
}

func (d *dataHandler) redisSetValue(key, value string, seconds time.Duration) error {
	return REDIS.Set(key, value, seconds).Err()
}

//InitDatabase setup db connection
func InitDatabase(dbinfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbinfo+" sslmode=disable")
	return db, err
}
