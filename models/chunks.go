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

// Chunk is an object representing the database table.
type Chunk struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	FileID    string    `boil:"file_id" json:"file_id" toml:"file_id" yaml:"file_id"`
	Size      int       `boil:"size" json:"size" toml:"size" yaml:"size"`
	Hash      string    `boil:"hash" json:"hash" toml:"hash" yaml:"hash"`
	Position  int       `boil:"position" json:"position" toml:"position" yaml:"position"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	//Relationships *ChunkRelationships `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// ChunkRelationships are where relationships are both cached
// and eagerly loaded.
type ChunkRelationships struct {
	File *File
}


var (
	chunkColumns                  = []string{"id", "file_id", "size", "hash", "position", "created_at", "updated_at"}
	chunkColumnsWithoutDefault    = []string{"file_id", "size", "hash", "position", "created_at", "updated_at"}
	chunkColumnsWithDefault       = []string{"id"}
	chunkColumnsWithSimpleDefault = []string{}
	chunkValidatedColumns         = []string{"id", "file_id"}
	chunkUniqueColumns            = []string{}
	chunkPrimaryKeyColumns        = []string{"id"}
	chunkAutoIncrementColumns     = []string{}
	chunkAutoIncPrimaryKey        = ""
)

type (
	ChunkSlice []*Chunk
	ChunkHook  func(*Chunk) error

	chunkQuery struct {
		*boil.Query
	}
)

var chunkBeforeCreateHooks []ChunkHook
var chunkBeforeUpdateHooks []ChunkHook
var chunkBeforeUpsertHooks []ChunkHook
var chunkAfterCreateHooks []ChunkHook
var chunkAfterUpdateHooks []ChunkHook
var chunkAfterUpsertHooks []ChunkHook

// doBeforeCreateHooks executes all "before create" hooks.
func (o *Chunk) doBeforeCreateHooks() (err error) {
	for _, hook := range chunkBeforeCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Chunk) doBeforeUpdateHooks() (err error) {
	for _, hook := range chunkBeforeUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Chunk) doBeforeUpsertHooks() (err error) {
	for _, hook := range chunkBeforeUpsertHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterCreateHooks executes all "after create" hooks.
func (o *Chunk) doAfterCreateHooks() (err error) {
	for _, hook := range chunkAfterCreateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Chunk) doAfterUpdateHooks() (err error) {
	for _, hook := range chunkAfterUpdateHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Chunk) doAfterUpsertHooks() (err error) {
	for _, hook := range chunkAfterUpsertHooks {
		if err := hook(o); err != nil {
			return err
		}
	}

	return nil
}

func ChunkAddHook(hookPoint boil.HookPoint, chunkHook ChunkHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		chunkBeforeCreateHooks = append(chunkBeforeCreateHooks, chunkHook)
	case boil.HookBeforeUpdate:
		chunkBeforeUpdateHooks = append(chunkBeforeUpdateHooks, chunkHook)
	case boil.HookBeforeUpsert:
		chunkBeforeUpsertHooks = append(chunkBeforeUpsertHooks, chunkHook)
	case boil.HookAfterCreate:
		chunkAfterCreateHooks = append(chunkAfterCreateHooks, chunkHook)
	case boil.HookAfterUpdate:
		chunkAfterUpdateHooks = append(chunkAfterUpdateHooks, chunkHook)
	case boil.HookAfterUpsert:
		chunkAfterUpsertHooks = append(chunkAfterUpsertHooks, chunkHook)
	}
}

// OneP returns a single chunk record from the query, and panics on error.
func (q chunkQuery) OneP() *Chunk {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single chunk record from the query.
func (q chunkQuery) One() (*Chunk, error) {
	o := &Chunk{}

	boil.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for chunks: %s", err)
	}

	return o, nil
}

// AllP returns all Chunk records from the query, and panics on error.
func (q chunkQuery) AllP() ChunkSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Chunk records from the query.
func (q chunkQuery) All() (ChunkSlice, error) {
	var o ChunkSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to Chunk slice: %s", err)
	}

	return o, nil
}

// CountP returns the count of all Chunk records in the query, and panics on error.
func (q chunkQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Chunk records in the query.
func (q chunkQuery) Count() (int64, error) {
	var count int64

	boil.SetCount(q.Query)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("models: failed to count chunks rows: %s", err)
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q chunkQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q chunkQuery) Exists() (bool, error) {
	var count int64

	boil.SetCount(q.Query)
	boil.SetLimit(q.Query, 1)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("models: failed to check if chunks exists: %s", err)
	}

	return count > 0, nil
}


// FileG pointed to by the foreign key.
func (c *Chunk) FileG(mods ...qm.QueryMod) (*File, error) {
	return c.File(boil.GetDB(), mods...)
}

