package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
	"gopkg.in/nullbio/null.v4"
)

// User is an object representing the database table.
type User struct {
	ID        string      `db:"user_id" json:"id"`
	FirstName null.String `db:"user_first_name" json:"first_name"`
	LastName  null.String `db:"user_last_name" json:"last_name"`
	Username  string      `db:"user_username" json:"username"`
	Phone     null.String `db:"user_phone" json:"phone"`
	Email     null.String `db:"user_email" json:"email"`
	Hash      null.String `db:"user_hash" json:"hash"`
	CreatedAt time.Time   `db:"user_created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"user_updated_at" json:"updated_at"`
	Banned    null.Bool   `db:"user_banned" json:"banned"`
}

var (
	userColumns                  = []string{"id", "first_name", "last_name", "username", "phone", "email", "hash", "created_at", "updated_at", "banned"}
	userColumnsWithoutDefault    = []string{"first_name", "last_name", "username", "phone", "email", "hash", "created_at", "updated_at"}
	userColumnsWithDefault       = []string{"id", "banned"}
	userColumnsWithSimpleDefault = []string{"banned"}
	userPrimaryKeyColumns        = []string{"id"}
	userAutoIncrementColumns     = []string{}
	userAutoIncPrimaryKey        = ""
)

type (
	UserSlice []*User
	UserHook  func(*User) error

	userQuery struct {
		*boil.Query
	}
)

var userBeforeCreateHooks []UserHook
var userBeforeUpdateHooks []UserHook
var userAfterCreateHooks []UserHook
var userAfterUpdateHooks []UserHook

// doBeforeCreateHooks executes all "before create" hooks.
func (o *User) doBeforeCreateHooks() (err error) {
	for _, hook := range userBeforeCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *User) doBeforeUpdateHooks() (err error) {
	for _, hook := range userBeforeUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterCreateHooks executes all "after create" hooks.
func (o *User) doAfterCreateHooks() (err error) {
	for _, hook := range userAfterCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *User) doAfterUpdateHooks() (err error) {
	for _, hook := range userAfterUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

func UserAddHook(hookPoint boil.HookPoint, userHook UserHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		userBeforeCreateHooks = append(userBeforeCreateHooks, userHook)
	case boil.HookBeforeUpdate:
		userBeforeUpdateHooks = append(userBeforeUpdateHooks, userHook)
	case boil.HookAfterCreate:
		userAfterCreateHooks = append(userAfterCreateHooks, userHook)
	case boil.HookAfterUpdate:
		userAfterUpdateHooks = append(userAfterUpdateHooks, userHook)
	}
}

// One returns a single user record from the query.
func (q userQuery) One() (*User, error) {
	o := &User{}

	boil.SetLimit(q.Query, 1)

	res := boil.ExecQueryOne(q.Query)
	err := boil.BindOne(res, boil.Select(q.Query), o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for users: %s", err)
	}

	return o, nil
}

// OneP returns a single user record from the query, and panics on error.
func (q userQuery) OneP() *User {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all User records from the query.
func (q userQuery) All() (UserSlice, error) {
	var o UserSlice

	res, err := boil.ExecQueryAll(q.Query)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute an all query for users: %s", err)
	}
	defer res.Close()

	err = boil.BindAll(res, boil.Select(q.Query), &o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to User slice: %s", err)
	}

	return o, nil
}

// AllP returns all User records from the query, and panics on error.
func (q userQuery) AllP() UserSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// Count returns the count of all User records in the query.
func (q userQuery) Count() (int64, error) {
	var count int64

	boil.SetCount(q.Query)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("models: failed to count users rows: %s", err)
	}

	return count, nil
}

// CountP returns the count of all User records in the query, and panics on error.
func (q userQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}



// UsersAll retrieves all records.
func Users(mods ...qm.QueryMod) userQuery {
	return UsersX(boil.GetDB(), mods...)
}

// UsersX retrieves all the records using an executor.
func UsersX(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	mods = append(mods, qm.From("users"))
	return userQuery{NewQueryX(exec, mods...)}
}


// UserFind retrieves a single record by ID.
func UserFind(id string, selectCols ...string) (*User, error) {
	return UserFindX(boil.GetDB(), id, selectCols...)
}

