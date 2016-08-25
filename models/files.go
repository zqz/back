package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
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

	Loaded *FileLoaded `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// FileLoaded are where relationships are eagerly loaded.
type FileLoaded struct {
	Chunks     ChunkSlice
	Thumbnails ThumbnailSlice
}

var (
	fileColumns               = []string{"id", "size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at", "slug"}
	fileColumnsWithoutDefault = []string{"size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at"}
	fileColumnsWithDefault    = []string{"id", "slug"}
	filePrimaryKeyColumns     = []string{"id"}
	fileTitleCases            = map[string]string{
		"id":         "ID",
		"size":       "Size",
		"num_chunks": "NumChunks",
		"state":      "State",
		"name":       "Name",
		"hash":       "Hash",
		"type":       "Type",
		"created_at": "CreatedAt",
		"updated_at": "UpdatedAt",
		"slug":       "Slug",
	}
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

	err := q.BindFast(o, fileTitleCases)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for files")
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

	err := q.BindFast(&o, fileTitleCases)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to File slice")
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

	boil.SetCount(q.Query)
	boil.SetLimit(q.Query, 1)

	err := boil.ExecQueryOne(q.Query).Scan(&count)
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
	return query
}

// ThumbnailsG retrieves all the file's thumbnails.
func (f *File) ThumbnailsG(mods ...qm.QueryMod) thumbnailQuery {
	return f.Thumbnails(boil.GetDB(), mods...)
}

// Thumbnails retrieves all the file's thumbnails with an executor.
func (f *File) Thumbnails(exec boil.Executor, mods ...qm.QueryMod) thumbnailQuery {
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
	return query
}



// LoadChunks allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (r *FileLoaded) LoadChunks(e boil.Executor, singular bool, maybeFile interface{}) error {
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
		`select * from "chunks" where "file_id" in (%s)`,
		strmangle.Placeholders(count, 1, 1),
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
	if err = boil.BindFast(results, &resultSlice, fileTitleCases); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice chunks")
	}

	if singular {
		if object.Loaded == nil {
			object.Loaded = &FileLoaded{}
		}
		object.Loaded.Chunks = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.FileID {
				if local.Loaded == nil {
					local.Loaded = &FileLoaded{}
				}
				local.Loaded.Chunks = append(local.Loaded.Chunks, foreign)
				break
			}
		}
	}

	return nil
}

// LoadThumbnails allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (r *FileLoaded) LoadThumbnails(e boil.Executor, singular bool, maybeFile interface{}) error {
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
		`select * from "thumbnails" where "file_id" in (%s)`,
		strmangle.Placeholders(count, 1, 1),
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
	if err = boil.BindFast(results, &resultSlice, fileTitleCases); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice thumbnails")
	}

	if singular {
		if object.Loaded == nil {
			object.Loaded = &FileLoaded{}
		}
		object.Loaded.Thumbnails = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.FileID {
				if local.Loaded == nil {
					local.Loaded = &FileLoaded{}
				}
				local.Loaded.Thumbnails = append(local.Loaded.Thumbnails, foreign)
				break
			}
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
	fileObj := &File{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(selectCols), ",")
	}
	query := fmt.Sprintf(
		`select %s from "files" where "id"=$1`, sel,
	)

	q := boil.SQL(query, id)
	boil.SetExecutor(q, exec)

	err := q.BindFast(fileObj, fileTitleCases)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from files")
	}

	return fileObj, nil
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

	wl, returnColumns := strmangle.InsertColumnSet(
		fileColumns,
		fileColumnsWithDefault,
		fileColumnsWithoutDefault,
		boil.NonZeroDefaultSet(fileColumnsWithDefault, fileTitleCases, o),
		whitelist,
	)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO files ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, fileTitleCases, wl...)...).Scan(boil.GetStructPointers(o, fileTitleCases, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, boil.GetStructValues(o, fileTitleCases, wl...)...)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, fileTitleCases, wl...))
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into files")
	}

	return o.doAfterCreateHooks()
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
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := strmangle.UpdateColumnSet(fileColumns, filePrimaryKeyColumns, whitelist)
	if len(wl) == 0 {
		return errors.New("models: unable to update files, could not build whitelist")
	}

	query = fmt.Sprintf(`UPDATE files SET %s WHERE %s`, strmangle.SetParamNames(wl), strmangle.WhereClause(len(wl)+1, filePrimaryKeyColumns))
	values = boil.GetStructValues(o, fileTitleCases, wl...)
	values = append(values, o.ID)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	result, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update files row")
	}

	if r, err := result.RowsAffected(); err == nil && r != 1 {
		return errors.Errorf("failed to update single row, updated %d rows", r)
	}

	return o.doAfterUpdateHooks()
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
		colNames[i] = strmangle.IdentQuote(name)
		args[i] = value
		i++
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
func (o *File) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no files provided for upsert")
	}

	var ret []string
	whitelist, ret = strmangle.InsertColumnSet(
		fileColumns,
		fileColumnsWithDefault,
		fileColumnsWithoutDefault,
		boil.NonZeroDefaultSet(fileColumnsWithDefault, fileTitleCases, o),
		whitelist,
	)
	update := strmangle.UpdateColumnSet(
		fileColumns,
		filePrimaryKeyColumns,
		updateColumns,
	)
	conflict := conflictColumns
	if len(conflict) == 0 {
		conflict = make([]string, len(filePrimaryKeyColumns))
		copy(conflict, filePrimaryKeyColumns)
	}

	query := generateUpsertQuery("files", updateOnConflict, ret, update, conflict, whitelist)

	var err error
	if err := o.doBeforeUpsertHooks(); err != nil {
		return err
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, fileTitleCases, whitelist...))
	}
	if len(ret) != 0 {
		err = exec.QueryRow(query, boil.GetStructValues(o, fileTitleCases, whitelist...)...).Scan(boil.GetStructPointers(o, fileTitleCases, ret...)...)
	} else {
		_, err = exec.Exec(query, boil.GetStructValues(o, fileTitleCases, whitelist...)...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for files")
	}

	if err := o.doAfterUpsertHooks(); err != nil {
		return err
	}

	return nil
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

	args := o.inPrimaryKeyArgs()

	sql := `DELETE FROM files WHERE "id"=$1`

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from files")
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
		return errors.Wrap(err, "models: unable to delete all from files")
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

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from file slice")
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
	if o == nil || len(*o) == 0 {
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

	err := q.BindFast(&files, fileTitleCases)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in FileSlice")
	}

	*o = files

	return nil
}

// FileExists checks if the File row exists.
func FileExists(exec boil.Executor, id string) (bool, error) {
	var exists bool

	sql := `select exists(select 1 from "files" where "id"=$1 limit 1)`

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

