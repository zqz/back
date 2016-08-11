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

// Thumbnail is an object representing the database table.
type Thumbnail struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	FileID    string    `boil:"file_id" json:"file_id" toml:"file_id" yaml:"file_id"`
	Size      int       `boil:"size" json:"size" toml:"size" yaml:"size"`
	Hash      string    `boil:"hash" json:"hash" toml:"hash" yaml:"hash"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	//Relationships *ThumbnailRelationships `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// ThumbnailRelationships are where relationships are both cached
// and eagerly loaded.
type ThumbnailRelationships struct {
	File *File
}


var (
	thumbnailColumns                  = []string{"id", "file_id", "size", "hash", "created_at", "updated_at"}
	thumbnailColumnsWithoutDefault    = []string{"file_id", "size", "hash", "created_at", "updated_at"}
	thumbnailColumnsWithDefault       = []string{"id"}
	thumbnailColumnsWithSimpleDefault = []string{}
	thumbnailValidatedColumns         = []string{"id", "file_id"}
	thumbnailUniqueColumns            = []string{}
	thumbnailPrimaryKeyColumns        = []string{"id"}
	thumbnailAutoIncrementColumns     = []string{}
	thumbnailAutoIncPrimaryKey        = ""
)

type (
	ThumbnailSlice []*Thumbnail
	ThumbnailHook  func(*Thumbnail) error

	thumbnailQuery struct {
		*boil.Query
	}
)

var thumbnailBeforeCreateHooks []ThumbnailHook
var thumbnailBeforeUpdateHooks []ThumbnailHook
var thumbnailBeforeUpsertHooks []ThumbnailHook
var thumbnailAfterCreateHooks []ThumbnailHook
var thumbnailAfterUpdateHooks []ThumbnailHook
var thumbnailAfterUpsertHooks []ThumbnailHook

// doBeforeCreateHooks executes all "before create" hooks.
func (o *Thumbnail) doBeforeCreateHooks() (err error) {
	for _, hook := range thumbnailBeforeCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Thumbnail) doBeforeUpdateHooks() (err error) {
	for _, hook := range thumbnailBeforeUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Thumbnail) doBeforeUpsertHooks() (err error) {
	for _, hook := range thumbnailBeforeUpsertHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterCreateHooks executes all "after create" hooks.
func (o *Thumbnail) doAfterCreateHooks() (err error) {
	for _, hook := range thumbnailAfterCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Thumbnail) doAfterUpdateHooks() (err error) {
	for _, hook := range thumbnailAfterUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Thumbnail) doAfterUpsertHooks() (err error) {
	for _, hook := range thumbnailAfterUpsertHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

func ThumbnailAddHook(hookPoint boil.HookPoint, thumbnailHook ThumbnailHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		thumbnailBeforeCreateHooks = append(thumbnailBeforeCreateHooks, thumbnailHook)
	case boil.HookBeforeUpdate:
		thumbnailBeforeUpdateHooks = append(thumbnailBeforeUpdateHooks, thumbnailHook)
	case boil.HookBeforeUpsert:
		thumbnailBeforeUpsertHooks = append(thumbnailBeforeUpsertHooks, thumbnailHook)
	case boil.HookAfterCreate:
		thumbnailAfterCreateHooks = append(thumbnailAfterCreateHooks, thumbnailHook)
	case boil.HookAfterUpdate:
		thumbnailAfterUpdateHooks = append(thumbnailAfterUpdateHooks, thumbnailHook)
	case boil.HookAfterUpsert:
		thumbnailAfterUpsertHooks = append(thumbnailAfterUpsertHooks, thumbnailHook)
	}
}

// OneP returns a single thumbnail record from the query, and panics on error.
func (q thumbnailQuery) OneP() *Thumbnail {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single thumbnail record from the query.
func (q thumbnailQuery) One() (*Thumbnail, error) {
	o := &Thumbnail{}

	boil.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for thumbnails: %s", err)
	}

	return o, nil
}

// AllP returns all Thumbnail records from the query, and panics on error.
func (q thumbnailQuery) AllP() ThumbnailSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Thumbnail records from the query.
func (q thumbnailQuery) All() (ThumbnailSlice, error) {
	var o ThumbnailSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to Thumbnail slice: %s", err)
	}

	return o, nil
}

// CountP returns the count of all Thumbnail records in the query, and panics on error.
func (q thumbnailQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Thumbnail records in the query.
func (q thumbnailQuery) Count() (int64, error) {
	var count int64

	boil.SetCount(q.Query)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("models: failed to count thumbnails rows: %s", err)
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q thumbnailQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q thumbnailQuery) Exists() (bool, error) {
	var count int64

	boil.SetCount(q.Query)
	boil.SetLimit(q.Query, 1)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("models: failed to check if thumbnails exists: %s", err)
	}

	return count > 0, nil
}


