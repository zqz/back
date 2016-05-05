package user

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"

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

const validCredentialsSQL = `
	SELECT hash
	FROM users
	WHERE username = $1
`

const setPasswordSQL = `
	UPDATE users
	SET hash = $2
	WHERE id = $1
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

// // Valid return true or false depending on whether or not the User is valid. It
// // additionally sets the errors field on the User to provide information about
// // why the user is not valid
// func (u *User) Valid() bool {
// 	result, err := govalidator.ValidateStruct(u)
// 	if err != nil {
// 		u.errors = strings.Split(strings.TrimRight(err.Error(), ";"), ";")
// 	}
// 	return result
// }

// func (u *User) Errors() *UserError {
// 	return &UserError{u.errors}
// }

// // UserFindByAuthToken adwad
// func UserFindByAuthToken(token string) (*User, error) {
// 	if len(token) == 0 {
// 		return nil, errors.New("Token unspecified")
// 	}

// 	var u User
// 	uc := userCollection()
// 	res := uc.Find(db.Cond{"api_key": token})
// 	res.One(&u)

// 	if len(u.ID) == 0 {
// 		return nil, errors.New("Failed to find User with Auth Token Provided")
// 	}

// 	return &u, nil
// }

// // UserCount is the count of users
// func UserCount() int {
// 	uc := userCollection()
// 	cnt, _ := uc.Find().Count()
// 	return int(cnt)
// }

// // UserFind finds by id
// func UserFind(id string) (*User, error) {
// 	if len(id) == 0 {
// 		return nil, errors.New("ID unspecified")
// 	}

// 	var u User
// 	uc := userCollection()
// 	res := uc.Find(db.Cond{"id": id})
// 	res.One(&u)

// 	if len(u.ID) == 0 {
// 		return nil, errors.New("User not found")
// 	}

// 	return &u, nil
// }

// // Save the user to the database
// func (u *User) Save() bool {
// 	if len(u.ID) == 0 {
// 		return u.Create()
// 	}

// 	return u.Update()
// }

// // Update a user
// func (u *User) Update() bool {
// 	uc := userCollection()

// 	if uc == nil {
// 		return false
// 	}

// 	res := uc.Find(db.Cond{"id": u.ID})

// 	if err := res.Update(u); err != nil {
// 		log.Println("failed to update user", err.Error())
// 		return false
// 	}

// 	return true
// }

// func LoadUser(r io.Reader) *User {
// 	u := User{}
// 	json.NewDecoder(r).Decode(&u)
// 	return &u
// }

// func (u *User) String() string {
// 	buf := new(bytes.Buffer)

// 	json.NewEncoder(buf).Encode(u)

// 	return buf.String()
// }

const minUsernameLength = 4
const maxUsernameLength = 14

func UsernameFree(tx *sql.Tx, username string) bool {
	length := utf8.RuneCount([]byte(username))

	fmt.Println("length", length)

	if length < minUsernameLength {
		return false
	}

	if length > maxUsernameLength {
		return false
	}

	var free bool
	err := tx.QueryRow(usernameFreeSQL, username).Scan(&free)

	fmt.Println(err)

	if err != nil {
		return false
	}

	return free
}
