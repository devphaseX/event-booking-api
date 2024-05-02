package models

import (
	"github.com/devphaseX/event-booking-api/db"
	"github.com/devphaseX/event-booking-api/encrypt/password"
)

type User struct {
	ID           int64  `json:"id"`
	Email        string `binding:"required" json:"email"`
	Password     string `binding:"required" json:"-"`
	PasswordSalt string `binding:"required" json:"-"`
}

type RegisterUserInput struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type SignInUserInput struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

var argon = password.NewArgon2idHash(3, 32, 12288, 1, 64)

func (u *User) Save() error {
	query := `
		INSERT INTO users(email, password, passwordSalt) VALUES (?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	hash, err := argon.GenerateHash([]byte(u.Password), []byte(""))

	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Email, hash.Hash, hash.Salt)

	if err != nil {
		return err
	}

	userId, err := res.LastInsertId()

	if err != nil {
		return err
	}

	u.ID = userId

	return nil
}

func (u *User) ComparePassword(cred SignInUserInput) error {
	return argon.Compare([]byte(u.Password), []byte(u.PasswordSalt), []byte(cred.Password))
}

func FindUserById(id int64) (*User, error) {
	query := `
		SELECT id, password, passwordSalt FROM users
		where id = ?
	`

	row := db.DB.QueryRow(query, id)

	var userWithId User
	err := row.Scan(&userWithId.ID, &userWithId.Password, &userWithId.PasswordSalt)

	if err != nil {
		return nil, err
	}

	return &userWithId, nil
}

func FindUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, password, passwordSalt FROM users
		where email = ?
	`

	row := db.DB.QueryRow(query, email)

	var userWithId User
	err := row.Scan(&userWithId.ID, &userWithId.Password, &userWithId.PasswordSalt)

	if err != nil {
		return nil, err
	}

	return &userWithId, nil
}
