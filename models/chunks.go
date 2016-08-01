package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
)

// Chunk is an object representing the database table.
type Chunk struct {
	ID        string    `db:"chunk_id" json:"id"`
	FileID    string    `db:"chunk_file_id" json:"file_id"`
	Size      int       `db:"chunk_size" json:"size"`
	Hash      string    `db:"chunk_hash" json:"hash"`
	Position  int       `db:"chunk_position" json:"position"`
	CreatedAt time.Time `db:"chunk_created_at" json:"created_at"`
	UpdatedAt time.Time `db:"chunk_updated_at" json:"updated_at"`
}

var (
	chunkColumns                  = []string{"id", "file_id", "size", "hash", "position", "created_at", "updated_at"}
	chunkColumnsWithoutDefault    = []string{"file_id", "size", "hash", "position", "created_at", "updated_at"}
	chunkColumnsWithDefault       = []string{"id"}
	chunkColumnsWithSimpleDefault = []string{}
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
var chunkAfterCreateHooks []ChunkHook
var chunkAfterUpdateHooks []ChunkHook

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

func ChunkAddHook(hookPoint boil.HookPoint, chunkHook ChunkHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		chunkBeforeCreateHooks = append(chunkBeforeCreateHooks, chunkHook)
	case boil.HookBeforeUpdate:
		chunkBeforeUpdateHooks = append(chunkBeforeUpdateHooks, chunkHook)
	case boil.HookAfterCreate:
		chunkAfterCreateHooks = append(chunkAfterCreateHooks, chunkHook)
	case boil.HookAfterUpdate:
		chunkAfterUpdateHooks = append(chunkAfterUpdateHooks, chunkHook)
	}
}

// One returns a single chunk record from the query.
func (q chunkQuery) One() (*Chunk, error) {
	o := &Chunk{}

	boil.SetLimit(q.Query, 1)

	res := boil.ExecQueryOne(q.Query)
	err := boil.BindOne(res, boil.Select(q.Query), o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for chunks: %s", err)
	}

	return o, nil
}

