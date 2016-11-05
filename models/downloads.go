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
	"gopkg.in/nullbio/null.v5"
)

// Download is an object representing the database table.
type Download struct {
	ID        int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Ip        null.String `boil:"ip" json:"ip,omitempty" toml:"ip" yaml:"ip,omitempty"`
	CacheHit  bool        `boil:"cache_hit" json:"cache_hit" toml:"cache_hit" yaml:"cache_hit"`
	FileID    null.String `boil:"file_id" json:"file_id,omitempty" toml:"file_id" yaml:"file_id,omitempty"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *downloadR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L downloadL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// downloadR is where relationships are stored.
type downloadR struct {
	File *File
}

// downloadL is where Load methods for each relationship are stored.
type downloadL struct{}

var (
	downloadColumns               = []string{"id", "ip", "cache_hit", "file_id", "created_at"}
	downloadColumnsWithoutDefault = []string{"ip", "cache_hit", "file_id", "created_at"}
	downloadColumnsWithDefault    = []string{"id"}
	downloadPrimaryKeyColumns     = []string{"id"}
)

type (
	// DownloadSlice is an alias for a slice of pointers to Download.
	// This should generally be used opposed to []Download.
	DownloadSlice []*Download
	// DownloadHook is the signature for custom Download hook methods
	DownloadHook func(boil.Executor, *Download) error

	downloadQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	downloadType                 = reflect.TypeOf(&Download{})
	downloadMapping              = queries.MakeStructMapping(downloadType)
	downloadPrimaryKeyMapping, _ = queries.BindMapping(downloadType, downloadMapping, downloadPrimaryKeyColumns)
	downloadInsertCacheMut       sync.RWMutex
	downloadInsertCache          = make(map[string]insertCache)
	downloadUpdateCacheMut       sync.RWMutex
	downloadUpdateCache          = make(map[string]updateCache)
	downloadUpsertCacheMut       sync.RWMutex
	downloadUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var downloadBeforeInsertHooks []DownloadHook
var downloadBeforeUpdateHooks []DownloadHook
var downloadBeforeDeleteHooks []DownloadHook
var downloadBeforeUpsertHooks []DownloadHook

var downloadAfterInsertHooks []DownloadHook
var downloadAfterSelectHooks []DownloadHook
var downloadAfterUpdateHooks []DownloadHook
var downloadAfterDeleteHooks []DownloadHook
var downloadAfterUpsertHooks []DownloadHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Download) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Download) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Download) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Download) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Download) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Download) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Download) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Download) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Download) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range downloadAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddDownloadHook registers your hook function for all future operations.
func AddDownloadHook(hookPoint boil.HookPoint, downloadHook DownloadHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		downloadBeforeInsertHooks = append(downloadBeforeInsertHooks, downloadHook)
	case boil.BeforeUpdateHook:
		downloadBeforeUpdateHooks = append(downloadBeforeUpdateHooks, downloadHook)
	case boil.BeforeDeleteHook:
		downloadBeforeDeleteHooks = append(downloadBeforeDeleteHooks, downloadHook)
	case boil.BeforeUpsertHook:
		downloadBeforeUpsertHooks = append(downloadBeforeUpsertHooks, downloadHook)
	case boil.AfterInsertHook:
		downloadAfterInsertHooks = append(downloadAfterInsertHooks, downloadHook)
	case boil.AfterSelectHook:
		downloadAfterSelectHooks = append(downloadAfterSelectHooks, downloadHook)
	case boil.AfterUpdateHook:
		downloadAfterUpdateHooks = append(downloadAfterUpdateHooks, downloadHook)
	case boil.AfterDeleteHook:
		downloadAfterDeleteHooks = append(downloadAfterDeleteHooks, downloadHook)
	case boil.AfterUpsertHook:
		downloadAfterUpsertHooks = append(downloadAfterUpsertHooks, downloadHook)
	}
}

// OneP returns a single download record from the query, and panics on error.
func (q downloadQuery) OneP() *Download {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single download record from the query.
func (q downloadQuery) One() (*Download, error) {
	o := &Download{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for downloads")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Download records from the query, and panics on error.
func (q downloadQuery) AllP() DownloadSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Download records from the query.
func (q downloadQuery) All() (DownloadSlice, error) {
	var o DownloadSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Download slice")
	}

	if len(downloadAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Download records in the query, and panics on error.
func (q downloadQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Download records in the query.
func (q downloadQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count downloads rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q downloadQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q downloadQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if downloads exists")
	}

	return count > 0, nil
}

// FileG pointed to by the foreign key.
func (o *Download) FileG(mods ...qm.QueryMod) fileQuery {
	return o.File(boil.GetDB(), mods...)
}

// File pointed to by the foreign key.
func (o *Download) File(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=$1", o.FileID),
	}

	queryMods = append(queryMods, mods...)

	query := Files(exec, queryMods...)
	queries.SetFrom(query.Query, "\"files\"")

	return query
}

// LoadFile allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (downloadL) LoadFile(e boil.Executor, singular bool, maybeDownload interface{}) error {
	var slice []*Download
	var object *Download

	count := 1
	if singular {
		object = maybeDownload.(*Download)
	} else {
		slice = *maybeDownload.(*DownloadSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		object.R = &downloadR{}
		args[0] = object.FileID
	} else {
		for i, obj := range slice {
			obj.R = &downloadR{}
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

	if len(downloadAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.File = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.FileID.String == foreign.ID {
				local.R.File = foreign
				break
			}
		}
	}

	return nil
}

// SetFile of the download to the related item.
// Sets o.R.File to related.
// Adds o to related.R.Downloads.
func (o *Download) SetFile(exec boil.Executor, insert bool, related *File) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"downloads\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"file_id"}),
		strmangle.WhereClause("\"", "\"", 2, downloadPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.FileID.String = related.ID
	o.FileID.Valid = true

	if o.R == nil {
		o.R = &downloadR{
			File: related,
		}
	} else {
		o.R.File = related
	}

	if related.R == nil {
		related.R = &fileR{
			Downloads: DownloadSlice{o},
		}
	} else {
		related.R.Downloads = append(related.R.Downloads, o)
	}

	return nil
}

// RemoveFile relationship.
// Sets o.R.File to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Download) RemoveFile(exec boil.Executor, related *File) error {
	var err error

	o.FileID.Valid = false
	if err = o.Update(exec, "file_id"); err != nil {
		o.FileID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.File = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Downloads {
		if o.FileID.String != ri.FileID.String {
			continue
		}

		ln := len(related.R.Downloads)
		if ln > 1 && i < ln-1 {
			related.R.Downloads[i] = related.R.Downloads[ln-1]
		}
		related.R.Downloads = related.R.Downloads[:ln-1]
		break
	}
	return nil
}

// DownloadsG retrieves all records.
func DownloadsG(mods ...qm.QueryMod) downloadQuery {
	return Downloads(boil.GetDB(), mods...)
}

// Downloads retrieves all the records using an executor.
func Downloads(exec boil.Executor, mods ...qm.QueryMod) downloadQuery {
	mods = append(mods, qm.From("\"downloads\""))
	return downloadQuery{NewQuery(exec, mods...)}
}

// FindDownloadG retrieves a single record by ID.
func FindDownloadG(id int, selectCols ...string) (*Download, error) {
	return FindDownload(boil.GetDB(), id, selectCols...)
}

// FindDownloadGP retrieves a single record by ID, and panics on error.
func FindDownloadGP(id int, selectCols ...string) *Download {
	retobj, err := FindDownload(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindDownload retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindDownload(exec boil.Executor, id int, selectCols ...string) (*Download, error) {
	downloadObj := &Download{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"downloads\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(downloadObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from downloads")
	}

	return downloadObj, nil
}

// FindDownloadP retrieves a single record by ID with an executor, and panics on error.
func FindDownloadP(exec boil.Executor, id int, selectCols ...string) *Download {
	retobj, err := FindDownload(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Download) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Download) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Download) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Download) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no downloads provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(downloadColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	downloadInsertCacheMut.RLock()
	cache, cached := downloadInsertCache[key]
	downloadInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			downloadColumns,
			downloadColumnsWithDefault,
			downloadColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(downloadType, downloadMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(downloadType, downloadMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"downloads\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into downloads")
	}

	if !cached {
		downloadInsertCacheMut.Lock()
		downloadInsertCache[key] = cache
		downloadInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Download record. See Update for
// whitelist behavior description.
func (o *Download) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Download record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Download) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Download, and panics on error.
// See Update for whitelist behavior description.
func (o *Download) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Download.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Download) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	downloadUpdateCacheMut.RLock()
	cache, cached := downloadUpdateCache[key]
	downloadUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(downloadColumns, downloadPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update downloads, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"downloads\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, downloadPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(downloadType, downloadMapping, append(wl, downloadPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update downloads row")
	}

	if !cached {
		downloadUpdateCacheMut.Lock()
		downloadUpdateCache[key] = cache
		downloadUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q downloadQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q downloadQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for downloads")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o DownloadSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o DownloadSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o DownloadSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o DownloadSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), downloadPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"downloads\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(downloadPrimaryKeyColumns), len(colNames)+1, len(downloadPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in download slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Download) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Download) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Download) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Download) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no downloads provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(downloadColumnsWithDefault, o)

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

	downloadUpsertCacheMut.RLock()
	cache, cached := downloadUpsertCache[key]
	downloadUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			downloadColumns,
			downloadColumnsWithDefault,
			downloadColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			downloadColumns,
			downloadPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert downloads, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(downloadPrimaryKeyColumns))
			copy(conflict, downloadPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"downloads\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(downloadType, downloadMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(downloadType, downloadMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for downloads")
	}

	if !cached {
		downloadUpsertCacheMut.Lock()
		downloadUpsertCache[key] = cache
		downloadUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Download record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Download) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Download record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Download) DeleteG() error {
	if o == nil {
		return errors.New("models: no Download provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Download record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Download) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Download record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Download) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Download provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), downloadPrimaryKeyMapping)
	sql := "DELETE FROM \"downloads\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from downloads")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q downloadQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q downloadQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no downloadQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from downloads")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o DownloadSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o DownloadSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Download slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o DownloadSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o DownloadSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Download slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(downloadBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), downloadPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"downloads\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, downloadPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(downloadPrimaryKeyColumns), 1, len(downloadPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from download slice")
	}

	if len(downloadAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Download) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Download) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Download) ReloadG() error {
	if o == nil {
		return errors.New("models: no Download provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Download) Reload(exec boil.Executor) error {
	ret, err := FindDownload(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DownloadSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DownloadSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DownloadSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty DownloadSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DownloadSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	downloads := DownloadSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), downloadPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"downloads\".* FROM \"downloads\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, downloadPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(downloadPrimaryKeyColumns), 1, len(downloadPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&downloads)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in DownloadSlice")
	}

	*o = downloads

	return nil
}

// DownloadExists checks if the Download row exists.
func DownloadExists(exec boil.Executor, id int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"downloads\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if downloads exists")
	}

	return exists, nil
}

// DownloadExistsG checks if the Download row exists.
func DownloadExistsG(id int) (bool, error) {
	return DownloadExists(boil.GetDB(), id)
}

// DownloadExistsGP checks if the Download row exists. Panics on error.
func DownloadExistsGP(id int) bool {
	e, err := DownloadExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// DownloadExistsP checks if the Download row exists. Panics on error.
func DownloadExistsP(exec boil.Executor, id int) bool {
	e, err := DownloadExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
