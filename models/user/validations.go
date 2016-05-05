package user

import (
	"database/sql"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

const validCredentialsSQL = `
	SELECT hash
	FROM users
	WHERE username = $1
`

const usernameFreeSQL = `
	SELECT NOT EXISTS (
		SELECT 1
		FROM users
		WHERE username = $1
	)
`

// Valid checks if the user had valid Data. It assigns to the errors field if
// the User is invalid.
func (u *User) Valid() bool {
	result, err := govalidator.ValidateStruct(u)
	if err != nil {
		u.errors = strings.Split(strings.TrimRight(err.Error(), ";"), ";")
	}
	return result
}

// ValidCredentials checks if a username and password combination exists.
func ValidCredentials(tx *sql.Tx, username string, password string) bool {
	if len(username) == 0 {
		return false
	}

	if len(password) == 0 {
		return false
	}

	var hash string

	err := tx.QueryRow(validCredentialsSQL, username).Scan(&hash)

	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

const minUsernameLength = 4
const maxUsernameLength = 14

// UsernameFree returns true if the provider username can be used.
func UsernameFree(tx *sql.Tx, username string) bool {
	length := utf8.RuneCount([]byte(username))

	if length < minUsernameLength {
		return false
	}

	if length > maxUsernameLength {
		return false
	}

	var free bool
	err := tx.QueryRow(usernameFreeSQL, username).Scan(&free)

	if err != nil {
		return false
	}

	return free
}
