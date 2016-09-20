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

// Thumbnail is an object representing the database table.
type Thumbnail struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	FileID    string    `boil:"file_id" json:"file_id" toml:"file_id" yaml:"file_id"`
	Size      int       `boil:"size" json:"size" toml:"size" yaml:"size"`
	Hash      string    `boil:"hash" json:"hash" toml:"hash" yaml:"hash"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *thumbnailR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L thumbnailL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// thumbnailR is where relationships are stored.
type thumbnailR struct {
	File *File
}

// thumbnailL is where Load methods for each relationship are stored.
type thumbnailL struct{}

var (
	thumbnailColumns               = []string{"id", "file_id", "size", "hash", "created_at", "updated_at"}
	thumbnailColumnsWithoutDefault = []string{"file_id", "size", "hash", "created_at", "updated_at"}
	thumbnailColumnsWithDefault    = []string{"id"}
	thumbnailPrimaryKeyColumns     = []string{"id"}
)

type (
	// ThumbnailSlice is an alias for a slice of pointers to Thumbnail.
	// This should generally be used opposed to []Thumbnail.
	ThumbnailSlice []*Thumbnail
	// ThumbnailHook is the signature for custom Thumbnail hook methods
	ThumbnailHook func(boil.Executor, *Thumbnail) error

	thumbnailQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	thumbnailType                 = reflect.TypeOf(&Thumbnail{})
	thumbnailMapping              = queries.MakeStructMapping(thumbnailType)
	thumbnailPrimaryKeyMapping, _ = queries.BindMapping(thumbnailType, thumbnailMapping, thumbnailPrimaryKeyColumns)
	thumbnailInsertCacheMut       sync.RWMutex
	thumbnailInsertCache          = make(map[string]insertCache)
	thumbnailUpdateCacheMut       sync.RWMutex
	thumbnailUpdateCache          = make(map[string]updateCache)
	thumbnailUpsertCacheMut       sync.RWMutex
	thumbnailUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

var thumbnailBeforeInsertHooks []ThumbnailHook
var thumbnailBeforeUpdateHooks []ThumbnailHook
var thumbnailBeforeDeleteHooks []ThumbnailHook
var thumbnailBeforeUpsertHooks []ThumbnailHook

var thumbnailAfterInsertHooks []ThumbnailHook
var thumbnailAfterSelectHooks []ThumbnailHook
var thumbnailAfterUpdateHooks []ThumbnailHook
var thumbnailAfterDeleteHooks []ThumbnailHook
var thumbnailAfterUpsertHooks []ThumbnailHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Thumbnail) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Thumbnail) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Thumbnail) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Thumbnail) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Thumbnail) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Thumbnail) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Thumbnail) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Thumbnail) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Thumbnail) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range thumbnailAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddThumbnailHook registers your hook function for all future operations.
func AddThumbnailHook(hookPoint boil.HookPoint, thumbnailHook ThumbnailHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		thumbnailBeforeInsertHooks = append(thumbnailBeforeInsertHooks, thumbnailHook)
	case boil.BeforeUpdateHook:
		thumbnailBeforeUpdateHooks = append(thumbnailBeforeUpdateHooks, thumbnailHook)
	case boil.BeforeDeleteHook:
		thumbnailBeforeDeleteHooks = append(thumbnailBeforeDeleteHooks, thumbnailHook)
	case boil.BeforeUpsertHook:
		thumbnailBeforeUpsertHooks = append(thumbnailBeforeUpsertHooks, thumbnailHook)
	case boil.AfterInsertHook:
		thumbnailAfterInsertHooks = append(thumbnailAfterInsertHooks, thumbnailHook)
	case boil.AfterSelectHook:
		thumbnailAfterSelectHooks = append(thumbnailAfterSelectHooks, thumbnailHook)
	case boil.AfterUpdateHook:
		thumbnailAfterUpdateHooks = append(thumbnailAfterUpdateHooks, thumbnailHook)
	case boil.AfterDeleteHook:
		thumbnailAfterDeleteHooks = append(thumbnailAfterDeleteHooks, thumbnailHook)
	case boil.AfterUpsertHook:
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

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for thumbnails")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
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
		return nil, errors.Wrap(err, "models: failed to assign all query results to Thumbnail slice")
	}

	if len(thumbnailAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
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

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count thumbnails rows")
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

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if thumbnails exists")
	}

	return count > 0, nil
}

// FileG pointed to by the foreign key.
func (t *Thumbnail) FileG(mods ...qm.QueryMod) fileQuery {
	return t.File(boil.GetDB(), mods...)
}

// File pointed to by the foreign key.
func (t *Thumbnail) File(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=$1", t.FileID),
	}

	queryMods = append(queryMods, mods...)

	query := Files(exec, queryMods...)
	queries.SetFrom(query.Query, "\"files\"")

	return query
}



// LoadFile allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (thumbnailL) LoadFile(e boil.Executor, singular bool, maybeThumbnail interface{}) error {
	var slice []*Thumbnail
	var object *Thumbnail

	count := 1
	if singular {
		object = maybeThumbnail.(*Thumbnail)
	} else {
		slice = *maybeThumbnail.(*ThumbnailSlice)
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
			object.R = &thumbnailR{}
		}
		object.R.File = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.FileID == foreign.ID {
				if local.R == nil {
					local.R = &thumbnailR{}
				}
				local.R.File = foreign
				break
			}
		}
	}

	return nil
}









// SetFile of the thumbnail to the related item.
// Sets t.R.File to related.
// Adds t to related.R.Thumbnails.
func (t *Thumbnail) SetFile(exec boil.Executor, insert bool, related *File) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	oldVal := t.FileID
	t.FileID = related.ID

	if err = t.Update(exec, "file_id"); err != nil {
		t.FileID = oldVal

		return errors.Wrap(err, "failed to update local table")
	}

	if t.R == nil {
		t.R = &thumbnailR{
			File: related,
		}
	} else {
		t.R.File = related
	}

	if related.R == nil {
		related.R = &fileR{
			Thumbnails: ThumbnailSlice{t},
		}
	} else {
		related.R.Thumbnails = append(related.R.Thumbnails, t)
	}

	return nil
}




// ThumbnailsG retrieves all records.
func ThumbnailsG(mods ...qm.QueryMod) thumbnailQuery {
	return Thumbnails(boil.GetDB(), mods...)
}

// Thumbnails retrieves all the records using an executor.
func Thumbnails(exec boil.Executor, mods ...qm.QueryMod) thumbnailQuery {
	mods = append(mods, qm.From("\"thumbnails\""))
	return thumbnailQuery{NewQuery(exec, mods...)}
}

// FindThumbnailG retrieves a single record by ID.
func FindThumbnailG(id string, selectCols ...string) (*Thumbnail, error) {
	return FindThumbnail(boil.GetDB(), id, selectCols...)
}

// FindThumbnailGP retrieves a single record by ID, and panics on error.
func FindThumbnailGP(id string, selectCols ...string) *Thumbnail {
	retobj, err := FindThumbnail(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindThumbnail retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindThumbnail(exec boil.Executor, id string, selectCols ...string) (*Thumbnail, error) {
	thumbnailObj := &Thumbnail{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"thumbnails\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(thumbnailObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from thumbnails")
	}

	return thumbnailObj, nil
}

// FindThumbnailP retrieves a single record by ID with an executor, and panics on error.
func FindThumbnailP(exec boil.Executor, id string, selectCols ...string) *Thumbnail {
	retobj, err := FindThumbnail(exec, id, selectCols...)
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
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Thumbnail) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no thumbnails provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(thumbnailColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	thumbnailInsertCacheMut.RLock()
	cache, cached := thumbnailInsertCache[key]
	thumbnailInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			thumbnailColumns,
			thumbnailColumnsWithDefault,
			thumbnailColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(thumbnailType, thumbnailMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(thumbnailType, thumbnailMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"thumbnails\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into thumbnails")
	}

	if !cached {
		thumbnailInsertCacheMut.Lock()
		thumbnailInsertCache[key] = cache
		thumbnailInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
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
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Thumbnail) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	thumbnailUpdateCacheMut.RLock()
	cache, cached := thumbnailUpdateCache[key]
	thumbnailUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(thumbnailColumns, thumbnailPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update thumbnails, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"thumbnails\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, thumbnailPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(thumbnailType, thumbnailMapping, append(wl, thumbnailPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update thumbnails row")
	}

	if !cached {
		thumbnailUpdateCacheMut.Lock()
		thumbnailUpdateCache[key] = cache
		thumbnailUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q thumbnailQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q thumbnailQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for thumbnails")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), thumbnailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"thumbnails\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(thumbnailPrimaryKeyColumns), len(colNames)+1, len(thumbnailPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in thumbnail slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Thumbnail) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Thumbnail) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Thumbnail) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Thumbnail) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no thumbnails provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(thumbnailColumnsWithDefault, o)

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

	thumbnailUpsertCacheMut.RLock()
	cache, cached := thumbnailUpsertCache[key]
	thumbnailUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			thumbnailColumns,
			thumbnailColumnsWithDefault,
			thumbnailColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			thumbnailColumns,
			thumbnailPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert thumbnails, could not build update column list")
		}

		var conflict []string
		if len(conflictColumns) == 0 {
			conflict = make([]string, len(thumbnailPrimaryKeyColumns))
			copy(conflict, thumbnailPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"thumbnails\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(thumbnailType, thumbnailMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(thumbnailType, thumbnailMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for thumbnails")
	}

	if !cached {
		thumbnailUpsertCacheMut.Lock()
		thumbnailUpsertCache[key] = cache
		thumbnailUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
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
		return errors.New("models: no Thumbnail provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), thumbnailPrimaryKeyMapping)
	sql := "DELETE FROM \"thumbnails\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from thumbnails")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
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

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from thumbnails")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
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

	if len(thumbnailBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), thumbnailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"thumbnails\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, thumbnailPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(thumbnailPrimaryKeyColumns), 1, len(thumbnailPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from thumbnail slice")
	}

	if len(thumbnailAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
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
	ret, err := FindThumbnail(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ThumbnailSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ThumbnailSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ThumbnailSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ThumbnailSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ThumbnailSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	thumbnails := ThumbnailSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), thumbnailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"thumbnails\".* FROM \"thumbnails\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, thumbnailPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(thumbnailPrimaryKeyColumns), 1, len(thumbnailPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&thumbnails)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ThumbnailSlice")
	}

	*o = thumbnails

	return nil
}

// ThumbnailExists checks if the Thumbnail row exists.
func ThumbnailExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"thumbnails\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if thumbnails exists")
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


