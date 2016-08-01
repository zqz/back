package models

import (
	"testing"
	"reflect"
	"time"

	"gopkg.in/nullbio/null.v4"
	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
	"github.com/nullbio/sqlboiler/strmangle"
)

func TestUsers(t *testing.T) {
	var err error

	o := make(UserSlice, 2)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	// insert two random objects to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Delete all rows to give a clean slate
	err = UsersG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Users: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	var c int64
	c, err = UsersG().Count()

	if c != 0 {
		t.Errorf("Expected users table to be empty, but got %d rows", c)
	}

	o = make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Ensure Count is valid
	c, err = UsersG().Count()
	if c != 3 {
		t.Errorf("Expected users table to have 3 rows, but got %d", c)
	}

	// Attempt to retrieve all objects
	res, err := UsersG().All()
	if err != nil {
		t.Errorf("Unable to retrieve all Users, err: %s", err)
	}

	if len(res) != 3 {
		t.Errorf("Expected 3 User rows, got %d", len(res))
	}

	usersDeleteAllRows(t)
}

func usersDeleteAllRows(t *testing.T) {
	// Delete all rows to give a clean slate
	err := UsersG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Users: %s", err)
	}
}

func TestUsersQueryDeleteAll(t *testing.T) {
	var err error
	var c int64

	// Start from a clean slate
	usersDeleteAllRows(t)

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = UsersG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after ObjDeleteAllRows() call, but got %d rows", c)
	}

	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Test DeleteAll() query function
	err = UsersG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Users: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = UsersG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after Obj().DeleteAll() call, but got %d rows", c)
	}
}

func TestUsersSliceDeleteAll(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test DeleteAll
	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	// test DeleteAll slice function
	if err = o.DeleteAllG(); err != nil {
		t.Errorf("Unable to objSlice.DeleteAll(): %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = UsersG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after objSlice.DeleteAll() call, but got %d rows", c)
	}
}

func TestUsersDelete(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test Delete
	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	o[0].DeleteG()

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = UsersG().Count()

	if c != 2 {
		t.Errorf("Expected 2 rows after obj.Delete() call, but got %d rows", c)
	}

	o[1].DeleteG()
	o[2].DeleteG()

	// Check number of rows in table to ensure Delete worked for all rows
	c, err = UsersG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after all obj.Delete() calls, but got %d rows", c)
	}
}

func TestUsersFind(t *testing.T) {
	var err error

	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(UserSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = UserFindG(o[i].ID)
		userCompareVals(o[i], j[i], t)
	}

	f, err := UserFindG(o[0].ID, userPrimaryKeyColumns...)

	if o[0].ID != f.ID {
		t.Errorf("Expected primary key values to match, ID did not match")
	}

	colsWithoutPrimKeys := boil.SetComplement(userColumns, userPrimaryKeyColumns)
	fRef := reflect.ValueOf(f).Elem()
	for _, v := range colsWithoutPrimKeys {
		val := fRef.FieldByName(v)
		if val.IsValid() {
			t.Errorf("Expected all other columns to be zero value, but column %s was %#v", v, val.Interface())
		}
	}

	usersDeleteAllRows(t)
}

func TestUsersBind(t *testing.T) {
	var err error

	o := User{}
	if err = boil.RandomizeStruct(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert User:\n%#v\nErr: %s", o, err)
	}

	j := User{}

	err = UsersG(qm.Where(`"id"=$1`, o.ID)).Bind(&j)
	if err != nil {
		t.Errorf("Unable to call Bind on User single object: %s", err)
	}

	userCompareVals(&o, &j, t)

	// insert 3 rows, attempt to bind into slice
	usersDeleteAllRows(t)

	y := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&y, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(y); i++ {
		err = y[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", y[i], err)
		}
	}

	k := UserSlice{}
	err = UsersG().Bind(&k)
	if err != nil {
		t.Errorf("Unable to call Bind on User slice of objects: %s", err)
	}

	if len(k) != 3 {
		t.Errorf("Expected 3 results, got %d", len(k))
	}

	for i := 0; i < len(y); i++ {
		userCompareVals(y[i], k[i], t)
	}

	usersDeleteAllRows(t)
}

func TestUsersOne(t *testing.T) {
	var err error

	o := User{}
	if err = boil.RandomizeStruct(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert User:\n%#v\nErr: %s", o, err)
	}

	j, err := UsersG().One()
	if err != nil {
		t.Errorf("Unable to fetch One User result:\n#%v\nErr: %s", j, err)
	}

	userCompareVals(&o, j, t)

	usersDeleteAllRows(t)
}

