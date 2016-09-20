package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
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

	R *chunkR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L chunkL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// chunkR is where relationships are stored.
type chunkR struct {
	File *File
}

// chunkL is where Load methods for each relationship are stored.
type chunkL struct{}

var (
	chunkColumns               = []string{"id", "file_id", "size", "hash", "position", "created_at", "updated_at"}
	chunkColumnsWithoutDefault = []string{"file_id", "size", "hash", "position", "created_at", "updated_at"}
	chunkColumnsWithDefault    = []string{"id"}
	chunkPrimaryKeyColumns     = []string{"id"}
)

type (
	// ChunkSlice is an alias for a slice of pointers to Chunk.
	// This should generally be used opposed to []Chunk.
	ChunkSlice []*Chunk
	// ChunkHook is the signature for custom Chunk hook methods
	ChunkHook func(boil.Executor, *Chunk) error

	chunkQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	chunkType                 = reflect.TypeOf(&Chunk{})
	chunkMapping              = queries.MakeStructMapping(chunkType)
	chunkPrimaryKeyMapping, _ = queries.BindMapping(chunkType, chunkMapping, chunkPrimaryKeyColumns)
	chunkInsertCacheMut       sync.RWMutex
	chunkInsertCache          = make(map[string]insertCache)
	chunkUpdateCacheMut       sync.RWMutex
	chunkUpdateCache          = make(map[string]updateCache)
	chunkUpsertCacheMut       sync.RWMutex
	chunkUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

var chunkBeforeInsertHooks []ChunkHook
var chunkBeforeUpdateHooks []ChunkHook
var chunkBeforeDeleteHooks []ChunkHook
var chunkBeforeUpsertHooks []ChunkHook

var chunkAfterInsertHooks []ChunkHook
var chunkAfterSelectHooks []ChunkHook
var chunkAfterUpdateHooks []ChunkHook
var chunkAfterDeleteHooks []ChunkHook
var chunkAfterUpsertHooks []ChunkHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Chunk) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Chunk) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Chunk) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Chunk) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Chunk) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Chunk) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Chunk) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Chunk) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Chunk) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range chunkAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddChunkHook registers your hook function for all future operations.
func AddChunkHook(hookPoint boil.HookPoint, chunkHook ChunkHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		chunkBeforeInsertHooks = append(chunkBeforeInsertHooks, chunkHook)
	case boil.BeforeUpdateHook:
		chunkBeforeUpdateHooks = append(chunkBeforeUpdateHooks, chunkHook)
	case boil.BeforeDeleteHook:
		chunkBeforeDeleteHooks = append(chunkBeforeDeleteHooks, chunkHook)
	case boil.BeforeUpsertHook:
		chunkBeforeUpsertHooks = append(chunkBeforeUpsertHooks, chunkHook)
	case boil.AfterInsertHook:
		chunkAfterInsertHooks = append(chunkAfterInsertHooks, chunkHook)
	case boil.AfterSelectHook:
		chunkAfterSelectHooks = append(chunkAfterSelectHooks, chunkHook)
	case boil.AfterUpdateHook:
		chunkAfterUpdateHooks = append(chunkAfterUpdateHooks, chunkHook)
	case boil.AfterDeleteHook:
		chunkAfterDeleteHooks = append(chunkAfterDeleteHooks, chunkHook)
	case boil.AfterUpsertHook:
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

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for chunks")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
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
		return nil, errors.Wrap(err, "models: failed to assign all query results to Chunk slice")
	}

	if len(chunkAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
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

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count chunks rows")
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

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if chunks exists")
	}

	return count > 0, nil
}

// FileG pointed to by the foreign key.
func (c *Chunk) FileG(mods ...qm.QueryMod) fileQuery {
	return c.File(boil.GetDB(), mods...)
}

// File pointed to by the foreign key.
func (c *Chunk) File(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=$1", c.FileID),
	}

	queryMods = append(queryMods, mods...)

	query := Files(exec, queryMods...)
	queries.SetFrom(query.Query, "\"files\"")

	return query
}



// LoadFile allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (chunkL) LoadFile(e boil.Executor, singular bool, maybeChunk interface{}) error {
	var slice []*Chunk
	var object *Chunk

	count := 1
	if singular {
		object = maybeChunk.(*Chunk)
	} else {
		slice = *maybeChunk.(*ChunkSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		args[0] = object.FileID
	} else {
		for i, obj := range slice {
			args[i] = obj.FileID
		}
	}

	query := fmt.Sprintf(
		"select * from \"files\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load File")
	}
	defer results.Close()

	var resultSlice []*File
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice File")
	}

	if len(fileAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		if object.R == nil {
			object.R = &chunkR{}
		}
		object.R.File = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.FileID == foreign.ID {
				if local.R == nil {
					local.R = &chunkR{}
				}
				local.R.File = foreign
				break
			}
		}
	}

	return nil
}









