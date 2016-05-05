package user

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is a foo.
type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name" valid:"required,alphanum"`
	LastName  string `json:"last_name" valid:"required,alphanum"`
	Username  string `json:"username" valid:"required"`
	Phone     string `json:"phone"`
	Email     string `json:"email" valid:"required,email"`
	Hash      string `json:"-"`
	Password  string
	Banned    bool `json:"banned"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	errors []string
}

const findByUsernameSQL = `
	SELECT
	id, first_name, last_name, phone, email, banned, created_at, updated_at
	FROM users
	WHERE
	username = $1
`

const insertSQL = `
	INSERT INTO users
	(first_name, last_name, username, phone, email, hash)
	VALUES
	($1, $2, $3, $4, $5, $6)
	RETURNING id
`

const setPasswordSQL = `
	UPDATE users
	SET hash = $2
	WHERE id = $1
`

// Create a user inside of a transaction.
func (u *User) Create(tx *sql.Tx) error {
	if len(u.Hash) == 0 {
		u.hashPassword()
	}

	err := tx.QueryRow(insertSQL,
		u.FirstName, u.LastName, u.Username, u.Phone, u.Email, u.Hash,
	).Scan(&u.ID)

	return err
}

func (u *User) hashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	u.Password = ""
	u.Hash = string(hash)
}

// FindByUsername returns a User with the specified username
func FindByUsername(tx *sql.Tx, username string) (*User, error) {
	var u User
	u.Username = username
	err := tx.QueryRow(findByUsernameSQL, username).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Phone, &u.Email, &u.Banned,
		&u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

// SetPassword changes the users password.
func (u *User) SetPassword(tx *sql.Tx, password string) bool {
	u.Password = password
	u.hashPassword()

	err, _ := tx.Exec(setPasswordSQL, u.ID, u.Hash)

	return err == nil
}
