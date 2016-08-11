package models

import (
	"testing"
	"reflect"
	"time"
	"errors"
	"fmt"

	"gopkg.in/vattle/null.v4"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/boil/qm"
	"github.com/vattle/sqlboiler/strmangle"
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

func TestThumbnailsExists(t *testing.T) {
	var err error

	o := Thumbnail{}
	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o, err)
	}

	// Check Exists finds existing rows
	e, err := ThumbnailExistsG(o.ID)
	if err != nil {
		t.Errorf("Unable to check if Thumbnail exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected ThumbnailExistsG to return true, but got false.")
	}

	whereClause := strmangle.WhereClause(1, thumbnailPrimaryKeyColumns)
	e, err = ThumbnailsG(qm.Where(whereClause, boil.GetStructValues(o, thumbnailPrimaryKeyColumns...)...)).Exists()
	if err != nil {
		t.Errorf("Unable to check if Thumbnail exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected ExistsG to return true, but got false.")
	}

	// Check Exists does not find non-existing rows
	o = Thumbnail{}
	e, err = ThumbnailExistsG(o.ID)
	if err != nil {
		t.Errorf("Unable to check if Thumbnail exists: %s", err)
	}
	if e != false {
		t.Errorf("Expected ThumbnailExistsG to return false, but got true.")
	}

	e, err = ThumbnailsG(qm.Where(whereClause, boil.GetStructValues(o, thumbnailPrimaryKeyColumns...)...)).Exists()
	if err != nil {
		t.Errorf("Unable to check if Thumbnail exists: %s", err)
	}
	if e != false {
		t.Errorf("Expected ExistsG to return false, but got true.")
	}

	thumbnailsDeleteAllRows(t)
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
		err = thumbnailCompareVals(o[i], j[i], true)
		if err != nil {
			t.Error(err)
		}
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

	err = thumbnailCompareVals(&o, &j, true)
	if err != nil {
		t.Error(err)
	}

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
		err = thumbnailCompareVals(y[i], k[i], true)
		if err != nil {
			t.Error(err)
		}
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

	err = thumbnailCompareVals(&o, j, true)
	if err != nil {
		t.Error(err)
	}

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
		err = thumbnailCompareVals(o[i], j[i], true)
		if err != nil {
			t.Error(err)
		}
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

