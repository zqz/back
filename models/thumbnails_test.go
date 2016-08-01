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

func TestThumbnails(t *testing.T) {
	var err error

	o := make(ThumbnailSlice, 2)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	// insert two random objects to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Delete all rows to give a clean slate
	err = ThumbnailsG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Thumbnails: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	var c int64
	c, err = ThumbnailsG().Count()

	if c != 0 {
		t.Errorf("Expected thumbnails table to be empty, but got %d rows", c)
	}

	o = make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Ensure Count is valid
	c, err = ThumbnailsG().Count()
	if c != 3 {
		t.Errorf("Expected thumbnails table to have 3 rows, but got %d", c)
	}

	// Attempt to retrieve all objects
	res, err := ThumbnailsG().All()
	if err != nil {
		t.Errorf("Unable to retrieve all Thumbnails, err: %s", err)
	}

	if len(res) != 3 {
		t.Errorf("Expected 3 Thumbnail rows, got %d", len(res))
	}

	thumbnailsDeleteAllRows(t)
}

func thumbnailsDeleteAllRows(t *testing.T) {
	// Delete all rows to give a clean slate
	err := ThumbnailsG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Thumbnails: %s", err)
	}
}

func TestThumbnailsQueryDeleteAll(t *testing.T) {
	var err error
	var c int64

	// Start from a clean slate
	thumbnailsDeleteAllRows(t)

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ThumbnailsG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after ObjDeleteAllRows() call, but got %d rows", c)
	}

	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Test DeleteAll() query function
	err = ThumbnailsG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Thumbnails: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ThumbnailsG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after Obj().DeleteAll() call, but got %d rows", c)
	}
}

func TestThumbnailsSliceDeleteAll(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test DeleteAll
	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	// test DeleteAll slice function
	if err = o.DeleteAllG(); err != nil {
		t.Errorf("Unable to objSlice.DeleteAll(): %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ThumbnailsG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after objSlice.DeleteAll() call, but got %d rows", c)
	}
}

func TestThumbnailsDelete(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test Delete
	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	o[0].DeleteG()

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ThumbnailsG().Count()

	if c != 2 {
		t.Errorf("Expected 2 rows after obj.Delete() call, but got %d rows", c)
	}

	o[1].DeleteG()
	o[2].DeleteG()

	// Check number of rows in table to ensure Delete worked for all rows
	c, err = ThumbnailsG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after all obj.Delete() calls, but got %d rows", c)
	}
}

