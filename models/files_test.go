package models

import (
	"testing"
	"reflect"
	"time"

	"gopkg.in/nullbio/null.v4"
	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
	"github.com/nullbio/sqlboiler/strmangle"
)

func TestFiles(t *testing.T) {
	var err error

	o := make(FileSlice, 2)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	// insert two random objects to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Delete all rows to give a clean slate
	err = FilesG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Files: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	var c int64
	c, err = FilesG().Count()

	if c != 0 {
		t.Errorf("Expected files table to be empty, but got %d rows", c)
	}

	o = make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Ensure Count is valid
	c, err = FilesG().Count()
	if c != 3 {
		t.Errorf("Expected files table to have 3 rows, but got %d", c)
	}

	// Attempt to retrieve all objects
	res, err := FilesG().All()
	if err != nil {
		t.Errorf("Unable to retrieve all Files, err: %s", err)
	}

	if len(res) != 3 {
		t.Errorf("Expected 3 File rows, got %d", len(res))
	}

	filesDeleteAllRows(t)
}

func filesDeleteAllRows(t *testing.T) {
	// Delete all rows to give a clean slate
	err := FilesG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Files: %s", err)
	}
}

func TestFilesQueryDeleteAll(t *testing.T) {
	var err error
	var c int64

	// Start from a clean slate
	filesDeleteAllRows(t)

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = FilesG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after ObjDeleteAllRows() call, but got %d rows", c)
	}

	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Test DeleteAll() query function
	err = FilesG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Files: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = FilesG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after Obj().DeleteAll() call, but got %d rows", c)
	}
}

func TestFilesSliceDeleteAll(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test DeleteAll
	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	// test DeleteAll slice function
	if err = o.DeleteAllG(); err != nil {
		t.Errorf("Unable to objSlice.DeleteAll(): %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = FilesG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after objSlice.DeleteAll() call, but got %d rows", c)
	}
}

func TestFilesDelete(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test Delete
	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	o[0].DeleteG()

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = FilesG().Count()

	if c != 2 {
		t.Errorf("Expected 2 rows after obj.Delete() call, but got %d rows", c)
	}

	o[1].DeleteG()
	o[2].DeleteG()

	// Check number of rows in table to ensure Delete worked for all rows
	c, err = FilesG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after all obj.Delete() calls, but got %d rows", c)
	}
}

func TestFilesFind(t *testing.T) {
	var err error

	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(FileSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = FileFindG(o[i].ID)
		fileCompareVals(o[i], j[i], t)
	}

	f, err := FileFindG(o[0].ID, filePrimaryKeyColumns...)

	if o[0].ID != f.ID {
		t.Errorf("Expected primary key values to match, ID did not match")
	}

	colsWithoutPrimKeys := boil.SetComplement(fileColumns, filePrimaryKeyColumns)
	fRef := reflect.ValueOf(f).Elem()
	for _, v := range colsWithoutPrimKeys {
		val := fRef.FieldByName(v)
		if val.IsValid() {
			t.Errorf("Expected all other columns to be zero value, but column %s was %#v", v, val.Interface())
		}
	}

	filesDeleteAllRows(t)
}

func TestFilesBind(t *testing.T) {
	var err error

	o := File{}
	if err = boil.RandomizeStruct(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert File:\n%#v\nErr: %s", o, err)
	}

	j := File{}

	err = FilesG(qm.Where(`"id"=$1`, o.ID)).Bind(&j)
	if err != nil {
		t.Errorf("Unable to call Bind on File single object: %s", err)
	}

	fileCompareVals(&o, &j, t)

	// insert 3 rows, attempt to bind into slice
	filesDeleteAllRows(t)

	y := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&y, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(y); i++ {
		err = y[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", y[i], err)
		}
	}

	k := FileSlice{}
	err = FilesG().Bind(&k)
	if err != nil {
		t.Errorf("Unable to call Bind on File slice of objects: %s", err)
	}

	if len(k) != 3 {
		t.Errorf("Expected 3 results, got %d", len(k))
	}

	for i := 0; i < len(y); i++ {
		fileCompareVals(y[i], k[i], t)
	}

	filesDeleteAllRows(t)
}

func TestFilesOne(t *testing.T) {
	var err error

	o := File{}
	if err = boil.RandomizeStruct(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert File:\n%#v\nErr: %s", o, err)
	}

	j, err := FilesG().One()
	if err != nil {
		t.Errorf("Unable to fetch One File result:\n#%v\nErr: %s", j, err)
	}

	fileCompareVals(&o, j, t)

	filesDeleteAllRows(t)
}

func TestFilesAll(t *testing.T) {
	var err error

	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	j, err := FilesG().All()
	if err != nil {
		t.Errorf("Unable to fetch All File results: %s", err)
	}

	if len(j) != 3 {
		t.Errorf("Expected 3 results, got %d", len(j))
	}

	for i := 0; i < len(o); i++ {
		fileCompareVals(o[i], j[i], t)
	}

	filesDeleteAllRows(t)
}

