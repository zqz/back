package models

import "testing"

func TestUser_Valid(t *testing.T) {
	// empty user is not valid
	u := &User{}

	if u.Valid() {
		t.Error("Empty user should not be valid.")
	}

	// factory is valid
	u = buildValidUser()

	if !u.Valid() {
		t.Error("buildValidUser should build a valid user.")
	}

	// missing email is not valid
	u.Email = ""

	if u.Valid() {
		t.Error("missing email should be invalid")
	}

	// missing name is not valid
	u = buildValidUser()
	u.FirstName = ""

	if u.Valid() {
		t.Error("missing first name should be invalid")
	}

	u = buildValidUser()
	u.LastName = ""

	if u.Valid() {
		t.Error("missing last name should be invalid")
	}

	u = buildValidUser()
	u.Username = ""

	if u.Valid() {
		t.Error("missing username should be invalid")
	}
}

func TestUser_Find(t *testing.T) {
	truncateUsers()

	_, err := UserFind("foo")
	expectNil(t, err)

	u := buildValidUser()
	u.Save()
	u2, err := UserFind(u.ID)

	expectNil(t, err)
	expectNotNil(t, u2)
}

func TestUser_Create(t *testing.T) {
	truncateUsers()

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

	expectPresent(t, "ID", u.ID)
}