func TestUsersAll(t *testing.T) {
	var err error

	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	j, err := UsersG().All()
	if err != nil {
		t.Errorf("Unable to fetch All User results: %s", err)
	}

	if len(j) != 3 {
		t.Errorf("Expected 3 results, got %d", len(j))
	}

	for i := 0; i < len(o); i++ {
		userCompareVals(o[i], j[i], t)
	}

	usersDeleteAllRows(t)
}

func TestUsersCount(t *testing.T) {
	var err error

	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	// insert random columns to test Count
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	c, err := UsersG().Count()
	if err != nil {
		t.Errorf("Unable to count query User: %s", err)
	}

	if c != 3 {
		t.Errorf("Expected 3 results from count User, got %d", c)
	}

	usersDeleteAllRows(t)
}

var userDBTypes = map[string]string{"LastName": "character varying", "Email": "character varying", "CreatedAt": "timestamp without time zone", "Banned": "boolean", "ID": "uuid", "Username": "character varying", "Phone": "character varying", "UpdatedAt": "timestamp without time zone", "FirstName": "character varying"}

func userCompareVals(o *User, j *User, t *testing.T) {
	if j.ID != o.ID {
		t.Errorf("Expected id columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.ID, j.ID)
	}

	if j.FirstName != o.FirstName {
		t.Errorf("Expected first_name columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.FirstName, j.FirstName)
	}

	if j.LastName != o.LastName {
		t.Errorf("Expected last_name columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.LastName, j.LastName)
	}

	if j.Username != o.Username {
		t.Errorf("Expected username columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Username, j.Username)
	}

	if j.Phone != o.Phone {
		t.Errorf("Expected phone columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Phone, j.Phone)
	}

	if j.Email != o.Email {
		t.Errorf("Expected email columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Email, j.Email)
	}

	if o.CreatedAt.Format("02/01/2006") != j.CreatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time created_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.CreatedAt.Format("02/01/2006"), j.CreatedAt.Format("02/01/2006"))
	}

	if o.UpdatedAt.Format("02/01/2006") != j.UpdatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time updated_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.UpdatedAt.Format("02/01/2006"), j.UpdatedAt.Format("02/01/2006"))
	}

	if j.Banned != o.Banned {
		t.Errorf("Expected banned columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Banned, j.Banned)
	}
}

func TestUsersInPrimaryKeyArgs(t *testing.T) {
	var err error
	var o User
	o = User{}

	if err = boil.RandomizeStruct(&o, userDBTypes, true); err != nil {
		t.Errorf("Could not randomize struct: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(userPrimaryKeyColumns) {
		t.Errorf("Expected args to be len %d, but got %d", len(userPrimaryKeyColumns), len(args))
	}

	if o.ID != args[0] {
		t.Errorf("Expected args[0] to be value of o.ID, but got %#v", args[0])
	}
}

func TestUsersSliceInPrimaryKeyArgs(t *testing.T) {
	var err error
	o := make(UserSlice, 3)

	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Could not randomize slice: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(userPrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(userPrimaryKeyColumns)*3, len(args))
	}

	for i := 0; i < len(userPrimaryKeyColumns)*3; i++ {

		if o[i].ID != args[i] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
	}
}

func userBeforeCreateHook(o *User) error {
	*o = User{}
	return nil
}

func userAfterCreateHook(o *User) error {
	*o = User{}
	return nil
}

func userBeforeUpdateHook(o *User) error {
	*o = User{}
	return nil
}

func userAfterUpdateHook(o *User) error {
	*o = User{}
	return nil
}

func TestUsersHooks(t *testing.T) {
	var err error

	empty := &User{}
	o := &User{}

	if err = boil.RandomizeStruct(o, userDBTypes, false); err != nil {
		t.Errorf("Unable to randomize User object: %s", err)
	}

	UserAddHook(boil.HookBeforeCreate, userBeforeCreateHook)
	if err = o.doBeforeCreateHooks(); err != nil {
		t.Errorf("Unable to execute doBeforeCreateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeCreateHook function to empty object, but got: %#v", o)
	}

	userBeforeCreateHooks = []UserHook{}
	usersDeleteAllRows(t)
}

func TestUsersInsert(t *testing.T) {
	var err error

	var errs []error
	_ = errs

	emptyTime := time.Time{}.String()
	_ = emptyTime

	nullTime := null.NewTime(time.Time{}, true)
	_ = nullTime

	o := make(UserSlice, 3)
	if err = boil.RandomizeSlice(&o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert User:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(UserSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = UserFindG(o[i].ID)
		userCompareVals(o[i], j[i], t)
	}

	usersDeleteAllRows(t)

	item := &User{}
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item User:\n%#v\nErr: %s", item, err)
	}

	for _, c := range userAutoIncrementColumns {
		// Ensure the auto increment columns are returned in the object.
		if errs = boil.IsZeroValue(item, false, c); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected auto-increment columns to be greater than 0, err: %s\n", e)
			}
		}
	}

	defaultValues := []interface{}{false}

	// Ensure the simple default column values are returned correctly.
	if len(userColumnsWithSimpleDefault) > 0 && len(defaultValues) > 0 {
		if len(userColumnsWithSimpleDefault) != len(defaultValues) {
			t.Fatalf("Mismatch between slice lengths: %d, %d", len(userColumnsWithSimpleDefault), len(defaultValues))
		}

		if errs = boil.IsValueMatch(item, userColumnsWithSimpleDefault, defaultValues); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected default value to match column value, err: %s\n", e)
			}
		}
	}

	regularCols := []string{"first_name", "last_name", "username", "phone", "email", "created_at", "updated_at"}

	// Ensure the non-defaultvalue columns and non-autoincrement columns are stored correctly as zero or null values.
	for _, c := range regularCols {
		rv := reflect.Indirect(reflect.ValueOf(item))
		field := rv.FieldByName(strmangle.TitleCase(c))

		zv := reflect.Zero(field.Type()).Interface()
		fv := field.Interface()

		if !reflect.DeepEqual(zv, fv) {
			t.Errorf("Expected column %s to be zero value, got: %v, wanted: %v", c, fv, zv)
		}
	}

	item = &User{}

	wl, rc := item.generateInsertColumns()
	if !reflect.DeepEqual(rc, userColumnsWithDefault) {
		t.Errorf("Expected return columns to contain all columns with default values:\n\nGot: %v\nWanted: %v", rc, userColumnsWithDefault)
	}

	if !reflect.DeepEqual(wl, userColumnsWithoutDefault) {
		t.Errorf("Expected whitelist to contain all columns without default values:\n\nGot: %v\nWanted: %v", wl, userColumnsWithoutDefault)
	}

	if err = boil.RandomizeStruct(item, userDBTypes, false); err != nil {
		t.Errorf("Unable to randomize item: %s", err)
	}

	wl, rc = item.generateInsertColumns()
	if len(rc) > 0 {
		t.Errorf("Expected return columns to contain no columns:\n\nGot: %v", rc)
	}

	if !reflect.DeepEqual(wl, userColumns) {
		t.Errorf("Expected whitelist to contain all columns values:\n\nGot: %v\nWanted: %v", wl, userColumns)
	}

	usersDeleteAllRows(t)
}