func TestThumbnailsFind(t *testing.T) {
	var err error

	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(ThumbnailSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = ThumbnailFindG(o[i].ID)
		thumbnailCompareVals(o[i], j[i], t)
	}

	f, err := ThumbnailFindG(o[0].ID, thumbnailPrimaryKeyColumns...)

	if o[0].ID != f.ID {
		t.Errorf("Expected primary key values to match, ID did not match")
	}

	colsWithoutPrimKeys := boil.SetComplement(thumbnailColumns, thumbnailPrimaryKeyColumns)
	fRef := reflect.ValueOf(f).Elem()
	for _, v := range colsWithoutPrimKeys {
		val := fRef.FieldByName(v)
		if val.IsValid() {
			t.Errorf("Expected all other columns to be zero value, but column %s was %#v", v, val.Interface())
		}
	}

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsBind(t *testing.T) {
	var err error

	o := Thumbnail{}
	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o, err)
	}

	j := Thumbnail{}

	err = ThumbnailsG(qm.Where(`"id"=$1`, o.ID)).Bind(&j)
	if err != nil {
		t.Errorf("Unable to call Bind on Thumbnail single object: %s", err)
	}

	thumbnailCompareVals(&o, &j, t)

	// insert 3 rows, attempt to bind into slice
	thumbnailsDeleteAllRows(t)

	y := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&y, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(y); i++ {
		err = y[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", y[i], err)
		}
	}

	k := ThumbnailSlice{}
	err = ThumbnailsG().Bind(&k)
	if err != nil {
		t.Errorf("Unable to call Bind on Thumbnail slice of objects: %s", err)
	}

	if len(k) != 3 {
		t.Errorf("Expected 3 results, got %d", len(k))
	}

	for i := 0; i < len(y); i++ {
		thumbnailCompareVals(y[i], k[i], t)
	}

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsOne(t *testing.T) {
	var err error

	o := Thumbnail{}
	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o, err)
	}

	j, err := ThumbnailsG().One()
	if err != nil {
		t.Errorf("Unable to fetch One Thumbnail result:\n#%v\nErr: %s", j, err)
	}

	thumbnailCompareVals(&o, j, t)

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsAll(t *testing.T) {
	var err error

	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	j, err := ThumbnailsG().All()
	if err != nil {
		t.Errorf("Unable to fetch All Thumbnail results: %s", err)
	}

	if len(j) != 3 {
		t.Errorf("Expected 3 results, got %d", len(j))
	}

	for i := 0; i < len(o); i++ {
		thumbnailCompareVals(o[i], j[i], t)
	}

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsCount(t *testing.T) {
	var err error

	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	// insert random columns to test Count
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	c, err := ThumbnailsG().Count()
	if err != nil {
		t.Errorf("Unable to count query Thumbnail: %s", err)
	}

	if c != 3 {
		t.Errorf("Expected 3 results from count Thumbnail, got %d", c)
	}

	thumbnailsDeleteAllRows(t)
}

var thumbnailDBTypes = map[string]string{"ID": "uuid", "FileID": "uuid", "Size": "integer", "Hash": "text", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone"}

func thumbnailCompareVals(o *Thumbnail, j *Thumbnail, t *testing.T) {
	if j.ID != o.ID {
		t.Errorf("Expected id columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.ID, j.ID)
	}

	if j.FileID != o.FileID {
		t.Errorf("Expected file_id columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.FileID, j.FileID)
	}

	if j.Size != o.Size {
		t.Errorf("Expected size columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Size, j.Size)
	}

	if j.Hash != o.Hash {
		t.Errorf("Expected hash columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Hash, j.Hash)
	}

	if o.CreatedAt.Format("02/01/2006") != j.CreatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time created_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.CreatedAt.Format("02/01/2006"), j.CreatedAt.Format("02/01/2006"))
	}

	if o.UpdatedAt.Format("02/01/2006") != j.UpdatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time updated_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.UpdatedAt.Format("02/01/2006"), j.UpdatedAt.Format("02/01/2006"))
	}
}

func TestThumbnailsInPrimaryKeyArgs(t *testing.T) {
	var err error
	var o Thumbnail
	o = Thumbnail{}

	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Could not randomize struct: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(thumbnailPrimaryKeyColumns) {
		t.Errorf("Expected args to be len %d, but got %d", len(thumbnailPrimaryKeyColumns), len(args))
	}

	if o.ID != args[0] {
		t.Errorf("Expected args[0] to be value of o.ID, but got %#v", args[0])
	}
}

func TestThumbnailsSliceInPrimaryKeyArgs(t *testing.T) {
	var err error
	o := make(ThumbnailSlice, 3)

	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Could not randomize slice: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(thumbnailPrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(thumbnailPrimaryKeyColumns)*3, len(args))
	}

	for i := 0; i < len(thumbnailPrimaryKeyColumns)*3; i++ {

		if o[i].ID != args[i] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
	}
}

func thumbnailBeforeCreateHook(o *Thumbnail) error {
	*o = Thumbnail{}
	return nil
}

func thumbnailAfterCreateHook(o *Thumbnail) error {
	*o = Thumbnail{}
	return nil
}

func thumbnailBeforeUpdateHook(o *Thumbnail) error {
	*o = Thumbnail{}
	return nil
}

func thumbnailAfterUpdateHook(o *Thumbnail) error {
	*o = Thumbnail{}
	return nil
}

func TestThumbnailsHooks(t *testing.T) {
	var err error

	empty := &Thumbnail{}
	o := &Thumbnail{}

	if err = boil.RandomizeStruct(o, thumbnailDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Thumbnail object: %s", err)
	}

	ThumbnailAddHook(boil.HookBeforeCreate, thumbnailBeforeCreateHook)
	if err = o.doBeforeCreateHooks(); err != nil {
		t.Errorf("Unable to execute doBeforeCreateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeCreateHook function to empty object, but got: %#v", o)
	}

	thumbnailBeforeCreateHooks = []ThumbnailHook{}
	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsInsert(t *testing.T) {
	var err error

	var errs []error
	_ = errs

	emptyTime := time.Time{}.String()
	_ = emptyTime

	nullTime := null.NewTime(time.Time{}, true)
	_ = nullTime

	o := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(ThumbnailSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = ThumbnailFindG(o[i].ID)
		thumbnailCompareVals(o[i], j[i], t)
	}

	thumbnailsDeleteAllRows(t)

	item := &Thumbnail{}
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item Thumbnail:\n%#v\nErr: %s", item, err)
	}

	for _, c := range thumbnailAutoIncrementColumns {
		// Ensure the auto increment columns are returned in the object.
		if errs = boil.IsZeroValue(item, false, c); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected auto-increment columns to be greater than 0, err: %s\n", e)
			}
		}
	}

	defaultValues := []interface{}{}

	// Ensure the simple default column values are returned correctly.
	if len(thumbnailColumnsWithSimpleDefault) > 0 && len(defaultValues) > 0 {
		if len(thumbnailColumnsWithSimpleDefault) != len(defaultValues) {
			t.Fatalf("Mismatch between slice lengths: %d, %d", len(thumbnailColumnsWithSimpleDefault), len(defaultValues))
		}

		if errs = boil.IsValueMatch(item, thumbnailColumnsWithSimpleDefault, defaultValues); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected default value to match column value, err: %s\n", e)
			}
		}
	}

	regularCols := []string{"file_id", "size", "hash", "created_at", "updated_at"}

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

	item = &Thumbnail{}

	wl, rc := item.generateInsertColumns()
	if !reflect.DeepEqual(rc, thumbnailColumnsWithDefault) {
		t.Errorf("Expected return columns to contain all columns with default values:\n\nGot: %v\nWanted: %v", rc, thumbnailColumnsWithDefault)
	}

	if !reflect.DeepEqual(wl, thumbnailColumnsWithoutDefault) {
		t.Errorf("Expected whitelist to contain all columns without default values:\n\nGot: %v\nWanted: %v", wl, thumbnailColumnsWithoutDefault)
	}

	if err = boil.RandomizeStruct(item, thumbnailDBTypes, false); err != nil {
		t.Errorf("Unable to randomize item: %s", err)
	}

	wl, rc = item.generateInsertColumns()
	if len(rc) > 0 {
		t.Errorf("Expected return columns to contain no columns:\n\nGot: %v", rc)
	}

	if !reflect.DeepEqual(wl, thumbnailColumns) {
		t.Errorf("Expected whitelist to contain all columns values:\n\nGot: %v\nWanted: %v", wl, thumbnailColumns)
	}

	thumbnailsDeleteAllRows(t)
}



