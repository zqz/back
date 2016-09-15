package models

import (
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

	R *fileR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L fileL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// fileR is where relationships are stored.
type fileR struct {
	Chunks     ChunkSlice
	Thumbnails ThumbnailSlice
}

// fileL is where Load methods for each relationship are stored.
type fileL struct{}

var (
	fileColumns               = []string{"id", "size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at", "slug"}
	fileColumnsWithoutDefault = []string{"size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at"}
	fileColumnsWithDefault    = []string{"id", "slug"}
	filePrimaryKeyColumns     = []string{"id"}
)

type (
	// FileSlice is an alias for a slice of pointers to File.
	// This should generally be used opposed to []File.
	FileSlice []*File
	// FileHook is the signature for custom File hook methods
	FileHook func(boil.Executor, *File) error

	fileQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	fileType           = reflect.TypeOf(&File{})
	fileMapping        = queries.MakeStructMapping(fileType)
	fileInsertCacheMut sync.RWMutex
	fileInsertCache    = make(map[string]insertCache)
	fileUpdateCacheMut sync.RWMutex
	fileUpdateCache    = make(map[string]updateCache)
	fileUpsertCacheMut sync.RWMutex
	fileUpsertCache    = make(map[string]insertCache)
)

// Force time package dependency for automated UpdatedAt/CreatedAt.
var _ = time.Second

var fileBeforeInsertHooks []FileHook
var fileBeforeUpdateHooks []FileHook
var fileBeforeDeleteHooks []FileHook
var fileBeforeUpsertHooks []FileHook

var fileAfterInsertHooks []FileHook
var fileAfterSelectHooks []FileHook
var fileAfterUpdateHooks []FileHook
var fileAfterDeleteHooks []FileHook
var fileAfterUpsertHooks []FileHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *File) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range fileBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *File) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range fileBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *File) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range fileBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *File) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range fileBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *File) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range fileAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *File) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range fileAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *File) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range fileAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *File) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range fileAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *File) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range fileAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddFileHook registers your hook function for all future operations.
func AddFileHook(hookPoint boil.HookPoint, fileHook FileHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		fileBeforeInsertHooks = append(fileBeforeInsertHooks, fileHook)
	case boil.BeforeUpdateHook:
		fileBeforeUpdateHooks = append(fileBeforeUpdateHooks, fileHook)
	case boil.BeforeDeleteHook:
		fileBeforeDeleteHooks = append(fileBeforeDeleteHooks, fileHook)
	case boil.BeforeUpsertHook:
		fileBeforeUpsertHooks = append(fileBeforeUpsertHooks, fileHook)
	case boil.AfterInsertHook:
		fileAfterInsertHooks = append(fileAfterInsertHooks, fileHook)
	case boil.AfterSelectHook:
		fileAfterSelectHooks = append(fileAfterSelectHooks, fileHook)
	case boil.AfterUpdateHook:
		fileAfterUpdateHooks = append(fileAfterUpdateHooks, fileHook)
	case boil.AfterDeleteHook:
		fileAfterDeleteHooks = append(fileAfterDeleteHooks, fileHook)
	case boil.AfterUpsertHook:
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

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for files")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
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
		return nil, errors.Wrap(err, "models: failed to assign all query results to File slice")
	}

	if len(fileAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
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

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count files rows")
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

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if files exists")
	}

	return count > 0, nil
}


// ChunksG retrieves all the file's chunks.
func (f *File) ChunksG(mods ...qm.QueryMod) chunkQuery {
	return f.Chunks(boil.GetDB(), mods...)
}

// Chunks retrieves all the file's chunks with an executor.
func (f *File) Chunks(exec boil.Executor, mods ...qm.QueryMod) chunkQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"file_id\"=$1", f.ID),
	)

	query := Chunks(exec, queryMods...)
	queries.SetFrom(query.Query, "\"chunks\" as \"a\"")
	return query
}

// ThumbnailsG retrieves all the file's thumbnails.
func (f *File) ThumbnailsG(mods ...qm.QueryMod) thumbnailQuery {
	return f.Thumbnails(boil.GetDB(), mods...)
}

// Thumbnails retrieves all the file's thumbnails with an executor.
func (f *File) Thumbnails(exec boil.Executor, mods ...qm.QueryMod) thumbnailQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"file_id\"=$1", f.ID),
	)

	query := Thumbnails(exec, queryMods...)
	queries.SetFrom(query.Query, "\"thumbnails\" as \"a\"")
	return query
}




// LoadChunks allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (fileL) LoadChunks(e boil.Executor, singular bool, maybeFile interface{}) error {
	var slice []*File
	var object *File

	count := 1
	if singular {
		object = maybeFile.(*File)
	} else {
		slice = *maybeFile.(*FileSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"chunks\" where \"file_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load chunks")
	}
	defer results.Close()

	var resultSlice []*Chunk
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice chunks")
	}

	if len(chunkAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		if object.R == nil {
			object.R = &fileR{}
		}
		object.R.Chunks = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.FileID {
				if local.R == nil {
					local.R = &fileR{}
				}
				local.R.Chunks = append(local.R.Chunks, foreign)
				break
			}
		}
	}

	return nil
}

// LoadThumbnails allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (fileL) LoadThumbnails(e boil.Executor, singular bool, maybeFile interface{}) error {
	var slice []*File
	var object *File

	count := 1
	if singular {
		object = maybeFile.(*File)
	} else {
		slice = *maybeFile.(*FileSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"thumbnails\" where \"file_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load thumbnails")
	}
	defer results.Close()

	var resultSlice []*Thumbnail
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice thumbnails")
	}

	if len(thumbnailAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		if object.R == nil {
			object.R = &fileR{}
		}
		object.R.Thumbnails = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.FileID {
				if local.R == nil {
					local.R = &fileR{}
				}
				local.R.Thumbnails = append(local.R.Thumbnails, foreign)
				break
			}
		}
	}

	return nil
}





