package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/boil/qm"
	"github.com/vattle/sqlboiler/strmangle"
)

// User is an object representing the database table.
type User struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	FirstName string    `boil:"first_name" json:"first_name" toml:"first_name" yaml:"first_name"`
	LastName  string    `boil:"last_name" json:"last_name" toml:"last_name" yaml:"last_name"`
	Username  string    `boil:"username" json:"username" toml:"username" yaml:"username"`
	Phone     string    `boil:"phone" json:"phone" toml:"phone" yaml:"phone"`
	Email     string    `boil:"email" json:"email" toml:"email" yaml:"email"`
	Hash      string    `boil:"hash" json:"hash" toml:"hash" yaml:"hash"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	Banned    bool      `boil:"banned" json:"banned" toml:"banned" yaml:"banned"`

	//Relationships *UserRelationships `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// UserRelationships are where relationships are both cached
// and eagerly loaded.
type UserRelationships struct {
}


var (
	userColumns                  = []string{"id", "first_name", "last_name", "username", "phone", "email", "hash", "created_at", "updated_at", "banned"}
	userColumnsWithoutDefault    = []string{"first_name", "last_name", "username", "phone", "email", "hash", "created_at", "updated_at"}
	userColumnsWithDefault       = []string{"id", "banned"}
	userColumnsWithSimpleDefault = []string{"banned"}
	userValidatedColumns         = []string{"id"}
	userUniqueColumns            = []string{"username"}
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
var userBeforeUpsertHooks []UserHook
var userAfterCreateHooks []UserHook
var userAfterUpdateHooks []UserHook
var userAfterUpsertHooks []UserHook

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

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *User) doBeforeUpsertHooks() (err error) {
	for _, hook := range userBeforeUpsertHooks {
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

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *User) doAfterUpsertHooks() (err error) {
	for _, hook := range userAfterUpsertHooks {
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
	case boil.HookBeforeUpsert:
		userBeforeUpsertHooks = append(userBeforeUpsertHooks, userHook)
	case boil.HookAfterCreate:
		userAfterCreateHooks = append(userAfterCreateHooks, userHook)
	case boil.HookAfterUpdate:
		userAfterUpdateHooks = append(userAfterUpdateHooks, userHook)
	case boil.HookAfterUpsert:
		userAfterUpsertHooks = append(userAfterUpsertHooks, userHook)
	}
}

// OneP returns a single user record from the query, and panics on error.
func (q userQuery) OneP() *User {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single user record from the query.
func (q userQuery) One() (*User, error) {
	o := &User{}

	boil.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for users: %s", err)
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

// All returns all User records from the query.
func (q userQuery) All() (UserSlice, error) {
	var o UserSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to User slice: %s", err)
	}

	return o, nil
}

// CountP returns the count of all User records in the query, and panics on error.
func (q userQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
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

// Exists checks if the row exists in the table, and panics on error.
func (q userQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q userQuery) Exists() (bool, error) {
	var count int64

	boil.SetCount(q.Query)
	boil.SetLimit(q.Query, 1)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("models: failed to check if users exists: %s", err)
	}

	return count > 0, nil
}



// UsersG retrieves all records.
func UsersG(mods ...qm.QueryMod) userQuery {
	return Users(boil.GetDB(), mods...)
}

// Users retrieves all the records using an executor.
func Users(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	mods = append(mods, qm.From("users"))
	return userQuery{NewQuery(exec, mods...)}
}


// UserFindG retrieves a single record by ID.
func UserFindG(id string, selectCols ...string) (*User, error) {
	return UserFind(boil.GetDB(), id, selectCols...)
}

// UserFindGP retrieves a single record by ID, and panics on error.
func UserFindGP(id string, selectCols ...string) *User {
	retobj, err := UserFind(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// UserFind retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func UserFind(exec boil.Executor, id string, selectCols ...string) (*User, error) {
	user := &User{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(selectCols), ",")
	}
	sql := fmt.Sprintf(
		`select %s from "users" where "id"=$1`, sel,
	)
	q := boil.SQL(sql, id)
	boil.SetExecutor(q, exec)

	err := q.Bind(user)
	if err != nil {
		return nil, fmt.Errorf("models: unable to select from users: %v", err)
	}

	return user, nil
}

// UserFindP retrieves a single record by ID with an executor, and panics on error.
func UserFindP(exec boil.Executor, id string, selectCols ...string) *User {
	retobj, err := UserFind(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *User) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *User) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *User) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are inferred (i.e. name, age)
// - All columns with a default, but non-zero are inferred (i.e. health = 75)
func (o *User) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no users provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO users ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.FirstName, o.LastName, o.Username, o.Phone, o.Email, o.Hash, o.CreatedAt, o.UpdatedAt, o.Banned)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into users: %s", err)
	}

	if err := o.doAfterCreateHooks(); err != nil {
		return err
	}

	return nil
}

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
// the return columns are used to get values that are assigned within the database during the
// insert to keep the struct in sync with what's in the db.
// with a whitelist:
// - the whitelist is used for the insert columns
// - the return columns are the result of (columns with default values - the whitelist)
// without a whitelist:
// - start with columns without a default as these always need to be inserted
// - add all columns that have a default in the database but that are non-zero in the struct
// - the return columns are the result of (columns with default values - the previous set)
func (o *User) generateInsertColumns(whitelist ...string) ([]string, []string) {
	if len(whitelist) > 0 {
		return whitelist, boil.SetComplement(userColumnsWithDefault, whitelist)
	}

	var wl []string

	wl = append(wl, userColumnsWithoutDefault...)

	wl = boil.SetMerge(boil.NonZeroDefaultSet(userColumnsWithDefault, o), wl)
	wl = boil.SortByKeys(userColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(userColumnsWithDefault, wl)

	return wl, rc
}


// UpdateG a single User record. See Update for
// whitelist behavior description.
func (o *User) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single User record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *User) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the User, and panics on error.
// See Update for whitelist behavior description.
func (o *User) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the User.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
func (o *User) Update(exec boil.Executor, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE users SET %s WHERE %s`, strmangle.SetParamNames(wl), strmangle.WhereClause(len(wl)+1, userPrimaryKeyColumns))
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

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q userQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q userQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for users: %s", err)
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UserSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o UserSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o UserSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserSlice) UpdateAll(exec boil.Executor, cols M) error {
	if o == nil {
		return errors.New("models: no User slice provided for update all")
	}

	if len(o) == 0 {
		return nil
	}

	colNames := make([]string, len(cols))
	var args []interface{}

	count := 0
	for name, value := range cols {
		colNames[count] = strmangle.IdentQuote(name)
		args = append(args, value)
		count++
	}

	// Append all of the primary key values for each column
	args = append(args, o.inPrimaryKeyArgs()...)

	sql := fmt.Sprintf(
		`UPDATE users SET (%s) = (%s) WHERE (%s) IN (%s)`,
		strings.Join(colNames, ", "),
		strmangle.Placeholders(len(colNames), 1, 1),
		strings.Join(strmangle.IdentQuoteSlice(userPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(userPrimaryKeyColumns), len(colNames)+1, len(userPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to update all in user slice: %s", err)
	}

	return nil
}

// generateUpdateColumns generates the whitelist columns for an update statement
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func (o *User) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	return boil.SetComplement(userColumns, userPrimaryKeyColumns)
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *User) UpsertG(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *User) UpsertGP(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *User) UpsertP(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *User) Upsert(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no users provided for upsert")
	}

	columns := o.generateUpsertColumns(conflictColumns, updateColumns, whitelist)
	query := o.generateUpsertQuery(update, columns)

	var err error
	if err := o.doBeforeUpsertHooks(); err != nil {
		return err
	}

	if len(columns.returning) != 0 {
		err = exec.QueryRow(query, boil.GetStructValues(o, columns.whitelist...)...).Scan(boil.GetStructPointers(o, columns.returning...)...)
	} else {
		_, err = exec.Exec(query, o.ID, o.FirstName, o.LastName, o.Username, o.Phone, o.Email, o.Hash, o.CreatedAt, o.UpdatedAt, o.Banned)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, columns.whitelist...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to upsert for users: %s", err)
	}

	if err := o.doAfterUpsertHooks(); err != nil {
		return err
	}

	return nil
}

// generateUpsertColumns builds an upsertData object, using generated values when necessary.
func (o *User) generateUpsertColumns(conflict []string, update []string, whitelist []string) upsertData {
	var upsertCols upsertData

	upsertCols.whitelist, upsertCols.returning = o.generateInsertColumns(whitelist...)

	upsertCols.conflict = make([]string, len(conflict))
	upsertCols.update = make([]string, len(update))

	// generates the ON CONFLICT() columns if none are provided
	upsertCols.conflict = o.generateConflictColumns(conflict...)

	// generate the UPDATE SET columns if none are provided
	upsertCols.update = o.generateUpdateColumns(update...)

	return upsertCols
}

// generateConflictColumns returns the user provided columns.
// If no columns are provided, it returns the primary key columns.
func (o *User) generateConflictColumns(columns ...string) []string {
	if len(columns) != 0 {
		return columns
	}

	c := make([]string, len(userPrimaryKeyColumns))
	copy(c, userPrimaryKeyColumns)

	return c
}

// generateUpsertQuery builds a SQL statement string using the upsertData provided.
func (o *User) generateUpsertQuery(update bool, columns upsertData) string {
	var set, query string

	conflict := strmangle.IdentQuoteSlice(columns.conflict)
	whitelist := strmangle.IdentQuoteSlice(columns.whitelist)
	returning := strmangle.IdentQuoteSlice(columns.returning)

	var sets []string
	// Generate the UPDATE SET clause
	for _, v := range columns.update {
		quoted := strmangle.IdentQuote(v)
		sets = append(sets, fmt.Sprintf("%s = EXCLUDED.%s", quoted, quoted))
	}
	set = strings.Join(sets, ", ")

	query = fmt.Sprintf(
		"INSERT INTO users (%s) VALUES (%s) ON CONFLICT",
		strings.Join(whitelist, ", "),
		strmangle.Placeholders(len(whitelist), 1, 1),
	)

	if !update {
		query = query + " DO NOTHING"
	} else {
		query = fmt.Sprintf("%s (%s) DO UPDATE SET %s", query, strings.Join(conflict, ", "), set)
	}

	if len(columns.returning) != 0 {
		query = fmt.Sprintf("%s RETURNING %s", query, strings.Join(returning, ", "))
	}

	return query
}

// DeleteP deletes a single User record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *User) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single User record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *User) DeleteG() error {
	if o == nil {
		return errors.New("models: no User provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single User record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *User) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single User record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *User) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User provided for deletion")
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.From("users"),
		qm.Where(`"id"=$1`, o.ID),
	)

	query := NewQuery(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete from users: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q userQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q userQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no userQuery provided for delete all")
	}

	boil.SetDelete(q.Query)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from users: %s", err)
	}

	return nil
}

// DeleteAll deletes all rows in the slice, and panics on error.
func (o UserSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o UserSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o UserSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`DELETE FROM users WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(userPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(userPrimaryKeyColumns), 1, len(userPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from user slice: %s", err)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *User) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *User) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *User) ReloadG() error {
	if o == nil {
		return errors.New("models: no User provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *User) Reload(exec boil.Executor) error {
	ret, err := UserFind(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *UserSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *UserSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *UserSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty UserSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSlice) ReloadAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User slice provided for reload all")
	}

	if len(*o) == 0 {
		return nil
	}

	users := UserSlice{}
	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`SELECT users.* FROM users WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(userPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(*o)*len(userPrimaryKeyColumns), 1, len(userPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	err := q.Bind(&users)
	if err != nil {
		return fmt.Errorf("models: unable to reload all in UserSlice: %v", err)
	}

	*o = users

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}


// UserExists checks if the User row exists.
func UserExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	row := exec.QueryRow(
		`select exists(select 1 from "users" where "id"=$1 limit 1)`,
		id,
	)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("models: unable to check if users exists: %v", err)
	}

	return exists, nil
}

// UserExistsG checks if the User row exists.
func UserExistsG(id string) (bool, error) {
	return UserExists(boil.GetDB(), id)
}

// UserExistsGP checks if the User row exists. Panics on error.
func UserExistsGP(id string) bool {
	e, err := UserExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// UserExistsP checks if the User row exists. Panics on error.
func UserExistsP(exec boil.Executor, id string) bool {
	e, err := UserExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
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

