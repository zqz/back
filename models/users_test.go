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

func TestUserFindByLogin(t *testing.T) {
	TruncateUsers()

	a := assert.New(t)

	var u *User
	u, _ = UserFindByLogin("foo", "bar")

	a.Nil(u, "No users obviously shouldnt work")

	u = &User{
		FirstName: "Foo",
		LastName:  "Last",
		Username:  "foo",
		Password:  "bar",
		Email:     "foo@bar.com",
	}

	u.Save()

	u, _ = UserFindByLogin("foo", "wrongpass")
	a.Nil(u, "Wrong pass should not work")

	u, _ = UserFindByLogin("foo", "bar")
	a.NotNil(u, "Correct username/pass works")
}

func TestUserNameValid(t *testing.T) {
	TruncateUsers()

	a := assert.New(t)

	a.Equal(UserNameValid("a23"), false, "Should not allow usenames < 4 chars")
	a.Equal(UserNameValid("abc3456789012345"), false,
		"Should not allow usenames > 14 chars",
	)
	a.Equal(UserNameValid("世界"), false, "No tricking unicode")
	a.Equal(UserNameValid("世界ada"), true, "No tricking unicode")

	a.Equal(UserNameValid("foobar"), true)

	u := &User{
		FirstName: "Foo",
		LastName:  "Last",
		Username:  "foobar",
		Password:  "bar",
		Email:     "foo@bar.com",
	}

	u.Save()

	a.Equal(UserNameValid("foobar"), false, "Cant use previously taken name")
}