// AddChunks adds the given related objects to the existing relationships
// of the file, optionally inserting them as new records.
// Appends related to f.R.Chunks.
// Sets related.R.File appropriately.
func (f *File) AddChunks(exec boil.Executor, insert bool, related ...*Chunk) error {
	var err error
	for _, rel := range related {
		rel.FileID = f.ID
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			if err = rel.Update(exec, "file_id"); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}
		}
	}

	if f.R == nil {
		f.R = &fileR{
			Chunks: related,
		}
	} else {
		f.R.Chunks = append(f.R.Chunks, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &chunkR{
				File: f,
			}
		} else {
			rel.R.File = f
		}
	}
	return nil
}

// AddThumbnails adds the given related objects to the existing relationships
// of the file, optionally inserting them as new records.
// Appends related to f.R.Thumbnails.
// Sets related.R.File appropriately.
func (f *File) AddThumbnails(exec boil.Executor, insert bool, related ...*Thumbnail) error {
	var err error
	for _, rel := range related {
		rel.FileID = f.ID
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			if err = rel.Update(exec, "file_id"); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}
		}
	}

	if f.R == nil {
		f.R = &fileR{
			Thumbnails: related,
		}
	} else {
		f.R.Thumbnails = append(f.R.Thumbnails, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &thumbnailR{
				File: f,
			}
		} else {
			rel.R.File = f
		}
	}
	return nil
}
// FilesG retrieves all records.
func FilesG(mods ...qm.QueryMod) fileQuery {
	return Files(boil.GetDB(), mods...)
}

// Files retrieves all the records using an executor.
func Files(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	mods = append(mods, qm.From("\"files\""))
	return fileQuery{NewQuery(exec, mods...)}
}

// FindFileG retrieves a single record by ID.
func FindFileG(id string, selectCols ...string) (*File, error) {
	return FindFile(boil.GetDB(), id, selectCols...)
}

// FindFileGP retrieves a single record by ID, and panics on error.
func FindFileGP(id string, selectCols ...string) *File {
	retobj, err := FindFile(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindFile retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFile(exec boil.Executor, id string, selectCols ...string) (*File, error) {
	fileObj := &File{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"files\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(fileObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from files")
	}

	return fileObj, nil
}

// FindFileP retrieves a single record by ID with an executor, and panics on error.
func FindFileP(exec boil.Executor, id string, selectCols ...string) *File {
	retobj, err := FindFile(exec, id, selectCols...)
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
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *File) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no files provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(fileColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	fileInsertCacheMut.RLock()
	cache, cached := fileInsertCache[key]
	fileInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			fileColumns,
			fileColumnsWithDefault,
			fileColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(fileType, fileMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(fileType, fileMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"files\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into files")
	}

	if !cached {
		fileInsertCacheMut.Lock()
		fileInsertCache[key] = cache
		fileInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
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
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *File) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	fileUpdateCacheMut.RLock()
	cache, cached := fileUpdateCache[key]
	fileUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(fileColumns, filePrimaryKeyColumns, whitelist)

		cache.query = fmt.Sprintf("UPDATE \"files\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, filePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(fileType, fileMapping, append(wl, filePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	if len(cache.valueMapping) == 0 {
		return errors.New("models: unable to update files, could not build whitelist")
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	result, err := exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update files row")
	}

	if r, err := result.RowsAffected(); err == nil && r != 1 {
		return errors.Errorf("failed to update single row, updated %d rows", r)
	}

	if !cached {
		fileUpdateCacheMut.Lock()
		fileUpdateCache[key] = cache
		fileUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q fileQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q fileQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for files")
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
	args = append(args, o.inPrimaryKeyArgs()...)

	sql := fmt.Sprintf(
		"UPDATE \"files\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(filePrimaryKeyColumns), len(colNames)+1, len(filePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in file slice")
	}

	if r, err := result.RowsAffected(); err == nil && r != ln {
		return errors.Errorf("failed to update %d rows, only affected %d", ln, r)
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *File) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *File) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *File) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *File) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no files provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

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
	key := buf.String()
	strmangle.PutBuffer(buf)

	fileUpsertCacheMut.RLock()
	cache, cached := fileUpsertCache[key]
	fileUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			fileColumns,
			fileColumnsWithDefault,
			fileColumnsWithoutDefault,
			queries.NonZeroDefaultSet(fileColumnsWithDefault, o),
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			fileColumns,
			filePrimaryKeyColumns,
			updateColumns,
		)

		var conflict []string
		if len(conflictColumns) == 0 {
			conflict = make([]string, len(filePrimaryKeyColumns))
			copy(conflict, filePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"files\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(fileType, fileMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(fileType, fileMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for files")
	}

	if !cached {
		fileUpsertCacheMut.Lock()
		fileUpsertCache[key] = cache
		fileUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
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
		return errors.New("models: no File provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := o.inPrimaryKeyArgs()

	sql := "DELETE FROM \"files\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from files")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
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

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from files")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
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

	if len(fileBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		"DELETE FROM \"files\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, filePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(filePrimaryKeyColumns), 1, len(filePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from file slice")
	}

	if len(fileAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
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
	ret, err := FindFile(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *FileSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *FileSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty FileSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	files := FileSlice{}
	args := o.inPrimaryKeyArgs()

	sql := fmt.Sprintf(
		"SELECT \"files\".* FROM \"files\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, filePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(filePrimaryKeyColumns), 1, len(filePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&files)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in FileSlice")
	}

	*o = files

	return nil
}

// FileExists checks if the File row exists.
func FileExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"files\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if files exists")
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


