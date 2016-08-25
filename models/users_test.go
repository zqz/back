package models

import (
	"testing"
	"reflect"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/boil/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testUsers(t *testing.T) {
	t.Parallel()

	query := Users(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = user.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Users(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := UserSlice{user}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := UserExists(tx, user.ID)
	if err != nil {
		t.Errorf("Unable to check if User exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected UserExistsG to return true, but got false.")
	}
}

func testUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	userFound, err := UserFind(tx, user.ID)
	if err != nil {
		t.Error(err)
	}

	if userFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Users(tx).Bind(user); err != nil {
		t.Error(err)
	}
}

func testUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Users(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	userOne := &User{}
	userTwo := &User{}
	if err = randomize.Struct(seed, userOne, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err = randomize.Struct(seed, userTwo, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = userOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = userTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Users(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	userOne := &User{}
	userTwo := &User{}
	if err = randomize.Struct(seed, userOne, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err = randomize.Struct(seed, userTwo, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = userOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = userTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

var userDBTypes = map[string]string{"FirstName": "character varying", "Username": "character varying", "Email": "character varying", "Hash": "character varying", "CreatedAt": "timestamp without time zone", "Banned": "boolean", "ID": "uuid", "LastName": "character varying", "Phone": "character varying", "UpdatedAt": "timestamp without time zone"}

func testUsersInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	var o User
	o = User{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &o, userDBTypes, true); err != nil {
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

func testUsersSliceInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	o := make(UserSlice, 3)

	seed := randomize.NewSeed()
	for i := range o {
		o[i] = &User{}
		if err = randomize.Struct(seed, o[i], userDBTypes, true); err != nil {
			t.Errorf("Could not randomize struct: %s", err)
		}
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(userPrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(userPrimaryKeyColumns)*3, len(args))
	}

	argC := 0
	for i := 0; i < 3; i++ {

		if o[i].ID != args[argC] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
		argC++
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

func testUsersHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &User{}
	o := &User{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, userDBTypes, false); err != nil {
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
}

func testUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx, userColumns...); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}



func testUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = user.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := UserSlice{user}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Users(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

func testUsersUpdate(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, user, userDBTypes, true, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	// If table only contains primary key columns, we need to pass
	// them into a whitelist to get a valid test result,
	// otherwise the Update method will error because it will not be able to
	// generate a whitelist (due to it excluding primary key columns).
	if strmangle.StringSliceMatch(userColumns, userPrimaryKeyColumns) {
		if err = user.Update(tx, userPrimaryKeyColumns...); err != nil {
			t.Error(err)
		}
	} else {
		if err = user.Update(tx); err != nil {
			t.Error(err)
		}
	}
}

func testUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, user, userDBTypes, true, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(userColumns, userPrimaryKeyColumns) {
		fields = userColumns
	} else {
		fields = strmangle.SetComplement(
			userColumns,
			userPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(user))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := UserSlice{user}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}

func testUsersUpsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	user := User{}
	if err = randomize.Struct(seed, &user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert User: %s", err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &user, userDBTypes, false, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = user.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert User: %s", err)
	}

	count, err = Users(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