// FileG pointed to by the foreign key.
func (t *Thumbnail) FileG(mods ...qm.QueryMod) (*File, error) {
	return t.File(boil.GetDB(), mods...)
}

// FileGP pointed to by the foreign key. Panics on error.
func (t *Thumbnail) FileGP(mods ...qm.QueryMod) *File {
	slice, err := t.File(boil.GetDB(), mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// FileP pointed to by the foreign key with exeuctor. Panics on error.
func (t *Thumbnail) FileP(exec boil.Executor, mods ...qm.QueryMod) *File {
	slice, err := t.File(exec, mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// File pointed to by the foreign key.
func (t *Thumbnail) File(exec boil.Executor, mods ...qm.QueryMod) (*File, error) {
	queryMods := []qm.QueryMod{
		qm.Where("id=$1", t.FileID),
	}

	queryMods = append(queryMods, mods...)

	query := Files(exec, queryMods...)
	boil.SetFrom(query.Query, "files")

	return query.One()
}



// ThumbnailsG retrieves all records.
func ThumbnailsG(mods ...qm.QueryMod) thumbnailQuery {
	return Thumbnails(boil.GetDB(), mods...)
}

// Thumbnails retrieves all the records using an executor.
func Thumbnails(exec boil.Executor, mods ...qm.QueryMod) thumbnailQuery {
	mods = append(mods, qm.From("thumbnails"))
	return thumbnailQuery{NewQuery(exec, mods...)}
}


// ThumbnailFindG retrieves a single record by ID.
func ThumbnailFindG(id string, selectCols ...string) (*Thumbnail, error) {
	return ThumbnailFind(boil.GetDB(), id, selectCols...)
}

// ThumbnailFindGP retrieves a single record by ID, and panics on error.
func ThumbnailFindGP(id string, selectCols ...string) *Thumbnail {
	retobj, err := ThumbnailFind(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// ThumbnailFind retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func ThumbnailFind(exec boil.Executor, id string, selectCols ...string) (*Thumbnail, error) {
	thumbnail := &Thumbnail{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(selectCols), ",")
	}
	sql := fmt.Sprintf(
		`select %s from "thumbnails" where "id"=$1`, sel,
	)
	q := boil.SQL(sql, id)
	boil.SetExecutor(q, exec)

	err := q.Bind(thumbnail)
	if err != nil {
		return nil, fmt.Errorf("models: unable to select from thumbnails: %v", err)
	}

	return thumbnail, nil
}

// ThumbnailFindP retrieves a single record by ID with an executor, and panics on error.
func ThumbnailFindP(exec boil.Executor, id string, selectCols ...string) *Thumbnail {
	retobj, err := ThumbnailFind(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Thumbnail) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Thumbnail) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Thumbnail) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are inferred (i.e. name, age)
// - All columns with a default, but non-zero are inferred (i.e. health = 75)
func (o *Thumbnail) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no thumbnails provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO thumbnails ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.FileID, o.Size, o.Hash, o.CreatedAt, o.UpdatedAt)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into thumbnails: %s", err)
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
func (o *Thumbnail) generateInsertColumns(whitelist ...string) ([]string, []string) {
	if len(whitelist) > 0 {
		return whitelist, boil.SetComplement(thumbnailColumnsWithDefault, whitelist)
	}

	var wl []string

	wl = append(wl, thumbnailColumnsWithoutDefault...)

	wl = boil.SetMerge(boil.NonZeroDefaultSet(thumbnailColumnsWithDefault, o), wl)
	wl = boil.SortByKeys(thumbnailColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(thumbnailColumnsWithDefault, wl)

	return wl, rc
}


// UpdateG a single Thumbnail record. See Update for
// whitelist behavior description.
func (o *Thumbnail) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Thumbnail record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Thumbnail) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Thumbnail, and panics on error.
// See Update for whitelist behavior description.
func (o *Thumbnail) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Thumbnail.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
func (o *Thumbnail) Update(exec boil.Executor, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE thumbnails SET %s WHERE %s`, strmangle.SetParamNames(wl), strmangle.WhereClause(len(wl)+1, thumbnailPrimaryKeyColumns))
		values = boil.GetStructValues(o, wl...)
		values = append(values, o.ID)
		_, err = exec.Exec(query, values...)
	} else {
		return fmt.Errorf("models: unable to update thumbnails, could not build whitelist")
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if err != nil {
		return fmt.Errorf("models: unable to update thumbnails row: %s", err)
	}

	if err := o.doAfterUpdateHooks(); err != nil {
		return err
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q thumbnailQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q thumbnailQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for thumbnails: %s", err)
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ThumbnailSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ThumbnailSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ThumbnailSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ThumbnailSlice) UpdateAll(exec boil.Executor, cols M) error {
	if o == nil {
		return errors.New("models: no Thumbnail slice provided for update all")
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
		`UPDATE thumbnails SET (%s) = (%s) WHERE (%s) IN (%s)`,
		strings.Join(colNames, ", "),
		strmangle.Placeholders(len(colNames), 1, 1),
		strings.Join(strmangle.IdentQuoteSlice(thumbnailPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(thumbnailPrimaryKeyColumns), len(colNames)+1, len(thumbnailPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to update all in thumbnail slice: %s", err)
	}

	return nil
}

// generateUpdateColumns generates the whitelist columns for an update statement
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func (o *Thumbnail) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	return boil.SetComplement(thumbnailColumns, thumbnailPrimaryKeyColumns)
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Thumbnail) UpsertG(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Thumbnail) UpsertGP(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Thumbnail) UpsertP(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Thumbnail) Upsert(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no thumbnails provided for upsert")
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
		_, err = exec.Exec(query, o.ID, o.FileID, o.Size, o.Hash, o.CreatedAt, o.UpdatedAt)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, columns.whitelist...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to upsert for thumbnails: %s", err)
	}

	if err := o.doAfterUpsertHooks(); err != nil {
		return err
	}

	return nil
}

// generateUpsertColumns builds an upsertData object, using generated values when necessary.
func (o *Thumbnail) generateUpsertColumns(conflict []string, update []string, whitelist []string) upsertData {
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
func (o *Thumbnail) generateConflictColumns(columns ...string) []string {
	if len(columns) != 0 {
		return columns
	}

	c := make([]string, len(thumbnailPrimaryKeyColumns))
	copy(c, thumbnailPrimaryKeyColumns)

	return c
}

// generateUpsertQuery builds a SQL statement string using the upsertData provided.
func (o *Thumbnail) generateUpsertQuery(update bool, columns upsertData) string {
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
		"INSERT INTO thumbnails (%s) VALUES (%s) ON CONFLICT",
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

// DeleteP deletes a single Thumbnail record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Thumbnail) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Thumbnail record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Thumbnail) DeleteG() error {
	if o == nil {
		return errors.New("models: no Thumbnail provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Thumbnail record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Thumbnail) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Thumbnail record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Thumbnail) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Thumbnail provided for deletion")
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.From("thumbnails"),
		qm.Where(`"id"=$1`, o.ID),
	)

	query := NewQuery(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete from thumbnails: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q thumbnailQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q thumbnailQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no thumbnailQuery provided for delete all")
	}

	boil.SetDelete(q.Query)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from thumbnails: %s", err)
	}

	return nil
}

// DeleteAll deletes all rows in the slice, and panics on error.
func (o ThumbnailSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ThumbnailSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Thumbnail slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ThumbnailSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ThumbnailSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Thumbnail slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`DELETE FROM thumbnails WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(thumbnailPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(thumbnailPrimaryKeyColumns), 1, len(thumbnailPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from thumbnail slice: %s", err)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Thumbnail) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Thumbnail) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Thumbnail) ReloadG() error {
	if o == nil {
		return errors.New("models: no Thumbnail provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Thumbnail) Reload(exec boil.Executor) error {
	ret, err := ThumbnailFind(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *ThumbnailSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *ThumbnailSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *ThumbnailSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ThumbnailSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ThumbnailSlice) ReloadAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Thumbnail slice provided for reload all")
	}

	if len(*o) == 0 {
		return nil
	}

	thumbnails := ThumbnailSlice{}
	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`SELECT thumbnails.* FROM thumbnails WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(thumbnailPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(*o)*len(thumbnailPrimaryKeyColumns), 1, len(thumbnailPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	err := q.Bind(&thumbnails)
	if err != nil {
		return fmt.Errorf("models: unable to reload all in ThumbnailSlice: %v", err)
	}

	*o = thumbnails

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}


// ThumbnailExists checks if the Thumbnail row exists.
func ThumbnailExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	row := exec.QueryRow(
		`select exists(select 1 from "thumbnails" where "id"=$1 limit 1)`,
		id,
	)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("models: unable to check if thumbnails exists: %v", err)
	}

	return exists, nil
}

// ThumbnailExistsG checks if the Thumbnail row exists.
func ThumbnailExistsG(id string) (bool, error) {
	return ThumbnailExists(boil.GetDB(), id)
}

// ThumbnailExistsGP checks if the Thumbnail row exists. Panics on error.
func ThumbnailExistsGP(id string) bool {
	e, err := ThumbnailExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ThumbnailExistsP checks if the Thumbnail row exists. Panics on error.
func ThumbnailExistsP(exec boil.Executor, id string) bool {
	e, err := ThumbnailExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

func (o Thumbnail) inPrimaryKeyArgs() []interface{} {
	var args []interface{}
	args = append(args, o.ID)
	return args
}

func (o ThumbnailSlice) inPrimaryKeyArgs() []interface{} {
	var args []interface{}

	for i := 0; i < len(o); i++ {
		args = append(args, o[i].ID)
	}

	return args
}

