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

// File is an object representing the database table.
type File struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Size      int       `boil:"size" json:"size" toml:"size" yaml:"size"`
	NumChunks int       `boil:"num_chunks" json:"num_chunks" toml:"num_chunks" yaml:"num_chunks"`
	State     int       `boil:"state" json:"state" toml:"state" yaml:"state"`
	Name      string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	Hash      string    `boil:"hash" json:"hash" toml:"hash" yaml:"hash"`
	Type      string    `boil:"type" json:"type" toml:"type" yaml:"type"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	Slug      string    `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`

	//Relationships *FileRelationships `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// FileRelationships are where relationships are both cached
// and eagerly loaded.
type FileRelationships struct {
	Chunks     ChunkSlice
	Thumbnails ThumbnailSlice
}


var (
	fileColumns                  = []string{"id", "size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at", "slug"}
	fileColumnsWithoutDefault    = []string{"size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at"}
	fileColumnsWithDefault       = []string{"id", "slug"}
	fileColumnsWithSimpleDefault = []string{}
	fileValidatedColumns         = []string{"id"}
	fileUniqueColumns            = []string{}
	filePrimaryKeyColumns        = []string{"id"}
	fileAutoIncrementColumns     = []string{}
	fileAutoIncPrimaryKey        = ""
)

type (
	FileSlice []*File
	FileHook  func(*File) error

	fileQuery struct {
		*boil.Query
	}
)

var fileBeforeCreateHooks []FileHook
var fileBeforeUpdateHooks []FileHook
var fileBeforeUpsertHooks []FileHook
var fileAfterCreateHooks []FileHook
var fileAfterUpdateHooks []FileHook
var fileAfterUpsertHooks []FileHook

// doBeforeCreateHooks executes all "before create" hooks.
func (o *File) doBeforeCreateHooks() (err error) {
	for _, hook := range fileBeforeCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *File) doBeforeUpdateHooks() (err error) {
	for _, hook := range fileBeforeUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *File) doBeforeUpsertHooks() (err error) {
	for _, hook := range fileBeforeUpsertHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterCreateHooks executes all "after create" hooks.
func (o *File) doAfterCreateHooks() (err error) {
	for _, hook := range fileAfterCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *File) doAfterUpdateHooks() (err error) {
	for _, hook := range fileAfterUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *File) doAfterUpsertHooks() (err error) {
	for _, hook := range fileAfterUpsertHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

func FileAddHook(hookPoint boil.HookPoint, fileHook FileHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		fileBeforeCreateHooks = append(fileBeforeCreateHooks, fileHook)
	case boil.HookBeforeUpdate:
		fileBeforeUpdateHooks = append(fileBeforeUpdateHooks, fileHook)
	case boil.HookBeforeUpsert:
		fileBeforeUpsertHooks = append(fileBeforeUpsertHooks, fileHook)
	case boil.HookAfterCreate:
		fileAfterCreateHooks = append(fileAfterCreateHooks, fileHook)
	case boil.HookAfterUpdate:
		fileAfterUpdateHooks = append(fileAfterUpdateHooks, fileHook)
	case boil.HookAfterUpsert:
		fileAfterUpsertHooks = append(fileAfterUpsertHooks, fileHook)
	}
}

// OneP returns a single file record from the query, and panics on error.
func (q fileQuery) OneP() *File {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single file record from the query.
func (q fileQuery) One() (*File, error) {
	o := &File{}

	boil.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for files: %s", err)
	}

	return o, nil
}

// AllP returns all File records from the query, and panics on error.
func (q fileQuery) AllP() FileSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all File records from the query.
func (q fileQuery) All() (FileSlice, error) {
	var o FileSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to File slice: %s", err)
	}

	return o, nil
}

// CountP returns the count of all File records in the query, and panics on error.
func (q fileQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all File records in the query.
func (q fileQuery) Count() (int64, error) {
	var count int64

	boil.SetCount(q.Query)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("models: failed to count files rows: %s", err)
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q fileQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q fileQuery) Exists() (bool, error) {
	var count int64

	boil.SetCount(q.Query)
	boil.SetLimit(q.Query, 1)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("models: failed to check if files exists: %s", err)
	}

	return count > 0, nil
}


// ChunksG retrieves all the file's chunks.
func (f *File) ChunksG(mods ...qm.QueryMod) (ChunkSlice, error) {
	return f.Chunks(boil.GetDB(), mods...)
}

