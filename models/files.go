package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
)

// File is an object representing the database table.
type File struct {
	ID        string    `db:"file_id" json:"id"`
	Size      int       `db:"file_size" json:"size"`
	NumChunks int       `db:"file_num_chunks" json:"num_chunks"`
	State     int       `db:"file_state" json:"state"`
	Name      string    `db:"file_name" json:"name"`
	Hash      string    `db:"file_hash" json:"hash"`
	Type      string    `db:"file_type" json:"type"`
	CreatedAt time.Time `db:"file_created_at" json:"created_at"`
	UpdatedAt time.Time `db:"file_updated_at" json:"updated_at"`
	Slug      string    `db:"file_slug" json:"slug"`
}

var (
	fileColumns                  = []string{"id", "size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at", "slug"}
	fileColumnsWithoutDefault    = []string{"size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at"}
	fileColumnsWithDefault       = []string{"id", "slug"}
	fileColumnsWithSimpleDefault = []string{}
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
var fileAfterCreateHooks []FileHook
var fileAfterUpdateHooks []FileHook

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

func FileAddHook(hookPoint boil.HookPoint, fileHook FileHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		fileBeforeCreateHooks = append(fileBeforeCreateHooks, fileHook)
	case boil.HookBeforeUpdate:
		fileBeforeUpdateHooks = append(fileBeforeUpdateHooks, fileHook)
	case boil.HookAfterCreate:
		fileAfterCreateHooks = append(fileAfterCreateHooks, fileHook)
	case boil.HookAfterUpdate:
		fileAfterUpdateHooks = append(fileAfterUpdateHooks, fileHook)
	}
}

// One returns a single file record from the query.
func (q fileQuery) One() (*File, error) {
	o := &File{}

	boil.SetLimit(q.Query, 1)

	res := boil.ExecQueryOne(q.Query)
	err := boil.BindOne(res, boil.Select(q.Query), o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for files: %s", err)
	}

	return o, nil
}