var thumbnailDBTypes = map[string]string{"Hash": "text", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone", "ID": "uuid", "FileID": "uuid", "Size": "integer"}

func thumbnailCompareVals(o *Thumbnail, j *Thumbnail, equal bool, blacklist ...string) error {
	if ((equal && j.ID != o.ID) ||
		(!equal && j.ID == o.ID)) &&
		!strmangle.HasElement("id", blacklist) {
		return errors.New(fmt.Sprintf("Expected id columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.ID, j.ID))
	}

	if ((equal && j.FileID != o.FileID) ||
		(!equal && j.FileID == o.FileID)) &&
		!strmangle.HasElement("file_id", blacklist) {
		return errors.New(fmt.Sprintf("Expected file_id columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.FileID, j.FileID))
	}

	if ((equal && j.Size != o.Size) ||
		(!equal && j.Size == o.Size)) &&
		!strmangle.HasElement("size", blacklist) {
		return errors.New(fmt.Sprintf("Expected size columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Size, j.Size))
	}

	if ((equal && j.Hash != o.Hash) ||
		(!equal && j.Hash == o.Hash)) &&
		!strmangle.HasElement("hash", blacklist) {
		return errors.New(fmt.Sprintf("Expected hash columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Hash, j.Hash))
	}

	if ((equal && o.CreatedAt.Format("02/01/2006") != j.CreatedAt.Format("02/01/2006")) ||
		(!equal && o.CreatedAt.Format("02/01/2006") == j.CreatedAt.Format("02/01/2006"))) &&
		!strmangle.HasElement("created_at", blacklist) {
		return errors.New(fmt.Sprintf("Time created_at unexpected value, got:\nStruct: %#v\nResponse: %#v\n\n", o.CreatedAt.Format("02/01/2006"), j.CreatedAt.Format("02/01/2006")))
	}

	if ((equal && o.UpdatedAt.Format("02/01/2006") != j.UpdatedAt.Format("02/01/2006")) ||
		(!equal && o.UpdatedAt.Format("02/01/2006") == j.UpdatedAt.Format("02/01/2006"))) &&
		!strmangle.HasElement("updated_at", blacklist) {
		return errors.New(fmt.Sprintf("Time updated_at unexpected value, got:\nStruct: %#v\nResponse: %#v\n\n", o.UpdatedAt.Format("02/01/2006"), j.UpdatedAt.Format("02/01/2006")))
	}
	return nil
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
	for i := 0; i < len(o); i++ {
		j[i], err = ThumbnailFindG(o[i].ID)
		if err != nil {
			t.Errorf("Unable to find Thumbnail row: %s", err)
		}
		err = thumbnailCompareVals(o[i], j[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	thumbnailsDeleteAllRows(t)

	item := &Thumbnail{}
	boil.RandomizeValidatedStruct(item, thumbnailValidatedColumns, thumbnailDBTypes)
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

	// Remove the validated columns, they can never be zero values
	regularCols = boil.SetComplement(regularCols, thumbnailValidatedColumns)

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

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.FileID = foreign.ID
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


func TestThumbnailsReload(t *testing.T) {
	var err error

	o := Thumbnail{}
	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o, err)
	}

	// Create another copy of the object
	o1, err := ThumbnailFindG(o.ID)
	if err != nil {
		t.Errorf("Unable to find Thumbnail row.")
	}

	// Randomize the struct values again, except for the primary key values, so we can call update.
	err = boil.RandomizeStruct(&o, thumbnailDBTypes, true, thumbnailPrimaryKeyColumns...)
	if err != nil {
		t.Errorf("Unable to randomize Thumbnail struct members excluding primary keys: %s", err)
	}

	colsWithoutPrimKeys := boil.SetComplement(thumbnailColumns, thumbnailPrimaryKeyColumns)

	if err = o.UpdateG(colsWithoutPrimKeys...); err != nil {
		t.Errorf("Unable to update the Thumbnail row: %s", err)
	}

	if err = o1.ReloadG(); err != nil {
		t.Errorf("Unable to reload Thumbnail object: %s", err)
	}
	err = thumbnailCompareVals(&o, o1, true)
	if err != nil {
		t.Error(err)
	}

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsReloadAll(t *testing.T) {
	var err error

	o1 := make(ThumbnailSlice, 3)
	o2 := make(ThumbnailSlice, 3)
	if err = boil.RandomizeSlice(&o1, thumbnailDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o1); i++ {
		if err = o1[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o1[i], err)
		}
	}

	for i := 0; i < len(o1); i++ {
		o2[i], err = ThumbnailFindG(o1[i].ID)
		if err != nil {
			t.Errorf("Unable to find Thumbnail row.")
		}
		err = thumbnailCompareVals(o1[i], o2[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	// Randomize the struct values again, except for the primary key values, so we can call update.
	err = boil.RandomizeSlice(&o1, thumbnailDBTypes, false, thumbnailPrimaryKeyColumns...)
	if err != nil {
		t.Errorf("Unable to randomize Thumbnail slice excluding primary keys: %s", err)
	}

	colsWithoutPrimKeys := boil.SetComplement(thumbnailColumns, thumbnailPrimaryKeyColumns)

	for i := 0; i < len(o1); i++ {
		if err = o1[i].UpdateG(colsWithoutPrimKeys...); err != nil {
			t.Errorf("Unable to update the Thumbnail row: %s", err)
		}
	}

	if err = o2.ReloadAllG(); err != nil {
		t.Errorf("Unable to reload Thumbnail object: %s", err)
	}

	for i := 0; i < len(o1); i++ {
		err = thumbnailCompareVals(o2[i], o1[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	thumbnailsDeleteAllRows(t)
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
	boil.RandomizeValidatedStruct(&item, thumbnailValidatedColumns, thumbnailDBTypes)
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

	err = thumbnailCompareVals(&item, j, true)
	if err != nil {
		t.Error(err)
	}

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

func TestThumbnailsSliceUpdateAll(t *testing.T) {
	var err error

	// insert random columns to test UpdateAll
	o := make(ThumbnailSlice, 3)
	j := make(ThumbnailSlice, 3)

	if err = boil.RandomizeSlice(&o, thumbnailDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Thumbnail slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Thumbnail:\n%#v\nErr: %s", o[i], err)
		}
	}

	vals := M{}

	tmp := Thumbnail{}
	blacklist := boil.SetMerge(thumbnailPrimaryKeyColumns, thumbnailUniqueColumns)
	if err = boil.RandomizeStruct(&tmp, thumbnailDBTypes, false, blacklist...); err != nil {
		t.Errorf("Unable to randomize struct Thumbnail: %s", err)
	}

	// Build the columns and column values from the randomized struct
	tmpVal := reflect.Indirect(reflect.ValueOf(tmp))
	nonBlacklist := boil.SetComplement(thumbnailColumns, blacklist)
	for _, col := range nonBlacklist {
		vals[col] = tmpVal.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	err = o.UpdateAllG(vals)
	if err != nil {
		t.Errorf("Failed to update all for Thumbnail: %s", err)
	}

	for i := 0; i < len(o); i++ {
		j[i], err = ThumbnailFindG(o[i].ID)
		if err != nil {
			t.Errorf("Unable to find Thumbnail row: %s", err)
		}

		err = thumbnailCompareVals(j[i], &tmp, true, blacklist...)
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < len(o); i++ {
		// Ensure Find found the correct primary key ID's
		orig := boil.GetStructValues(o[i], thumbnailPrimaryKeyColumns...)
		new := boil.GetStructValues(j[i], thumbnailPrimaryKeyColumns...)

		if !reflect.DeepEqual(orig, new) {
			t.Errorf("object %d): primary keys do not match:\n\n%#v\n%#v", i, orig, new)
		}
	}

	thumbnailsDeleteAllRows(t)
}

func TestThumbnailsUpsert(t *testing.T) {
	var err error

	o := Thumbnail{}

	columns := o.generateUpsertColumns([]string{"one", "two"}, []string{"three", "four"}, []string{"five", "six"})
	if columns.conflict[0] != "one" || columns.conflict[1] != "two" {
		t.Errorf("Expected conflict to be %v, got %v", []string{"one", "two"}, columns.conflict)
	}

	if columns.update[0] != "three" || columns.update[1] != "four" {
		t.Errorf("Expected update to be %v, got %v", []string{"three", "four"}, columns.update)
	}

	if columns.whitelist[0] != "five" || columns.whitelist[1] != "six" {
		t.Errorf("Expected whitelist to be %v, got %v", []string{"five", "six"}, columns.whitelist)
	}

	columns = o.generateUpsertColumns(nil, nil, nil)
	if len(columns.whitelist) == 0 {
		t.Errorf("Expected whitelist to contain columns, but got len 0")
	}

	if len(columns.conflict) == 0 {
		t.Errorf("Expected conflict to contain columns, but got len 0")
	}

	if len(columns.update) == 0 {
		t.Errorf("expected update to contain columns, but got len 0")
	}

	upsertCols := upsertData{
		conflict:  []string{"key1", `"key2"`},
		update:    []string{"aaa", `"bbb"`},
		whitelist: []string{"thing", `"stuff"`},
		returning: []string{},
	}

	query := o.generateUpsertQuery(false, upsertCols)
	expectedQuery := `INSERT INTO thumbnails ("thing", "stuff") VALUES ($1, $2) ON CONFLICT DO NOTHING`

	if query != expectedQuery {
		t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
	}

	query = o.generateUpsertQuery(true, upsertCols)
	expectedQuery = `INSERT INTO thumbnails ("thing", "stuff") VALUES ($1, $2) ON CONFLICT ("key1", "key2") DO UPDATE SET "aaa" = EXCLUDED."aaa", "bbb" = EXCLUDED."bbb"`

	if query != expectedQuery {
		t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
	}

	upsertCols.returning = []string{"stuff"}
	query = o.generateUpsertQuery(true, upsertCols)
	expectedQuery = expectedQuery + ` RETURNING "stuff"`

	if query != expectedQuery {
		t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
	}

	// Attempt the INSERT side of an UPSERT
	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = o.UpsertG(false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Thumbnail: %s", err)
	}

	compare, err := ThumbnailFindG(o.ID)
	if err != nil {
		t.Errorf("Unable to find Thumbnail: %s", err)
	}
	err = thumbnailCompareVals(&o, compare, true)
	if err != nil {
		t.Error(err)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = boil.RandomizeStruct(&o, thumbnailDBTypes, false, thumbnailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = o.UpsertG(true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Thumbnail: %s", err)
	}

	compare, err = ThumbnailFindG(o.ID)
	if err != nil {
		t.Errorf("Unable to find Thumbnail: %s", err)
	}
	err = thumbnailCompareVals(&o, compare, true)
	if err != nil {
		t.Error(err)
	}

	thumbnailsDeleteAllRows(t)
}