func TestFilesCount(t *testing.T) {
	var err error

	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	// insert random columns to test Count
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	c, err := FilesG().Count()
	if err != nil {
		t.Errorf("Unable to count query File: %s", err)
	}

	if c != 3 {
		t.Errorf("Expected 3 results from count File, got %d", c)
	}

	filesDeleteAllRows(t)
}

var fileDBTypes = map[string]string{"Hash": "text", "Type": "text", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone", "Slug": "text", "ID": "uuid", "Name": "text", "State": "integer", "Size": "integer", "NumChunks": "integer"}

func fileCompareVals(o *File, j *File, t *testing.T) {
	if j.ID != o.ID {
		t.Errorf("Expected id columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.ID, j.ID)
	}

	if j.Size != o.Size {
		t.Errorf("Expected size columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Size, j.Size)
	}

	if j.NumChunks != o.NumChunks {
		t.Errorf("Expected num_chunks columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.NumChunks, j.NumChunks)
	}

	if j.State != o.State {
		t.Errorf("Expected state columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.State, j.State)
	}

	if j.Name != o.Name {
		t.Errorf("Expected name columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Name, j.Name)
	}

	if j.Hash != o.Hash {
		t.Errorf("Expected hash columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Hash, j.Hash)
	}

	if j.Type != o.Type {
		t.Errorf("Expected type columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Type, j.Type)
	}

	if o.CreatedAt.Format("02/01/2006") != j.CreatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time created_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.CreatedAt.Format("02/01/2006"), j.CreatedAt.Format("02/01/2006"))
	}

	if o.UpdatedAt.Format("02/01/2006") != j.UpdatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time updated_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.UpdatedAt.Format("02/01/2006"), j.UpdatedAt.Format("02/01/2006"))
	}

	if j.Slug != o.Slug {
		t.Errorf("Expected slug columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Slug, j.Slug)
	}
}

func TestFilesInPrimaryKeyArgs(t *testing.T) {
	var err error
	var o File
	o = File{}

	if err = boil.RandomizeStruct(&o, fileDBTypes, true); err != nil {
		t.Errorf("Could not randomize struct: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(filePrimaryKeyColumns) {
		t.Errorf("Expected args to be len %d, but got %d", len(filePrimaryKeyColumns), len(args))
	}

	if o.ID != args[0] {
		t.Errorf("Expected args[0] to be value of o.ID, but got %#v", args[0])
	}
}

func TestFilesSliceInPrimaryKeyArgs(t *testing.T) {
	var err error
	o := make(FileSlice, 3)

	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Could not randomize slice: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(filePrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(filePrimaryKeyColumns)*3, len(args))
	}

	for i := 0; i < len(filePrimaryKeyColumns)*3; i++ {

		if o[i].ID != args[i] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
	}
}

func fileBeforeCreateHook(o *File) error {
	*o = File{}
	return nil
}

func fileAfterCreateHook(o *File) error {
	*o = File{}
	return nil
}

func fileBeforeUpdateHook(o *File) error {
	*o = File{}
	return nil
}

func fileAfterUpdateHook(o *File) error {
	*o = File{}
	return nil
}

func TestFilesHooks(t *testing.T) {
	var err error

	empty := &File{}
	o := &File{}

	if err = boil.RandomizeStruct(o, fileDBTypes, false); err != nil {
		t.Errorf("Unable to randomize File object: %s", err)
	}

	FileAddHook(boil.HookBeforeCreate, fileBeforeCreateHook)
	if err = o.doBeforeCreateHooks(); err != nil {
		t.Errorf("Unable to execute doBeforeCreateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeCreateHook function to empty object, but got: %#v", o)
	}

	fileBeforeCreateHooks = []FileHook{}
	filesDeleteAllRows(t)
}

func TestFilesInsert(t *testing.T) {
	var err error

	var errs []error
	_ = errs

	emptyTime := time.Time{}.String()
	_ = emptyTime

	nullTime := null.NewTime(time.Time{}, true)
	_ = nullTime

	o := make(FileSlice, 3)
	if err = boil.RandomizeSlice(&o, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert File:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(FileSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = FileFindG(o[i].ID)
		fileCompareVals(o[i], j[i], t)
	}

	filesDeleteAllRows(t)

	item := &File{}
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item File:\n%#v\nErr: %s", item, err)
	}

	for _, c := range fileAutoIncrementColumns {
		// Ensure the auto increment columns are returned in the object.
		if errs = boil.IsZeroValue(item, false, c); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected auto-increment columns to be greater than 0, err: %s\n", e)
			}
		}
	}

	defaultValues := []interface{}{}

	// Ensure the simple default column values are returned correctly.
	if len(fileColumnsWithSimpleDefault) > 0 && len(defaultValues) > 0 {
		if len(fileColumnsWithSimpleDefault) != len(defaultValues) {
			t.Fatalf("Mismatch between slice lengths: %d, %d", len(fileColumnsWithSimpleDefault), len(defaultValues))
		}

		if errs = boil.IsValueMatch(item, fileColumnsWithSimpleDefault, defaultValues); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected default value to match column value, err: %s\n", e)
			}
		}
	}

	regularCols := []string{"size", "num_chunks", "state", "name", "hash", "type", "created_at", "updated_at"}

	// Ensure the non-defaultvalue columns and non-autoincrement columns are stored correctly as zero or null values.
	for _, c := range regularCols {
		rv := reflect.Indirect(reflect.ValueOf(item))
		field := rv.FieldByName(strmangle.TitleCase(c))

		zv := reflect.Zero(field.Type()).Interface()
		fv := field.Interface()

		if !reflect.DeepEqual(zv, fv) {
			t.Errorf("Expected column %s to be zero value, got: %v, wanted: %v", c, fv, zv)
		}
	}

	item = &File{}

	wl, rc := item.generateInsertColumns()
	if !reflect.DeepEqual(rc, fileColumnsWithDefault) {
		t.Errorf("Expected return columns to contain all columns with default values:\n\nGot: %v\nWanted: %v", rc, fileColumnsWithDefault)
	}

	if !reflect.DeepEqual(wl, fileColumnsWithoutDefault) {
		t.Errorf("Expected whitelist to contain all columns without default values:\n\nGot: %v\nWanted: %v", wl, fileColumnsWithoutDefault)
	}

	if err = boil.RandomizeStruct(item, fileDBTypes, false); err != nil {
		t.Errorf("Unable to randomize item: %s", err)
	}

	wl, rc = item.generateInsertColumns()
	if len(rc) > 0 {
		t.Errorf("Expected return columns to contain no columns:\n\nGot: %v", rc)
	}

	if !reflect.DeepEqual(wl, fileColumns) {
		t.Errorf("Expected whitelist to contain all columns values:\n\nGot: %v\nWanted: %v", wl, fileColumns)
	}

	filesDeleteAllRows(t)
}

func TestFileToManyChunks(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Chunk

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	boil.RandomizeStruct(&b, chunkDBTypes, true, "file_id")
	boil.RandomizeStruct(&c, chunkDBTypes, true, "file_id")

	b.FileID = a.ID
	c.FileID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	chunks, err := a.Chunks(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range chunks {
		if v.FileID == b.FileID {
			bFound = true
		}
		if v.FileID == c.FileID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	if t.Failed() {
		t.Logf("%#v", chunks)
	}
}

func TestFileToManyThumbnails(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Thumbnail

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	boil.RandomizeStruct(&b, thumbnailDBTypes, true, "file_id")
	boil.RandomizeStruct(&c, thumbnailDBTypes, true, "file_id")
	b.FileID.Valid = true
	c.FileID.Valid = true
	b.FileID.String = a.ID
	c.FileID.String = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	thumbnails, err := a.Thumbnails(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range thumbnails {
		if v.FileID.String == b.FileID.String {
			bFound = true
		}
		if v.FileID.String == c.FileID.String {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	if t.Failed() {
		t.Logf("%#v", thumbnails)
	}
}



func TestFilesSelect(t *testing.T) {
	// Only run this test if there are ample cols to test on
	if len(fileAutoIncrementColumns) == 0 {
		return
	}

	var err error

	x := &struct {
	}{}

	item := File{}

	blacklistCols := boil.SetMerge(fileAutoIncrementColumns, filePrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, fileDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert item File:\n%#v\nErr: %s", item, err)
	}

	err = FilesG(qm.Select(fileAutoIncrementColumns...), qm.Where(`"id"=$1`, item.ID)).Bind(x)
	if err != nil {
		t.Errorf("Unable to select insert results with bind: %s", err)
	}

	filesDeleteAllRows(t)
}

func TestFilesUpdate(t *testing.T) {
	var err error

	item := File{}
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item File:\n%#v\nErr: %s", item, err)
	}

	blacklistCols := boil.SetMerge(fileAutoIncrementColumns, filePrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, fileDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	whitelist := boil.SetComplement(fileColumns, fileAutoIncrementColumns)
	if err = item.UpdateG(whitelist...); err != nil {
		t.Errorf("Unable to update File: %s", err)
	}

	var j *File
	j, err = FileFindG(item.ID)
	if err != nil {
		t.Errorf("Unable to find File row: %s", err)
	}

	fileCompareVals(&item, j, t)

	wl := item.generateUpdateColumns("test")
	if len(wl) != 1 && wl[0] != "test" {
		t.Errorf("Expected generateUpdateColumns whitelist to match expected whitelist")
	}

	wl = item.generateUpdateColumns()
	if len(wl) == 0 && len(fileColumnsWithoutDefault) > 0 {
		t.Errorf("Expected generateUpdateColumns to build a whitelist for File, but got 0 results")
	}

	filesDeleteAllRows(t)
}