// OneP returns a single file record from the query, and panics on error.
func (q fileQuery) OneP() *File {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all File records from the query.
func (q fileQuery) All() (FileSlice, error) {
	var o FileSlice

	res, err := boil.ExecQueryAll(q.Query)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute an all query for files: %s", err)
	}
	defer res.Close()

	err = boil.BindAll(res, boil.Select(q.Query), &o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to File slice: %s", err)
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

// CountP returns the count of all File records in the query, and panics on error.
func (q fileQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}


// ChunksG retrieves all the file's chunks.
func (f *File) ChunksG(selectCols ...string) (ChunkSlice, error) {
	return f.Chunks(boil.GetDB(), selectCols...)
}

// ChunksGP panics on error. Retrieves all the file's chunks.
func (f *File) ChunksGP(selectCols ...string) ChunkSlice {
	o, err := f.Chunks(boil.GetDB(), selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// ChunksP panics on error. Retrieves all the file's chunks with an executor.
func (f *File) ChunksP(exec boil.Executor, selectCols ...string) ChunkSlice {
	o, err := f.Chunks(exec, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// Chunks retrieves all the file's chunks with an executor.
func (f *File) Chunks(exec boil.Executor, selectCols ...string) (ChunkSlice, error) {
	var ret ChunkSlice

	selectColumns := `"a".*`
	if len(selectCols) != 0 {
		selectColumns = `"a".` + strings.Join(selectCols, `","a"."`)
	}
	query := fmt.Sprintf(`select %s from chunks "a" where "a"."file_id"=$1`, selectColumns)

	rows, err := exec.Query(query, f.ID)
	if err != nil {
		return nil, fmt.Errorf(`models: unable to select from chunks: %v`, err)
	}
	defer rows.Close()

	for rows.Next() {
		next := new(Chunk)

		err = rows.Scan(boil.GetStructPointers(next, selectCols...)...)
		if err != nil {
			return nil, fmt.Errorf(`models: unable to scan into Chunk: %v`, err)
		}

		ret = append(ret, next)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(`models: unable to select from Chunk: %v`, err)
	}

	return ret, nil
}

// ThumbnailsG retrieves all the file's thumbnails.
func (f *File) ThumbnailsG(selectCols ...string) (ThumbnailSlice, error) {
	return f.Thumbnails(boil.GetDB(), selectCols...)
}

// ThumbnailsGP panics on error. Retrieves all the file's thumbnails.
func (f *File) ThumbnailsGP(selectCols ...string) ThumbnailSlice {
	o, err := f.Thumbnails(boil.GetDB(), selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// ThumbnailsP panics on error. Retrieves all the file's thumbnails with an executor.
func (f *File) ThumbnailsP(exec boil.Executor, selectCols ...string) ThumbnailSlice {
	o, err := f.Thumbnails(exec, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// Thumbnails retrieves all the file's thumbnails with an executor.
func (f *File) Thumbnails(exec boil.Executor, selectCols ...string) (ThumbnailSlice, error) {
	var ret ThumbnailSlice

	selectColumns := `"a".*`
	if len(selectCols) != 0 {
		selectColumns = `"a".` + strings.Join(selectCols, `","a"."`)
	}
	query := fmt.Sprintf(`select %s from thumbnails "a" where "a"."file_id"=$1`, selectColumns)

	rows, err := exec.Query(query, f.ID)
	if err != nil {
		return nil, fmt.Errorf(`models: unable to select from thumbnails: %v`, err)
	}
	defer rows.Close()

	for rows.Next() {
		next := new(Thumbnail)

		err = rows.Scan(boil.GetStructPointers(next, selectCols...)...)
		if err != nil {
			return nil, fmt.Errorf(`models: unable to scan into Thumbnail: %v`, err)
		}

		ret = append(ret, next)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(`models: unable to select from Thumbnail: %v`, err)
	}

	return ret, nil
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
	o, err := FileFind(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// FileFind retrieves a single record by ID with an executor.
func FileFind(exec boil.Executor, id string, selectCols ...string) (*File, error) {
	file := &File{}

	mods := []qm.QueryMod{
		qm.Select(selectCols...),
		qm.From("files"),
		qm.Where(`"id"=$1`, id),
	}

	q := NewQuery(exec, mods...)

	err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers(file, selectCols...)...)

	if err != nil {
		return nil, fmt.Errorf("models: unable to select from files: %v", err)
	}

	return file, nil
}

// FileFindP retrieves a single record by ID with an executor, and panics on error.
func FileFindP(exec boil.Executor, id string, selectCols ...string) *File {
	o, err := FileFind(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// InsertG a single record.
func (o *File) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error.
func (o *File) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
func (o *File) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no files provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO files ("%s") VALUES (%s)`, strings.Join(wl, `","`), boil.GenerateParamFlags(len(wl), 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.Size, o.NumChunks, o.State, o.Name, o.Hash, o.Type, o.CreatedAt, o.UpdatedAt, o.Slug)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into files: %s", err)
	}

	if err := o.doAfterCreateHooks(); err != nil {
		return err
	}

	return nil
}

// InsertP a single record using an executor, and panics on error.
func (o *File) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
func (o *File) generateInsertColumns(whitelist ...string) ([]string, []string) {
	var wl []string

	wl = append(wl, whitelist...)
	if len(whitelist) == 0 {
		wl = append(wl, fileColumnsWithoutDefault...)
	}

	wl = append(boil.NonZeroDefaultSet(fileColumnsWithDefault, o), wl...)
	wl = boil.SortByKeys(fileColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(fileColumnsWithDefault, wl)

	return wl, rc
}


// UpdateG a single File record.
// UpdateG takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
func (o *File) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single File record.
// UpdateGP takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
// Panics on error.
func (o *File) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the File.
func (o *File) Update(exec boil.Executor, whitelist ...string) error {
	return o.UpdateAt(exec, o.ID, whitelist...)
}

// UpdateP uses an executor to update the File, and panics on error.
func (o *File) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.UpdateAt(exec, o.ID, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAtG updates the File using the primary key to find the row to update.
func (o *File) UpdateAtG(id string, whitelist ...string) error {
	return o.UpdateAt(boil.GetDB(), id, whitelist...)
}

// UpdateAtGP updates the File using the primary key to find the row to update. Panics on error.
func (o *File) UpdateAtGP(id string, whitelist ...string) {
	if err := o.UpdateAt(boil.GetDB(), id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAt uses an executor to update the File using the primary key to find the row to update.
func (o *File) UpdateAt(exec boil.Executor, id string, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE files SET %s WHERE %s`, boil.SetParamNames(wl), boil.WherePrimaryKey(len(wl)+1, "id"))
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

// UpdateAtP uses an executor to update the File using the primary key to find the row to update.
// Panics on error.
func (o *File) UpdateAtP(exec boil.Executor, id string, whitelist ...string) {
	if err := o.UpdateAt(exec, id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with matching column names.
func (q fileQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for files: %s", err)
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q fileQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateUpdateColumns generates the whitelist columns for an update statement
func (o *File) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	var wl []string
	cols := fileColumnsWithoutDefault
	cols = append(boil.NonZeroDefaultSet(fileColumnsWithDefault, o), cols...)
	// Subtract primary keys and autoincrement columns
	cols = boil.SetComplement(cols, filePrimaryKeyColumns)
	cols = boil.SetComplement(cols, fileAutoIncrementColumns)

	wl = make([]string, len(cols))
	copy(wl, cols)

	return wl
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

// DeleteP deletes a single File record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *File) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows.
func (o fileQuery) DeleteAll() error {
	if o.Query == nil {
		return errors.New("models: no fileQuery provided for delete all")
	}

	boil.SetDelete(o.Query)

	_, err := boil.ExecQuery(o.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from files: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (o fileQuery) DeleteAllP() {
	if err := o.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
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

// DeleteAll deletes all rows in the slice with an executor.
func (o FileSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no File slice provided for delete all")
	}

	var mods []qm.QueryMod

	args := o.inPrimaryKeyArgs()
	in := boil.WherePrimaryKeyIn(len(o), "id")

	mods = append(mods,
		qm.From("files"),
		qm.Where(in, args...),
	)

	query := NewQuery(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from file slice: %s", err)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// DeleteAllP deletes all rows in the slice with an executor, and panics on error.
func (o FileSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
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