// UserFindP retrieves a single record by ID, and panics on error.
func UserFindP(id string, selectCols ...string) *User {
	o, err := UserFindX(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// UserFindX retrieves a single record by ID with an executor.
func UserFindX(exec boil.Executor, id string, selectCols ...string) (*User, error) {
	user := &User{}

	mods := []qm.QueryMod{
		qm.Select(selectCols...),
		qm.From("users"),
		qm.Where(`"id"=$1`, id),
	}

	q := NewQueryX(exec, mods...)

	err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers(user, selectCols...)...)

	if err != nil {
		return nil, fmt.Errorf("models: unable to select from users: %v", err)
	}

	return user, nil
}

// UserFindXP retrieves a single record by ID with an executor, and panics on error.
func UserFindXP(exec boil.Executor, id string, selectCols ...string) *User {
	o, err := UserFindX(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// Insert a single record.
func (o *User) Insert(whitelist ...string) error {
	return o.InsertX(boil.GetDB(), whitelist...)
}

// InsertP a single record, and panics on error.
func (o *User) InsertP(whitelist ...string) {
	if err := o.InsertX(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertX a single record using an executor.
func (o *User) InsertX(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no users provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO users ("%s") VALUES (%s)`, strings.Join(wl, `","`), boil.GenerateParamFlags(len(wl), 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.FirstName, o.LastName, o.Username, o.Phone, o.Email, o.Hash, o.CreatedAt, o.UpdatedAt, o.Banned)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into users: %s", err)
	}

	if err := o.doAfterCreateHooks(); err != nil {
		return err
	}

	return nil
}

// InsertXP a single record using an executor, and panics on error.
func (o *User) InsertXP(exec boil.Executor, whitelist ...string) {
	if err := o.InsertX(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
func (o *User) generateInsertColumns(whitelist ...string) ([]string, []string) {
	var wl []string

	wl = append(wl, whitelist...)
	if len(whitelist) == 0 {
		wl = append(wl, userColumnsWithoutDefault...)
	}

	wl = append(boil.NonZeroDefaultSet(userColumnsWithDefault, o), wl...)
	wl = boil.SortByKeys(userColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(userColumnsWithDefault, wl)

	return wl, rc
}


// Update a single User record.
// Update takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
func (o *User) Update(whitelist ...string) error {
	return o.UpdateX(boil.GetDB(), whitelist...)
}

// Update a single User record.
// UpdateP takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
// Panics on error.
func (o *User) UpdateP(whitelist ...string) {
	if err := o.UpdateX(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateX uses an executor to update the User.
func (o *User) UpdateX(exec boil.Executor, whitelist ...string) error {
	return o.UpdateAtX(exec, o.ID, whitelist...)
}

// UpdateXP uses an executor to update the User, and panics on error.
func (o *User) UpdateXP(exec boil.Executor, whitelist ...string) {
	err := o.UpdateAtX(exec, o.ID, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAt updates the User using the primary key to find the row to update.
func (o *User) UpdateAt(id string, whitelist ...string) error {
	return o.UpdateAtX(boil.GetDB(), id, whitelist...)
}

// UpdateAtP updates the User using the primary key to find the row to update. Panics on error.
func (o *User) UpdateAtP(id string, whitelist ...string) {
	if err := o.UpdateAtX(boil.GetDB(), id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAtX uses an executor to update the User using the primary key to find the row to update.
func (o *User) UpdateAtX(exec boil.Executor, id string, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE users SET %s WHERE %s`, boil.SetParamNames(wl), boil.WherePrimaryKey(len(wl)+1, "id"))
		values = boil.GetStructValues(o, wl...)
		values = append(values, o.ID)
		_, err = exec.Exec(query, values...)
	} else {
		return fmt.Errorf("models: unable to update users, could not build whitelist")
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if err != nil {
		return fmt.Errorf("models: unable to update users row: %s", err)
	}

	if err := o.doAfterUpdateHooks(); err != nil {
		return err
	}

	return nil
}

// UpdateAtXP uses an executor to update the User using the primary key to find the row to update.
// Panics on error.
func (o *User) UpdateAtXP(exec boil.Executor, id string, whitelist ...string) {
	if err := o.UpdateAtX(exec, id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with matching column names.
func (q userQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for users: %s", err)
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q userQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateUpdateColumns generates the whitelist columns for an update statement
func (o *User) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	var wl []string
	cols := userColumnsWithoutDefault
	cols = append(boil.NonZeroDefaultSet(userColumnsWithDefault, o), cols...)
	// Subtract primary keys and autoincrement columns
	cols = boil.SetComplement(cols, userPrimaryKeyColumns)
	cols = boil.SetComplement(cols, userAutoIncrementColumns)

	wl = make([]string, len(cols))
	copy(wl, cols)

	return wl
}

// Delete deletes a single User record.
// Delete will match against the primary key column to find the record to delete.
func (o *User) Delete() error {
	if o == nil {
		return errors.New("models: no User provided for deletion")
	}

	return o.DeleteX(boil.GetDB())
}

// DeleteP deletes a single User record.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *User) DeleteP() {
	if err := o.Delete(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteX deletes a single User record with an executor.
// DeleteX will match against the primary key column to find the record to delete.
func (o *User) DeleteX(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User provided for deletion")
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.From("users"),
		qm.Where(`"id"=$1`, o.ID),
	)

	query := NewQueryX(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete from users: %s", err)
	}

	return nil
}

// DeleteXP deletes a single User record with an executor.
// DeleteXP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *User) DeleteXP(exec boil.Executor) {
	if err := o.DeleteX(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows.
func (o userQuery) DeleteAll() error {
	if o.Query == nil {
		return errors.New("models: no userQuery provided for delete all")
	}

	boil.SetDelete(o.Query)

	_, err := boil.ExecQuery(o.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from users: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (o userQuery) DeleteAllP() {
	if err := o.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice.
func (o UserSlice) DeleteAll() error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}
	return o.DeleteAllX(boil.GetDB())
}

// DeleteAll deletes all rows in the slice.
func (o UserSlice) DeleteAllP() {
	if err := o.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllX deletes all rows in the slice with an executor.
func (o UserSlice) DeleteAllX(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}

	var mods []qm.QueryMod

	args := o.inPrimaryKeyArgs()
	in := boil.WherePrimaryKeyIn(len(o), "id")

	mods = append(mods,
		qm.From("users"),
		qm.Where(in, args...),
	)

	query := NewQueryX(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from user slice: %s", err)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// DeleteAllXP deletes all rows in the slice with an executor, and panics on error.
func (o UserSlice) DeleteAllXP(exec boil.Executor) {
	if err := o.DeleteAllX(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o User) inPrimaryKeyArgs() []interface{} {
	var args []interface{}
	args = append(args, o.ID)
	return args
}

func (o UserSlice) inPrimaryKeyArgs() []interface{} {
	var args []interface{}

	for i := 0; i < len(o); i++ {
		args = append(args, o[i].ID)
	}

	return args
}

