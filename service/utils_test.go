package service

import (
	"testing"
	"time"
)

func TestUserLogin(t *testing.T) {
	user := User{Username: "Testname", Password: "testpassword", Email: "test@mail.com", ID: 1}
	oldPassword := user.Password
	database := &testDatabase{}
	user.Save(database)

	_, err := UserLogin("fake", "fake", database)
	if err == nil {
		t.Error("Expected an error and empty token")
	}

	_, err = UserLogin(user.Username, "password", database)
	if err == nil {
		t.Error("Expected error because wrong password, got none")
	}

	correctPasswordToken, err := UserLogin(user.Username, oldPassword, database)
	if err != nil && correctPasswordToken != database.tokens[0] {
		t.Error("Expect token from user but got an error and wrong token")
	}
}

func TestGetTokenOrCreateNewOneIfExpired(t *testing.T) {
	var invalidUserID uint
	expiresAt := time.Now().AddDate(0, 2, 0).Unix()
	notExpiredToken := &Token{Key: "testtoken", UserID: 1, ExpiresAt: expiresAt}
	expiredToken := &Token{Key: "testtoken2", UserID: 2, ExpiresAt: time.Now().AddDate(0, 0, -1).Unix()}
	invalidUserID = 0
	database := &testDatabase{}
	database.addToken(notExpiredToken)
	database.addToken(expiredToken)

	_, err := getTokenOrCreateNewOneIfExpired(invalidUserID, database)
	if err == nil {
		t.Error("Expected error for not found token")
	}

	token, err := getTokenOrCreateNewOneIfExpired(1, database)
	if err != nil && token != database.tokens[0] {
		t.Error("Expected valid token and no error")
	}

	newToken, err := getTokenOrCreateNewOneIfExpired(2, database)
	if err != nil && len(database.tokens) != 3 && newToken.UserID != 2 {
		t.Error("Expeted no error and a newly created token")
	}

}
