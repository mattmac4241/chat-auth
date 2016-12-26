package service

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//GenerateToken creates token
func GenerateToken(userID uint) (Token, error) {
	key, err := generateKey(userID)
	if err != nil {
		return Token{}, err
	}

	token := Token{
		Key:       key,
		UserID:    userID,
		ExpiresAt: getExpiresAtTime(),
	}

	return token, nil
}

func UserLogin(username, password string, database Database) (Token, error) {
	user, err := database.getUserByUsername(username)
	if err != nil {
		return Token{}, err
	}
	canLogin := user.CheckPasswordEqual(password)
	if canLogin == false {
		return Token{}, errors.New("Passwords do not match")
	}
	newToken, err := getTokenOrCreateNewOneIfExpired(user.ID, database)
	return newToken, err
}

func CheckTokenKey(key string, database Database) (Token, error) {
	token, err := database.getTokenByKey(key)
	return token, err
}

func generateKey(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userID,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenString, err
}

func getExpiresAtTime() int64 {
	now := time.Now().AddDate(0, 2, 0).Unix()
	return now
}

func getTokenOrCreateNewOneIfExpired(userID uint, database Database) (Token, error) {
	token, err := database.getTokenByUserID(userID)
	if err != nil {
		return Token{}, err
	}
	if token.isValid() {
		return token, nil
	}

	newToken, err := GenerateToken(token.UserID)
	if err != nil {
		return Token{}, err
	}
	err = database.addToken(&newToken)
	return newToken, err
}