// OneP returns a single chunk record from the query, and panics on error.
func (q chunkQuery) OneP() *Chunk {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Chunk records from the query.
func (q chunkQuery) All() (ChunkSlice, error) {
	var o ChunkSlice

	res, err := boil.ExecQueryAll(q.Query)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute an all query for chunks: %s", err)
	}
	defer res.Close()

	err = boil.BindAll(res, boil.Select(q.Query), &o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to Chunk slice: %s", err)
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

// CountP returns the count of all Chunk records in the query, and panics on error.
func (q chunkQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}


// FileG pointed to by the foreign key.
func (c *Chunk) FileG(selectCols ...string) (*File, error) {
	return c.File(boil.GetDB(), selectCols...)
}

// FileGP pointed to by the foreign key. Panics on error.
func (c *Chunk) FileGP(selectCols ...string) *File {
	o, err := c.File(boil.GetDB(), selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// FileP pointed to by the foreign key with exeuctor. Panics on error.
func (c *Chunk) FileP(exec boil.Executor, selectCols ...string) *File {
	o, err := c.File(exec, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// File pointed to by the foreign key.
func (c *Chunk) File(exec boil.Executor, selectCols ...string) (*File, error) {
	file := &File{}

	selectColumns := `*`
	if len(selectCols) != 0 {
		selectColumns = fmt.Sprintf(`"%s"`, strings.Join(selectCols, `","`))
	}

	query := fmt.Sprintf(`select %s from files where "id" = $1`, selectColumns)
	err := exec.QueryRow(query, c.FileID).Scan(boil.GetStructPointers(file, selectCols...)...)
	if err != nil {
		return nil, fmt.Errorf(`models: unable to select from files: %v`, err)
	}

	return file, nil
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
	o, err := ChunkFind(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// ChunkFind retrieves a single record by ID with an executor.
func ChunkFind(exec boil.Executor, id string, selectCols ...string) (*Chunk, error) {
	chunk := &Chunk{}

	mods := []qm.QueryMod{
		qm.Select(selectCols...),
		qm.From("chunks"),
		qm.Where(`"id"=$1`, id),
	}

	q := NewQuery(exec, mods...)

	err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers(chunk, selectCols...)...)

	if err != nil {
		return nil, fmt.Errorf("models: unable to select from chunks: %v", err)
	}

	return chunk, nil
}

// ChunkFindP retrieves a single record by ID with an executor, and panics on error.
func ChunkFindP(exec boil.Executor, id string, selectCols ...string) *Chunk {
	o, err := ChunkFind(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// InsertG a single record.
func (o *Chunk) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error.
func (o *Chunk) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
func (o *Chunk) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no chunks provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO chunks ("%s") VALUES (%s)`, strings.Join(wl, `","`), boil.GenerateParamFlags(len(wl), 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.FileID, o.Size, o.Hash, o.Position, o.CreatedAt, o.UpdatedAt)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into chunks: %s", err)
	}

	if err := o.doAfterCreateHooks(); err != nil {
		return err
	}

	return nil
}

// InsertP a single record using an executor, and panics on error.
func (o *Chunk) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
func (o *Chunk) generateInsertColumns(whitelist ...string) ([]string, []string) {
	var wl []string

	wl = append(wl, whitelist...)
	if len(whitelist) == 0 {
		wl = append(wl, chunkColumnsWithoutDefault...)
	}

	wl = append(boil.NonZeroDefaultSet(chunkColumnsWithDefault, o), wl...)
	wl = boil.SortByKeys(chunkColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(chunkColumnsWithDefault, wl)

	return wl, rc
}


// UpdateG a single Chunk record.
// UpdateG takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
func (o *Chunk) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Chunk record.
// UpdateGP takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
// Panics on error.
func (o *Chunk) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Chunk.
func (o *Chunk) Update(exec boil.Executor, whitelist ...string) error {
	return o.UpdateAt(exec, o.ID, whitelist...)
}

// UpdateP uses an executor to update the Chunk, and panics on error.
func (o *Chunk) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.UpdateAt(exec, o.ID, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAtG updates the Chunk using the primary key to find the row to update.
func (o *Chunk) UpdateAtG(id string, whitelist ...string) error {
	return o.UpdateAt(boil.GetDB(), id, whitelist...)
}

// UpdateAtGP updates the Chunk using the primary key to find the row to update. Panics on error.
func (o *Chunk) UpdateAtGP(id string, whitelist ...string) {
	if err := o.UpdateAt(boil.GetDB(), id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAt uses an executor to update the Chunk using the primary key to find the row to update.
func (o *Chunk) UpdateAt(exec boil.Executor, id string, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE chunks SET %s WHERE %s`, boil.SetParamNames(wl), boil.WherePrimaryKey(len(wl)+1, "id"))
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

// UpdateAtP uses an executor to update the Chunk using the primary key to find the row to update.
// Panics on error.
func (o *Chunk) UpdateAtP(exec boil.Executor, id string, whitelist ...string) {
	if err := o.UpdateAt(exec, id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with matching column names.
func (q chunkQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for chunks: %s", err)
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q chunkQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateUpdateColumns generates the whitelist columns for an update statement
func (o *Chunk) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	var wl []string
	cols := chunkColumnsWithoutDefault
	cols = append(boil.NonZeroDefaultSet(chunkColumnsWithDefault, o), cols...)
	// Subtract primary keys and autoincrement columns
	cols = boil.SetComplement(cols, chunkPrimaryKeyColumns)
	cols = boil.SetComplement(cols, chunkAutoIncrementColumns)

	wl = make([]string, len(cols))
	copy(wl, cols)

	return wl
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

// DeleteP deletes a single Chunk record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Chunk) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows.
func (o chunkQuery) DeleteAll() error {
	if o.Query == nil {
		return errors.New("models: no chunkQuery provided for delete all")
	}

	boil.SetDelete(o.Query)

	_, err := boil.ExecQuery(o.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from chunks: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (o chunkQuery) DeleteAllP() {
	if err := o.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
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

// DeleteAll deletes all rows in the slice with an executor.
func (o ChunkSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Chunk slice provided for delete all")
	}

	var mods []qm.QueryMod

	args := o.inPrimaryKeyArgs()
	in := boil.WherePrimaryKeyIn(len(o), "id")

	mods = append(mods,
		qm.From("chunks"),
		qm.Where(in, args...),
	)

	query := NewQuery(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from chunk slice: %s", err)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// DeleteAllP deletes all rows in the slice with an executor, and panics on error.
func (o ChunkSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
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