func TestThumbnailToOneFile_File(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var foreign File
	var local Thumbnail
	local.FileID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.FileID.String = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}
	check, err := local.File(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}
}


func TestThumbnailsSelect(t *testing.T) {
	// Only run this test if there are ample cols to test on
	if len(thumbnailAutoIncrementColumns) == 0 {
		return
	}

	var err error

	x := &struct {
	}{}

	item := Thumbnail{}

	blacklistCols := boil.SetMerge(thumbnailAutoIncrementColumns, thumbnailPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, thumbnailDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert item Thumbnail:\n%#v\nErr: %s", item, err)
	}

	err = ThumbnailsG(qm.Select(thumbnailAutoIncrementColumns...), qm.Where(`"id"=$1`, item.ID)).Bind(x)
	if err != nil {
		t.Errorf("Unable to select insert results with bind: %s", err)
	}

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsUpdate(t *testing.T) {
	var err error

	item := Thumbnail{}
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item Thumbnail:\n%#v\nErr: %s", item, err)
	}

	blacklistCols := boil.SetMerge(thumbnailAutoIncrementColumns, thumbnailPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, thumbnailDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	whitelist := boil.SetComplement(thumbnailColumns, thumbnailAutoIncrementColumns)
	if err = item.UpdateG(whitelist...); err != nil {
		t.Errorf("Unable to update Thumbnail: %s", err)
	}

	var j *Thumbnail
	j, err = ThumbnailFindG(item.ID)
	if err != nil {
		t.Errorf("Unable to find Thumbnail row: %s", err)
	}

	thumbnailCompareVals(&item, j, t)

	wl := item.generateUpdateColumns("test")
	if len(wl) != 1 && wl[0] != "test" {
		t.Errorf("Expected generateUpdateColumns whitelist to match expected whitelist")
	}

	wl = item.generateUpdateColumns()
	if len(wl) == 0 && len(thumbnailColumnsWithoutDefault) > 0 {
		t.Errorf("Expected generateUpdateColumns to build a whitelist for Thumbnail, but got 0 results")
	}

	thumbnailsDeleteAllRows(t)
}

