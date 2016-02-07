package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Valid(t *testing.T) {
	// empty user is not valid
	u := &User{}

	if u.Valid() {
		t.Fatal("Empty user should not be valid.")
	}

	// factory is valid
	u = buildValidUser()

	if !u.Valid() {
		t.Fatal("buildValidUser should build a valid user.")
	}

	// missing email is not valid
	u.Email = ""

	if u.Valid() {
		t.Fatal("missing email should be invalid")
	}

	// missing name is not valid
	u = buildValidUser()
	u.FirstName = ""

	if u.Valid() {
		t.Fatal("missing first name should be invalid")
	}

	u = buildValidUser()
	u.LastName = ""

	if u.Valid() {
		t.Fatal("missing last name should be invalid")
	}

	u = buildValidUser()
	u.Username = ""

	if u.Valid() {
		t.Fatal("missing username should be invalid")
	}
}

func TestUser_Find(t *testing.T) {
	TruncateUsers()

	a := assert.New(t)

	_, err := UserFind("foo")

	a.NotNil(err)
	a.Equal(err.Error(), "User not found")

	u := buildValidUser()
	u.Save()
	u2, err := UserFind(u.ID)

	a.Nil(err)
	a.NotNil(u2)
}

func TestUser_Create(t *testing.T) {
	TruncateUsers()

	a := assert.New(t)
	u := buildValidUser()

	// invalid user
	u.Username = ""
	succ := u.Save()

	if succ == true {
		t.Error("Save should have failed as user is invalid")
	}

	u.Username = "jc"
	succ = u.Save()

	if succ == false {
		t.Error("Save should have been successful")
	}

	a.NotEmpty(u.ID)
}
