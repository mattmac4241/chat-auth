package service

import (
	"errors"
	"testing"
	"time"
)

type testDatabase struct {
	users  []User
	tokens []Token
}

func (t *testDatabase) addToken(token *Token) error {
	t.tokens = append(t.tokens, *token)
	return nil
}

func (t *testDatabase) addUser(user *User) (uint, error) {
	t.users = append(t.users, *user)
	return 1, nil
}

func (t *testDatabase) getUserByUsername(username string) (User, error) {
	for _, user := range t.users {
		if username == user.Username {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}

func (t *testDatabase) getTokenByKey(key string) (Token, error) {
	for _, token := range t.tokens {
		if token.Key == key {
			return token, nil
		}
	}
	return Token{}, errors.New("Token not found")
}

func (t *testDatabase) getTokenByUserID(userID uint) (Token, error) {
	for _, token := range t.tokens {
		if token.UserID == userID {
			return token, nil
		}
	}

	return Token{}, errors.New("Token not found")
}

func TestUserSave(t *testing.T) {
	database := &testDatabase{}
	user := User{Username: "Testname", Password: "testpassword", Email: "test@mail.com"}
	unhashedPassword := user.Password
	user.Save(database)
	if len(user.Password) <= 0 || user.Password == unhashedPassword {
		t.Errorf("Failed to hash password")
		return
	}

	if len(database.tokens) != 1 {
		t.Errorf("Failed to create token")
		return
	}

	if database.tokens[0].UserID != 1 {
		t.Errorf("Failed to set userID")
	}
}

func TestUserBeforeSave(t *testing.T) {
	user := User{Username: "Testname", Password: "testpassword", Email: "test@mail.com"}
	unhashedPassword := user.Password
	user.beforeSave()

	if len(user.Password) <= 0 || user.Password == unhashedPassword {
		t.Errorf("Failed to hash password")
		return
	}
}

func TestUserAfterSave(t *testing.T) {
	database := &testDatabase{}
	user := User{Username: "Testname", Password: "testpassword", Email: "test@mail.com", ID: 1}
	user.afterSave(database)

	if len(database.tokens) != 1 {
		t.Errorf("Failed to create token")
		return
	}

	if database.tokens[0].UserID != 1 {
		t.Errorf("Failed to set userID")
	}
}

func TestHashPassword(t *testing.T) {
	user := User{Username: "Testname", Password: "testpassword", Email: "test@mail.com", ID: 1}
	unhashedPassword := user.Password
	err := user.hashPassword()
	if err != nil {
		t.Errorf("Failed to hash password")
		return
	}

	if len(user.Password) <= 0 || user.Password == unhashedPassword {
		t.Errorf("Passwords still match")
		return
	}
}

func TestCheckPasswordEqual(t *testing.T) {
	user := User{Username: "Testname", Password: "testpassword", Email: "test@mail.com", ID: 1}
	oldPassword := user.Password
	user.hashPassword()
	isEqual := user.CheckPasswordEqual(oldPassword)
	if isEqual == false {
		t.Errorf("Password should be equal")
	}

	notEqual := user.CheckPasswordEqual("notpassword")
	if notEqual == true {
		t.Errorf("Password should not equal")
	}
}

func TestTokenSave(t *testing.T) {
	token := Token{Key: "testkey"}
	database := &testDatabase{}

	token.Save(database)

	if len(database.tokens) == 0 || database.tokens[0].Key != token.Key {
		t.Error("Token did not save")
	}
}

func TestTokenIsValid(t *testing.T) {
	expiresAt := time.Now().AddDate(0, 2, 0).Unix()
	validToken := Token{Key: "testtoken", UserID: 1, ExpiresAt: expiresAt}
	invalidToken := Token{Key: "testtoken2", UserID: 1, ExpiresAt: time.Now().AddDate(0, 0, -1).Unix()}

	if validToken.isValid() == false {
		t.Error("Expected token to be valid")
	}
	if invalidToken.isValid() {
		t.Error("Expected token to be invalid")
	}
}