// FileGP pointed to by the foreign key. Panics on error.
func (c *Chunk) FileGP(mods ...qm.QueryMod) *File {
	slice, err := c.File(boil.GetDB(), mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// FileP pointed to by the foreign key with exeuctor. Panics on error.
func (c *Chunk) FileP(exec boil.Executor, mods ...qm.QueryMod) *File {
	slice, err := c.File(exec, mods...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return slice
}

// File pointed to by the foreign key.
func (c *Chunk) File(exec boil.Executor, mods ...qm.QueryMod) (*File, error) {
	queryMods := []qm.QueryMod{
		qm.Where("id=$1", c.FileID),
	}

	queryMods = append(queryMods, mods...)

	query := Files(exec, queryMods...)
	boil.SetFrom(query.Query, "files")

	return query.One()
}



// ChunksG retrieves all records.
func ChunksG(mods ...qm.QueryMod) chunkQuery {
	return Chunks(boil.GetDB(), mods...)
}

// Chunks retrieves all the records using an executor.
func Chunks(exec boil.Executor, mods ...qm.QueryMod) chunkQuery {
	mods = append(mods, qm.From("chunks"))
	return chunkQuery{NewQuery(exec, mods...)}
}


// ChunkFindG retrieves a single record by ID.
func ChunkFindG(id string, selectCols ...string) (*Chunk, error) {
	return ChunkFind(boil.GetDB(), id, selectCols...)
}

// ChunkFindGP retrieves a single record by ID, and panics on error.
func ChunkFindGP(id string, selectCols ...string) *Chunk {
	retobj, err := ChunkFind(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// ChunkFind retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func ChunkFind(exec boil.Executor, id string, selectCols ...string) (*Chunk, error) {
	chunk := &Chunk{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(selectCols), ",")
	}
	sql := fmt.Sprintf(
		`select %s from "chunks" where "id"=$1`, sel,
	)
	q := boil.SQL(sql, id)
	boil.SetExecutor(q, exec)

	err := q.Bind(chunk)
	if err != nil {
		return nil, fmt.Errorf("models: unable to select from chunks: %v", err)
	}

	return chunk, nil
}

// ChunkFindP retrieves a single record by ID with an executor, and panics on error.
func ChunkFindP(exec boil.Executor, id string, selectCols ...string) *Chunk {
	retobj, err := ChunkFind(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Chunk) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Chunk) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Chunk) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are inferred (i.e. name, age)
// - All columns with a default, but non-zero are inferred (i.e. health = 75)
func (o *Chunk) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no chunks provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO chunks ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.FileID, o.Size, o.Hash, o.Position, o.CreatedAt, o.UpdatedAt)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into chunks: %s", err)
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
func (o *Chunk) generateInsertColumns(whitelist ...string) ([]string, []string) {
	if len(whitelist) > 0 {
		return whitelist, boil.SetComplement(chunkColumnsWithDefault, whitelist)
	}

	var wl []string

	wl = append(wl, chunkColumnsWithoutDefault...)

	wl = boil.SetMerge(boil.NonZeroDefaultSet(chunkColumnsWithDefault, o), wl)
	wl = boil.SortByKeys(chunkColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(chunkColumnsWithDefault, wl)

	return wl, rc
}


// UpdateG a single Chunk record. See Update for
// whitelist behavior description.
func (o *Chunk) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Chunk record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Chunk) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Chunk, and panics on error.
// See Update for whitelist behavior description.
func (o *Chunk) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Chunk.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
func (o *Chunk) Update(exec boil.Executor, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE chunks SET %s WHERE %s`, strmangle.SetParamNames(wl), strmangle.WhereClause(len(wl)+1, chunkPrimaryKeyColumns))
		values = boil.GetStructValues(o, wl...)
		values = append(values, o.ID)
		_, err = exec.Exec(query, values...)
	} else {
		return fmt.Errorf("models: unable to update chunks, could not build whitelist")
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if err != nil {
		return fmt.Errorf("models: unable to update chunks row: %s", err)
	}

	if err := o.doAfterUpdateHooks(); err != nil {
		return err
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q chunkQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q chunkQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for chunks: %s", err)
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ChunkSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ChunkSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ChunkSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ChunkSlice) UpdateAll(exec boil.Executor, cols M) error {
	if o == nil {
		return errors.New("models: no Chunk slice provided for update all")
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
		`UPDATE chunks SET (%s) = (%s) WHERE (%s) IN (%s)`,
		strings.Join(colNames, ", "),
		strmangle.Placeholders(len(colNames), 1, 1),
		strings.Join(strmangle.IdentQuoteSlice(chunkPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(chunkPrimaryKeyColumns), len(colNames)+1, len(chunkPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to update all in chunk slice: %s", err)
	}

	return nil
}

// generateUpdateColumns generates the whitelist columns for an update statement
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func (o *Chunk) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	return boil.SetComplement(chunkColumns, chunkPrimaryKeyColumns)
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Chunk) UpsertG(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Chunk) UpsertGP(update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Chunk) UpsertP(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, update, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Chunk) Upsert(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no chunks provided for upsert")
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
		_, err = exec.Exec(query, o.ID, o.FileID, o.Size, o.Hash, o.Position, o.CreatedAt, o.UpdatedAt)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, columns.whitelist...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to upsert for chunks: %s", err)
	}

	if err := o.doAfterUpsertHooks(); err != nil {
		return err
	}

	return nil
}

// generateUpsertColumns builds an upsertData object, using generated values when necessary.
func (o *Chunk) generateUpsertColumns(conflict []string, update []string, whitelist []string) upsertData {
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
func (o *Chunk) generateConflictColumns(columns ...string) []string {
	if len(columns) != 0 {
		return columns
	}

	c := make([]string, len(chunkPrimaryKeyColumns))
	copy(c, chunkPrimaryKeyColumns)

	return c
}

// generateUpsertQuery builds a SQL statement string using the upsertData provided.
func (o *Chunk) generateUpsertQuery(update bool, columns upsertData) string {
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
		"INSERT INTO chunks (%s) VALUES (%s) ON CONFLICT",
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

// DeleteP deletes a single Chunk record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Chunk) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Chunk record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Chunk) DeleteG() error {
	if o == nil {
		return errors.New("models: no Chunk provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Chunk record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Chunk) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Chunk record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Chunk) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Chunk provided for deletion")
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.From("chunks"),
		qm.Where(`"id"=$1`, o.ID),
	)

	query := NewQuery(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete from chunks: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q chunkQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q chunkQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no chunkQuery provided for delete all")
	}

	boil.SetDelete(q.Query)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from chunks: %s", err)
	}

	return nil
}

// DeleteAll deletes all rows in the slice, and panics on error.
func (o ChunkSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ChunkSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Chunk slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ChunkSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ChunkSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Chunk slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`DELETE FROM chunks WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(chunkPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(o)*len(chunkPrimaryKeyColumns), 1, len(chunkPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	_, err := boil.ExecQuery(q)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from chunk slice: %s", err)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Chunk) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Chunk) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Chunk) ReloadG() error {
	if o == nil {
		return errors.New("models: no Chunk provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Chunk) Reload(exec boil.Executor) error {
	ret, err := ChunkFind(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *ChunkSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *ChunkSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

func (o *ChunkSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ChunkSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ChunkSlice) ReloadAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Chunk slice provided for reload all")
	}

	if len(*o) == 0 {
		return nil
	}

	chunks := ChunkSlice{}
	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		`SELECT chunks.* FROM chunks WHERE (%s) IN (%s)`,
		strings.Join(strmangle.IdentQuoteSlice(chunkPrimaryKeyColumns), ","),
		strmangle.Placeholders(len(*o)*len(chunkPrimaryKeyColumns), 1, len(chunkPrimaryKeyColumns)),
	)

	q := boil.SQL(sql, args...)
	boil.SetExecutor(q, exec)

	err := q.Bind(&chunks)
	if err != nil {
		return fmt.Errorf("models: unable to reload all in ChunkSlice: %v", err)
	}

	*o = chunks

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}


// ChunkExists checks if the Chunk row exists.
func ChunkExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	row := exec.QueryRow(
		`select exists(select 1 from "chunks" where "id"=$1 limit 1)`,
		id,
	)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("models: unable to check if chunks exists: %v", err)
	}

	return exists, nil
}

// ChunkExistsG checks if the Chunk row exists.
func ChunkExistsG(id string) (bool, error) {
	return ChunkExists(boil.GetDB(), id)
}

// ChunkExistsGP checks if the Chunk row exists. Panics on error.
func ChunkExistsGP(id string) bool {
	e, err := ChunkExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ChunkExistsP checks if the Chunk row exists. Panics on error.
func ChunkExistsP(exec boil.Executor, id string) bool {
	e, err := ChunkExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

func (o Chunk) inPrimaryKeyArgs() []interface{} {
	var args []interface{}
	args = append(args, o.ID)
	return args
}

func (o ChunkSlice) inPrimaryKeyArgs() []interface{} {
	var args []interface{}

	for i := 0; i < len(o); i++ {
		args = append(args, o[i].ID)
	}

	return args
}

