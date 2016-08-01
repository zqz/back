package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
	"gopkg.in/nullbio/null.v4"
)

// Thumbnail is an object representing the database table.
type Thumbnail struct {
	ID        string      `db:"thumbnail_id" json:"id"`
	FileID    null.String `db:"thumbnail_file_id" json:"file_id"`
	Size      null.Int32  `db:"thumbnail_size" json:"size"`
	Hash      null.String `db:"thumbnail_hash" json:"hash"`
	CreatedAt time.Time   `db:"thumbnail_created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"thumbnail_updated_at" json:"updated_at"`
}

var (
	thumbnailColumns                  = []string{"id", "file_id", "size", "hash", "created_at", "updated_at"}
	thumbnailColumnsWithoutDefault    = []string{"file_id", "size", "hash", "created_at", "updated_at"}
	thumbnailColumnsWithDefault       = []string{"id"}
	thumbnailColumnsWithSimpleDefault = []string{}
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
var thumbnailAfterCreateHooks []ThumbnailHook
var thumbnailAfterUpdateHooks []ThumbnailHook

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

func ThumbnailAddHook(hookPoint boil.HookPoint, thumbnailHook ThumbnailHook) {
	switch hookPoint {
	case boil.HookBeforeCreate:
		thumbnailBeforeCreateHooks = append(thumbnailBeforeCreateHooks, thumbnailHook)
	case boil.HookBeforeUpdate:
		thumbnailBeforeUpdateHooks = append(thumbnailBeforeUpdateHooks, thumbnailHook)
	case boil.HookAfterCreate:
		thumbnailAfterCreateHooks = append(thumbnailAfterCreateHooks, thumbnailHook)
	case boil.HookAfterUpdate:
		thumbnailAfterUpdateHooks = append(thumbnailAfterUpdateHooks, thumbnailHook)
	}
}

// One returns a single thumbnail record from the query.
func (q thumbnailQuery) One() (*Thumbnail, error) {
	o := &Thumbnail{}

	boil.SetLimit(q.Query, 1)

	res := boil.ExecQueryOne(q.Query)
	err := boil.BindOne(res, boil.Select(q.Query), o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute a one query for thumbnails: %s", err)
	}

	return o, nil
}

// OneP returns a single thumbnail record from the query, and panics on error.
func (q thumbnailQuery) OneP() *Thumbnail {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Thumbnail records from the query.
func (q thumbnailQuery) All() (ThumbnailSlice, error) {
	var o ThumbnailSlice

	res, err := boil.ExecQueryAll(q.Query)
	if err != nil {
		return nil, fmt.Errorf("models: failed to execute an all query for thumbnails: %s", err)
	}
	defer res.Close()

	err = boil.BindAll(res, boil.Select(q.Query), &o)
	if err != nil {
		return nil, fmt.Errorf("models: failed to assign all query results to Thumbnail slice: %s", err)
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

// CountP returns the count of all Thumbnail records in the query, and panics on error.
func (q thumbnailQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}


// File pointed to by the foreign key.
func (t *Thumbnail) File(selectCols ...string) (*File, error) {
	return t.FileX(boil.GetDB(), selectCols...)
}

// FileP pointed to by the foreign key. Panics on error.
func (t *Thumbnail) FileP(selectCols ...string) *File {
	o, err := t.FileX(boil.GetDB(), selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// FileXP pointed to by the foreign key with exeuctor. Panics on error.
func (t *Thumbnail) FileXP(exec boil.Executor, selectCols ...string) *File {
	o, err := t.FileX(exec, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// FileX pointed to by the foreign key.
func (t *Thumbnail) FileX(exec boil.Executor, selectCols ...string) (*File, error) {
	file := &File{}

	selectColumns := `*`
	if len(selectCols) != 0 {
		selectColumns = fmt.Sprintf(`"%s"`, strings.Join(selectCols, `","`))
	}

	query := fmt.Sprintf(`select %s from files where "id" = $1`, selectColumns)
	err := exec.QueryRow(query, t.FileID).Scan(boil.GetStructPointers(file, selectCols...)...)
	if err != nil {
		return nil, fmt.Errorf(`models: unable to select from files: %v`, err)
	}

	return file, nil
}



// ThumbnailsAll retrieves all records.
func Thumbnails(mods ...qm.QueryMod) thumbnailQuery {
	return ThumbnailsX(boil.GetDB(), mods...)
}

// ThumbnailsX retrieves all the records using an executor.
func ThumbnailsX(exec boil.Executor, mods ...qm.QueryMod) thumbnailQuery {
	mods = append(mods, qm.Table("thumbnails"))
	return thumbnailQuery{NewQueryX(exec, mods...)}
}


// ThumbnailFind retrieves a single record by ID.
func ThumbnailFind(id string, selectCols ...string) (*Thumbnail, error) {
	return ThumbnailFindX(boil.GetDB(), id, selectCols...)
}

// ThumbnailFindP retrieves a single record by ID, and panics on error.
func ThumbnailFindP(id string, selectCols ...string) *Thumbnail {
	o, err := ThumbnailFindX(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// ThumbnailFindX retrieves a single record by ID with an executor.
func ThumbnailFindX(exec boil.Executor, id string, selectCols ...string) (*Thumbnail, error) {
	thumbnail := &Thumbnail{}

	mods := []qm.QueryMod{
		qm.Select(selectCols...),
		qm.Table("thumbnails"),
		qm.Where(`"id"=$1`, id),
	}

	q := NewQueryX(exec, mods...)

	err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers(thumbnail, selectCols...)...)

	if err != nil {
		return nil, fmt.Errorf("models: unable to select from thumbnails: %v", err)
	}

	return thumbnail, nil
}

// ThumbnailFindXP retrieves a single record by ID with an executor, and panics on error.
func ThumbnailFindXP(exec boil.Executor, id string, selectCols ...string) *Thumbnail {
	o, err := ThumbnailFindX(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// Insert a single record.
func (o *Thumbnail) Insert(whitelist ...string) error {
	return o.InsertX(boil.GetDB(), whitelist...)
}

// InsertP a single record, and panics on error.
func (o *Thumbnail) InsertP(whitelist ...string) {
	if err := o.InsertX(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertX a single record using an executor.
func (o *Thumbnail) InsertX(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no thumbnails provided for insertion")
	}

	wl, returnColumns := o.generateInsertColumns(whitelist...)

	var err error
	if err := o.doBeforeCreateHooks(); err != nil {
		return err
	}

	ins := fmt.Sprintf(`INSERT INTO thumbnails ("%s") VALUES (%s)`, strings.Join(wl, `","`), boil.GenerateParamFlags(len(wl), 1))

	if len(returnColumns) != 0 {
		ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
		err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
	} else {
		_, err = exec.Exec(ins, o.ID, o.FileID, o.Size, o.Hash, o.CreatedAt, o.UpdatedAt)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, wl...))
	}

	if err != nil {
		return fmt.Errorf("models: unable to insert into thumbnails: %s", err)
	}

	if err := o.doAfterCreateHooks(); err != nil {
		return err
	}

	return nil
}

// InsertXP a single record using an executor, and panics on error.
func (o *Thumbnail) InsertXP(exec boil.Executor, whitelist ...string) {
	if err := o.InsertX(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
func (o *Thumbnail) generateInsertColumns(whitelist ...string) ([]string, []string) {
	var wl []string

	wl = append(wl, whitelist...)
	if len(whitelist) == 0 {
		wl = append(wl, thumbnailColumnsWithoutDefault...)
	}

	wl = append(boil.NonZeroDefaultSet(thumbnailColumnsWithDefault, o), wl...)
	wl = boil.SortByKeys(thumbnailColumns, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := boil.SetComplement(thumbnailColumnsWithDefault, wl)

	return wl, rc
}


// Update a single Thumbnail record.
// Update takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
func (o *Thumbnail) Update(whitelist ...string) error {
	return o.UpdateX(boil.GetDB(), whitelist...)
}

// Update a single Thumbnail record.
// UpdateP takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
// Panics on error.
func (o *Thumbnail) UpdateP(whitelist ...string) {
	if err := o.UpdateX(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateX uses an executor to update the Thumbnail.
func (o *Thumbnail) UpdateX(exec boil.Executor, whitelist ...string) error {
	return o.UpdateAtX(exec, o.ID, whitelist...)
}

// UpdateXP uses an executor to update the Thumbnail, and panics on error.
func (o *Thumbnail) UpdateXP(exec boil.Executor, whitelist ...string) {
	err := o.UpdateAtX(exec, o.ID, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAt updates the Thumbnail using the primary key to find the row to update.
func (o *Thumbnail) UpdateAt(id string, whitelist ...string) error {
	return o.UpdateAtX(boil.GetDB(), id, whitelist...)
}

// UpdateAtP updates the Thumbnail using the primary key to find the row to update. Panics on error.
func (o *Thumbnail) UpdateAtP(id string, whitelist ...string) {
	if err := o.UpdateAtX(boil.GetDB(), id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAtX uses an executor to update the Thumbnail using the primary key to find the row to update.
func (o *Thumbnail) UpdateAtX(exec boil.Executor, id string, whitelist ...string) error {
	if err := o.doBeforeUpdateHooks(); err != nil {
		return err
	}

	var err error
	var query string
	var values []interface{}

	wl := o.generateUpdateColumns(whitelist...)

	if len(wl) != 0 {
		query = fmt.Sprintf(`UPDATE thumbnails SET %s WHERE %s`, boil.SetParamNames(wl), boil.WherePrimaryKey(len(wl)+1, "id"))
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

// UpdateAtXP uses an executor to update the Thumbnail using the primary key to find the row to update.
// Panics on error.
func (o *Thumbnail) UpdateAtXP(exec boil.Executor, id string, whitelist ...string) {
	if err := o.UpdateAtX(exec, id, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with matching column names.
func (q thumbnailQuery) UpdateAll(cols M) error {
	boil.SetUpdate(q.Query, cols)

	_, err := boil.ExecQuery(q.Query)
	if err != nil {
		return fmt.Errorf("models: unable to update all for thumbnails: %s", err)
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q thumbnailQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// generateUpdateColumns generates the whitelist columns for an update statement
func (o *Thumbnail) generateUpdateColumns(whitelist ...string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	var wl []string
	cols := thumbnailColumnsWithoutDefault
	cols = append(boil.NonZeroDefaultSet(thumbnailColumnsWithDefault, o), cols...)
	// Subtract primary keys and autoincrement columns
	cols = boil.SetComplement(cols, thumbnailPrimaryKeyColumns)
	cols = boil.SetComplement(cols, thumbnailAutoIncrementColumns)

	wl = make([]string, len(cols))
	copy(wl, cols)

	return wl
}

// Delete deletes a single Thumbnail record.
// Delete will match against the primary key column to find the record to delete.
func (o *Thumbnail) Delete() error {
	if o == nil {
		return errors.New("models: no Thumbnail provided for deletion")
	}

	return o.DeleteX(boil.GetDB())
}

// DeleteP deletes a single Thumbnail record.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Thumbnail) DeleteP() {
	if err := o.Delete(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteX deletes a single Thumbnail record with an executor.
// DeleteX will match against the primary key column to find the record to delete.
func (o *Thumbnail) DeleteX(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Thumbnail provided for deletion")
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.Table("thumbnails"),
		qm.Where(`"id"=$1`, o.ID),
	)

	query := NewQueryX(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete from thumbnails: %s", err)
	}

	return nil
}

// DeleteXP deletes a single Thumbnail record with an executor.
// DeleteXP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Thumbnail) DeleteXP(exec boil.Executor) {
	if err := o.DeleteX(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows.
func (o thumbnailQuery) DeleteAll() error {
	if o.Query == nil {
		return errors.New("models: no thumbnailQuery provided for delete all")
	}

	boil.SetDelete(o.Query)

	_, err := boil.ExecQuery(o.Query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from thumbnails: %s", err)
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (o thumbnailQuery) DeleteAllP() {
	if err := o.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice.
func (o ThumbnailSlice) DeleteAll() error {
	if o == nil {
		return errors.New("models: no Thumbnail slice provided for delete all")
	}
	return o.DeleteAllX(boil.GetDB())
}

// DeleteAll deletes all rows in the slice.
func (o ThumbnailSlice) DeleteAllP() {
	if err := o.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllX deletes all rows in the slice with an executor.
func (o ThumbnailSlice) DeleteAllX(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Thumbnail slice provided for delete all")
	}

	var mods []qm.QueryMod

	args := o.inPrimaryKeyArgs()
	in := boil.WherePrimaryKeyIn(len(o), "id")

	mods = append(mods,
		qm.Table("thumbnails"),
		qm.Where(in, args...),
	)

	query := NewQueryX(exec, mods...)
	boil.SetDelete(query)

	_, err := boil.ExecQuery(query)
	if err != nil {
		return fmt.Errorf("models: unable to delete all from thumbnail slice: %s", err)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, args)
	}

	return nil
}

// DeleteAllXP deletes all rows in the slice with an executor, and panics on error.
func (o ThumbnailSlice) DeleteAllXP(exec boil.Executor) {
	if err := o.DeleteAllX(exec); err != nil {
		panic(boil.WrapErr(err))
	}
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

