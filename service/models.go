package service

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

//User struct
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//Token struct
type Token struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	UserID    uint      `json:"userID"`
	ExpiresAt int64     `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	DelatedAt time.Time `json:"deleted_at"`
}

//Save handles before save functions
func (u *User) Save(db Database) error {
	u.beforeSave()
	id, err := db.addUser(u)
	u.ID = id
	if err != nil {
		return err
	}
	u.afterSave(db)
	return nil
}

func (u *User) beforeSave() {
	u.hashPassword()
}

func (u *User) afterSave(db Database) {
	token, err := GenerateToken(u.ID)
	if err != nil {
		return
	}
	token.Save(db)
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

//CheckPasswordEqual compares passwords
func (u *User) CheckPasswordEqual(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Save create token
func (t *Token) Save(database Database) error {
	err := database.addToken(t)
	return err
}

func (t *Token) isValid() bool {
	now := time.Now().Unix()
	if t.ExpiresAt < now {
		return false
	}
	return true
}