// ChunksGP panics on error. Retrieves all the file's chunks.
func (f *File) ChunksGP(mods ...qm.QueryMod) ChunkSlice {
	slice, err := f.Chunks(boil.GetDB(), mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// ChunksP panics on error. Retrieves all the file's chunks with an executor.
func (f *File) ChunksP(exec boil.Executor, mods ...qm.QueryMod) ChunkSlice {
	slice, err := f.Chunks(exec, mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// Chunks retrieves all the file's chunks with an executor.
func (f *File) Chunks(exec boil.Executor, mods ...qm.QueryMod) (ChunkSlice, error) {
	queryMods := []qm.QueryMod{
		qm.Select(`"a".*`),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where(`"a"."file_id"=$1`, f.ID),
	)

	query := Chunks(exec, queryMods...)
	boil.SetFrom(query.Query, `"chunks" as "a"`)
	return query.All()
}

// ThumbnailsG retrieves all the file's thumbnails.
func (f *File) ThumbnailsG(mods ...qm.QueryMod) (ThumbnailSlice, error) {
	return f.Thumbnails(boil.GetDB(), mods...)
}

// ThumbnailsGP panics on error. Retrieves all the file's thumbnails.
func (f *File) ThumbnailsGP(mods ...qm.QueryMod) ThumbnailSlice {
	slice, err := f.Thumbnails(boil.GetDB(), mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// ThumbnailsP panics on error. Retrieves all the file's thumbnails with an executor.
func (f *File) ThumbnailsP(exec boil.Executor, mods ...qm.QueryMod) ThumbnailSlice {
	slice, err := f.Thumbnails(exec, mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// Thumbnails retrieves all the file's thumbnails with an executor.
func (f *File) Thumbnails(exec boil.Executor, mods ...qm.QueryMod) (ThumbnailSlice, error) {
	queryMods := []qm.QueryMod{
		qm.Select(`"a".*`),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where(`"a"."file_id"=$1`, f.ID),
	)

	query := Thumbnails(exec, queryMods...)
	boil.SetFrom(query.Query, `"thumbnails" as "a"`)
	return query.All()
}


// FilesG retrieves all records.
func FilesG(mods ...qm.QueryMod) fileQuery {
	return Files(boil.GetDB(), mods...)
}

// Files retrieves all the records using an executor.
func Files(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	mods = append(mods, qm.From("files"))
	return fileQuery{NewQuery(exec, mods...)}
}


// FileFindG retrieves a single record by ID.
func FileFindG(id string, selectCols ...string) (*File, error) {
	return FileFind(boil.GetDB(), id, selectCols...)
}

// FileFindGP retrieves a single record by ID, and panics on error.
func FileFindGP(id string, selectCols ...string) *File {
	retobj, err := FileFind(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FileFind retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FileFind(exec boil.Executor, id string, selectCols ...string) (*File, error) {
	file := &File{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(selectCols), ",")
	}
	sql := fmt.Sprintf(
		`select %s from "files" where "id"=$1`, sel,
	)
	q := boil.SQL(sql, id)
	boil.SetExecutor(q, exec)

	err := q.Bind(file)
	if err != nil {
		return nil, fmt.Errorf("models: unable to select from files: %v", err)
	}

	return file, nil
}

// FileFindP retrieves a single record by ID with an executor, and panics on error.
func FileFindP(exec boil.Executor, id string, selectCols ...string) *File {
	retobj, err := FileFind(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *File) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *File) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *File) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are inferred (i.e. name, age)
// - All columns with a default, but non-zero are inferred (i.e. health = 75)
func (o *File) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no files provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO files ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.Size, o.NumChunks, o.State, o.Name, o.Hash, o.Type, o.CreatedAt, o.UpdatedAt, o.Slug)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into files: %s", err)
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
func (o *File) generateInsertColumns(whitelist ...string) ([]string, []string) {
	if len(whitelist) > 0 {
		return whitelist, boil.SetComplement(fileColumnsWithDefault, whitelist)
	}

	var wl []string

	wl = append(wl, fileColumnsWithoutDefault...)

	wl = boil.SetMerge(boil.NonZeroDefaultSet(fileColumnsWithDefault, o), wl)
	wl = boil.SortByKeys(fileColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(fileColumnsWithDefault, wl)

	return wl, rc
}


// UpdateG a single File record. See Update for
// whitelist behavior description.
func (o *File) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single File record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *File) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the File, and panics on error.
// See Update for whitelist behavior description.
func (o *File) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the File.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
func (o *File) Update(exec boil.Executor, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE files SET %s WHERE %s`, strmangle.SetParamNames(wl), strmangle.WhereClause(len(wl)+1, filePrimaryKeyColumns))
		values = boil.GetStructValues(o, wl...)
		values = append(values, o.ID)
		_, err = exec.Exec(query, values...)
	} else {
		return fmt.Errorf("models: unable to update files, could not build whitelist")
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if err != nil {
		return fmt.Errorf("models: unable to update files row: %s", err)
	}

	if err := o.doAfterUpdateHooks(); err != nil {
		return err
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q fileQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q fileQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for files: %s", err)
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o FileSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o FileSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o FileSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FileSlice) UpdateAll(exec boil.Executor, cols M) error {
	if o == nil {
		return errors.New("models: no File slice provided for update all")
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
		`UPDATE files SET (%s) = (%s) WHERE (%s) IN (%s)`,
		strings.Join(colNames, ", "),
		strmangle.Placeholders(len(colNames), 1, 1),
		strings.Join(strmangle.IdentQuoteSlice(filePrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(filePrimaryKeyColumns), len(colNames)+1, len(filePrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to update all in file slice: %s", err)
	}

	return nil
}

// generateUpdateColumns generates the whitelist columns for an update statement
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func (o *File) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	return boil.SetComplement(fileColumns, filePrimaryKeyColumns)
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *File) UpsertG(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *File) UpsertGP(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *File) UpsertP(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *File) Upsert(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no files provided for upsert")
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
		_, err = exec.Exec(query, o.ID, o.Size, o.NumChunks, o.State, o.Name, o.Hash, o.Type, o.CreatedAt, o.UpdatedAt, o.Slug)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, columns.whitelist...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to upsert for files: %s", err)
	}

	if err := o.doAfterUpsertHooks(); err != nil {
		return err
	}

	return nil
}

// generateUpsertColumns builds an upsertData object, using generated values when necessary.
func (o *File) generateUpsertColumns(conflict []string, update []string, whitelist []string) upsertData {
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
func (o *File) generateConflictColumns(columns ...string) []string {
	if len(columns) != 0 {
		return columns
	}

	c := make([]string, len(filePrimaryKeyColumns))
	copy(c, filePrimaryKeyColumns)

	return c
}

// generateUpsertQuery builds a SQL statement string using the upsertData provided.
func (o *File) generateUpsertQuery(update bool, columns upsertData) string {
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
		"INSERT INTO files (%s) VALUES (%s) ON CONFLICT",
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

// DeleteP deletes a single File record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *File) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single File record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *File) DeleteG() error {
	if o == nil {
		return errors.New("models: no File provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single File record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *File) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single File record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *File) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no File provided for deletion")
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.From("files"),
		qm.Where(`"id"=$1`, o.ID),
	)

	query := NewQuery(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete from files: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q fileQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q fileQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no fileQuery provided for delete all")
	}

	boil.SetDelete(q.Query)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from files: %s", err)
	}

	return nil
}

// DeleteAll deletes all rows in the slice, and panics on error.
func (o FileSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o FileSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no File slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o FileSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FileSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no File slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`DELETE FROM files WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(filePrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(filePrimaryKeyColumns), 1, len(filePrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from file slice: %s", err)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *File) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *File) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *File) ReloadG() error {
	if o == nil {
		return errors.New("models: no File provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *File) Reload(exec boil.Executor) error {
	ret, err := FileFind(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *FileSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *FileSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *FileSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty FileSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileSlice) ReloadAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no File slice provided for reload all")
	}

	if len(*o) == 0 {
		return nil
	}

	files := FileSlice{}
	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`SELECT files.* FROM files WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(filePrimaryKeyColumns), ","),
		strmangle.Placeholders(len(*o)*len(filePrimaryKeyColumns), 1, len(filePrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	err := q.Bind(&files)
	if err != nil {
		return fmt.Errorf("models: unable to reload all in FileSlice: %v", err)
	}

	*o = files

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}


// FileExists checks if the File row exists.
func FileExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	row := exec.QueryRow(
		`select exists(select 1 from "files" where "id"=$1 limit 1)`,
		id,
	)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("models: unable to check if files exists: %v", err)
	}

	return exists, nil
}

// FileExistsG checks if the File row exists.
func FileExistsG(id string) (bool, error) {
	return FileExists(boil.GetDB(), id)
}

// FileExistsGP checks if the File row exists. Panics on error.
func FileExistsGP(id string) bool {
	e, err := FileExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// FileExistsP checks if the File row exists. Panics on error.
func FileExistsP(exec boil.Executor, id string) bool {
	e, err := FileExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

func (o File) inPrimaryKeyArgs() []interface{} {
	var args []interface{}
	args = append(args, o.ID)
	return args
}

func (o FileSlice) inPrimaryKeyArgs() []interface{} {
	var args []interface{}

	for i := 0; i < len(o); i++ {
		args = append(args, o[i].ID)
	}

	return args
}