// SetFile of the chunk to the related item.
// Sets c.R.File to related.
// Adds c to related.R.Chunks.
func (c *Chunk) SetFile(exec boil.Executor, insert bool, related *File) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	oldVal := c.FileID
	c.FileID = related.ID

	if err = c.Update(exec, "file_id"); err != nil {
		c.FileID = oldVal

		return errors.Wrap(err, "failed to update local table")
	}

	if c.R == nil {
		c.R = &chunkR{
			File: related,
		}
	} else {
		c.R.File = related
	}

	if related.R == nil {
		related.R = &fileR{
			Chunks: ChunkSlice{c},
		}
	} else {
		related.R.Chunks = append(related.R.Chunks, c)
	}

	return nil
}




// ChunksG retrieves all records.
func ChunksG(mods ...qm.QueryMod) chunkQuery {
	return Chunks(boil.GetDB(), mods...)
}

// Chunks retrieves all the records using an executor.
func Chunks(exec boil.Executor, mods ...qm.QueryMod) chunkQuery {
	mods = append(mods, qm.From("\"chunks\""))
	return chunkQuery{NewQuery(exec, mods...)}
}

// FindChunkG retrieves a single record by ID.
func FindChunkG(id string, selectCols ...string) (*Chunk, error) {
	return FindChunk(boil.GetDB(), id, selectCols...)
}

// FindChunkGP retrieves a single record by ID, and panics on error.
func FindChunkGP(id string, selectCols ...string) *Chunk {
	retobj, err := FindChunk(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindChunk retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindChunk(exec boil.Executor, id string, selectCols ...string) (*Chunk, error) {
	chunkObj := &Chunk{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"chunks\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(chunkObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from chunks")
	}

	return chunkObj, nil
}

// FindChunkP retrieves a single record by ID with an executor, and panics on error.
func FindChunkP(exec boil.Executor, id string, selectCols ...string) *Chunk {
	retobj, err := FindChunk(exec, id, selectCols...)
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
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Chunk) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no chunks provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(chunkColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	chunkInsertCacheMut.RLock()
	cache, cached := chunkInsertCache[key]
	chunkInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			chunkColumns,
			chunkColumnsWithDefault,
			chunkColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(chunkType, chunkMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(chunkType, chunkMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"chunks\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into chunks")
	}

	if !cached {
		chunkInsertCacheMut.Lock()
		chunkInsertCache[key] = cache
		chunkInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
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
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Chunk) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	chunkUpdateCacheMut.RLock()
	cache, cached := chunkUpdateCache[key]
	chunkUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(chunkColumns, chunkPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update chunks, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"chunks\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, chunkPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(chunkType, chunkMapping, append(wl, chunkPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update chunks row")
	}

	if !cached {
		chunkUpdateCacheMut.Lock()
		chunkUpdateCache[key] = cache
		chunkUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q chunkQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q chunkQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for chunks")
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
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), chunkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"chunks\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(chunkPrimaryKeyColumns), len(colNames)+1, len(chunkPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in chunk slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Chunk) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Chunk) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Chunk) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Chunk) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no chunks provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(chunkColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	chunkUpsertCacheMut.RLock()
	cache, cached := chunkUpsertCache[key]
	chunkUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			chunkColumns,
			chunkColumnsWithDefault,
			chunkColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			chunkColumns,
			chunkPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert chunks, could not build update column list")
		}

		var conflict []string
		if len(conflictColumns) == 0 {
			conflict = make([]string, len(chunkPrimaryKeyColumns))
			copy(conflict, chunkPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"chunks\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(chunkType, chunkMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(chunkType, chunkMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	values := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, values...).Scan(returns...)
	} else {
		_, err = exec.Exec(cache.query, values...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for chunks")
	}

	if !cached {
		chunkUpsertCacheMut.Lock()
		chunkUpsertCache[key] = cache
		chunkUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
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
		return errors.New("models: no Chunk provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), chunkPrimaryKeyMapping)
	sql := "DELETE FROM \"chunks\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from chunks")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
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

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from chunks")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
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

	if len(chunkBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), chunkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"chunks\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, chunkPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(chunkPrimaryKeyColumns), 1, len(chunkPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from chunk slice")
	}

	if len(chunkAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
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
	ret, err := FindChunk(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ChunkSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ChunkSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ChunkSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ChunkSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ChunkSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	chunks := ChunkSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), chunkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"chunks\".* FROM \"chunks\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, chunkPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(chunkPrimaryKeyColumns), 1, len(chunkPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&chunks)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ChunkSlice")
	}

	*o = chunks

	return nil
}

// ChunkExists checks if the Chunk row exists.
func ChunkExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"chunks\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if chunks exists")
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