func TestUsersSelect(t *testing.T) {
	// Only run this test if there are ample cols to test on
	if len(userAutoIncrementColumns) == 0 {
		return
	}

	var err error

	x := &struct {
	}{}

	item := User{}

	blacklistCols := boil.SetMerge(userAutoIncrementColumns, userPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, userDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert item User:\n%#v\nErr: %s", item, err)
	}

	err = UsersG(qm.Select(userAutoIncrementColumns...), qm.Where(`"id"=$1`, item.ID)).Bind(x)
	if err != nil {
		t.Errorf("Unable to select insert results with bind: %s", err)
	}

	usersDeleteAllRows(t)
}

func TestUsersUpdate(t *testing.T) {
	var err error

	item := User{}
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item User:\n%#v\nErr: %s", item, err)
	}

	blacklistCols := boil.SetMerge(userAutoIncrementColumns, userPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, userDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	whitelist := boil.SetComplement(userColumns, userAutoIncrementColumns)
	if err = item.UpdateG(whitelist...); err != nil {
		t.Errorf("Unable to update User: %s", err)
	}

	var j *User
	j, err = UserFindG(item.ID)
	if err != nil {
		t.Errorf("Unable to find User row: %s", err)
	}

	userCompareVals(&item, j, t)

	wl := item.generateUpdateColumns("test")
	if len(wl) != 1 && wl[0] != "test" {
		t.Errorf("Expected generateUpdateColumns whitelist to match expected whitelist")
	}

	wl = item.generateUpdateColumns()
	if len(wl) == 0 && len(userColumnsWithoutDefault) > 0 {
		t.Errorf("Expected generateUpdateColumns to build a whitelist for User, but got 0 results")
	}

	usersDeleteAllRows(t)
}

