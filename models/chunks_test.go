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

func TestChunks(t *testing.T) {
	var err error

	o := make(ChunkSlice, 2)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert two random objects to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Delete all rows to give a clean slate
	err = ChunksG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Chunks: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	var c int64
	c, err = ChunksG().Count()

	if c != 0 {
		t.Errorf("Expected chunks table to be empty, but got %d rows", c)
	}

	o = make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Ensure Count is valid
	c, err = ChunksG().Count()
	if c != 3 {
		t.Errorf("Expected chunks table to have 3 rows, but got %d", c)
	}

	// Attempt to retrieve all objects
	res, err := ChunksG().All()
	if err != nil {
		t.Errorf("Unable to retrieve all Chunks, err: %s", err)
	}

	if len(res) != 3 {
		t.Errorf("Expected 3 Chunk rows, got %d", len(res))
	}

	chunksDeleteAllRows(t)
}

func chunksDeleteAllRows(t *testing.T) {
	// Delete all rows to give a clean slate
	err := ChunksG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Chunks: %s", err)
	}
}

func TestChunksQueryDeleteAll(t *testing.T) {
	var err error
	var c int64

	// Start from a clean slate
	chunksDeleteAllRows(t)

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ChunksG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after ObjDeleteAllRows() call, but got %d rows", c)
	}

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Test DeleteAll() query function
	err = ChunksG().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Chunks: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ChunksG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after Obj().DeleteAll() call, but got %d rows", c)
	}
}

func TestChunksSliceDeleteAll(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test DeleteAll
	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// test DeleteAll slice function
	if err = o.DeleteAllG(); err != nil {
		t.Errorf("Unable to objSlice.DeleteAll(): %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ChunksG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after objSlice.DeleteAll() call, but got %d rows", c)
	}
}

func TestChunksDelete(t *testing.T) {
	var err error
	var c int64

	// insert random columns to test Delete
	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	o[0].DeleteG()

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = ChunksG().Count()

	if c != 2 {
		t.Errorf("Expected 2 rows after obj.Delete() call, but got %d rows", c)
	}

	o[1].DeleteG()
	o[2].DeleteG()

	// Check number of rows in table to ensure Delete worked for all rows
	c, err = ChunksG().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after all obj.Delete() calls, but got %d rows", c)
	}
}

func TestChunksExists(t *testing.T) {
	var err error

	o := Chunk{}
	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o, err)
	}

	// Check Exists finds existing rows
	e, err := ChunkExistsG(o.ID)
	if err != nil {
		t.Errorf("Unable to check if Chunk exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected ChunkExistsG to return true, but got false.")
	}

	whereClause := strmangle.WhereClause(1, chunkPrimaryKeyColumns)
	e, err = ChunksG(qm.Where(whereClause, boil.GetStructValues(o, chunkPrimaryKeyColumns...)...)).Exists()
	if err != nil {
		t.Errorf("Unable to check if Chunk exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected ExistsG to return true, but got false.")
	}

	// Check Exists does not find non-existing rows
	o = Chunk{}
	e, err = ChunkExistsG(o.ID)
	if err != nil {
		t.Errorf("Unable to check if Chunk exists: %s", err)
	}
	if e != false {
		t.Errorf("Expected ChunkExistsG to return false, but got true.")
	}

	e, err = ChunksG(qm.Where(whereClause, boil.GetStructValues(o, chunkPrimaryKeyColumns...)...)).Exists()
	if err != nil {
		t.Errorf("Unable to check if Chunk exists: %s", err)
	}
	if e != false {
		t.Errorf("Expected ExistsG to return false, but got true.")
	}

	chunksDeleteAllRows(t)
}

func TestChunksFind(t *testing.T) {
	var err error

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(ChunkSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = ChunkFindG(o[i].ID)
		err = chunkCompareVals(o[i], j[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	f, err := ChunkFindG(o[0].ID, chunkPrimaryKeyColumns...)

	if o[0].ID != f.ID {
		t.Errorf("Expected primary key values to match, ID did not match")
	}

	colsWithoutPrimKeys := boil.SetComplement(chunkColumns, chunkPrimaryKeyColumns)
	fRef := reflect.ValueOf(f).Elem()
	for _, v := range colsWithoutPrimKeys {
		val := fRef.FieldByName(v)
		if val.IsValid() {
			t.Errorf("Expected all other columns to be zero value, but column %s was %#v", v, val.Interface())
		}
	}

	chunksDeleteAllRows(t)
}

func TestChunksBind(t *testing.T) {
	var err error

	o := Chunk{}
	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o, err)
	}

	j := Chunk{}

	err = ChunksG(qm.Where(`"id"=$1`, o.ID)).Bind(&j)
	if err != nil {
		t.Errorf("Unable to call Bind on Chunk single object: %s", err)
	}

	err = chunkCompareVals(&o, &j, true)
	if err != nil {
		t.Error(err)
	}

	// insert 3 rows, attempt to bind into slice
	chunksDeleteAllRows(t)

	y := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&y, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(y); i++ {
		err = y[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", y[i], err)
		}
	}

	k := ChunkSlice{}
	err = ChunksG().Bind(&k)
	if err != nil {
		t.Errorf("Unable to call Bind on Chunk slice of objects: %s", err)
	}

	if len(k) != 3 {
		t.Errorf("Expected 3 results, got %d", len(k))
	}

	for i := 0; i < len(y); i++ {
		err = chunkCompareVals(y[i], k[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	chunksDeleteAllRows(t)
}

func TestChunksOne(t *testing.T) {
	var err error

	o := Chunk{}
	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o, err)
	}

	j, err := ChunksG().One()
	if err != nil {
		t.Errorf("Unable to fetch One Chunk result:\n#%v\nErr: %s", j, err)
	}

	err = chunkCompareVals(&o, j, true)
	if err != nil {
		t.Error(err)
	}

	chunksDeleteAllRows(t)
}

func TestChunksAll(t *testing.T) {
	var err error

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	j, err := ChunksG().All()
	if err != nil {
		t.Errorf("Unable to fetch All Chunk results: %s", err)
	}

	if len(j) != 3 {
		t.Errorf("Expected 3 results, got %d", len(j))
	}

	for i := 0; i < len(o); i++ {
		err = chunkCompareVals(o[i], j[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	chunksDeleteAllRows(t)
}

func TestChunksCount(t *testing.T) {
	var err error

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert random columns to test Count
	for i := 0; i < len(o); i++ {
		err = o[i].InsertG()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	c, err := ChunksG().Count()
	if err != nil {
		t.Errorf("Unable to count query Chunk: %s", err)
	}

	if c != 3 {
		t.Errorf("Expected 3 results from count Chunk, got %d", c)
	}

	chunksDeleteAllRows(t)
}

var chunkDBTypes = map[string]string{"Hash": "text", "Position": "integer", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone", "ID": "uuid", "FileID": "uuid", "Size": "integer"}

func chunkCompareVals(o *Chunk, j *Chunk, equal bool, blacklist ...string) error {
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

	if ((equal && j.Position != o.Position) ||
		(!equal && j.Position == o.Position)) &&
		!strmangle.HasElement("position", blacklist) {
		return errors.New(fmt.Sprintf("Expected position columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Position, j.Position))
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

func TestChunksInPrimaryKeyArgs(t *testing.T) {
	var err error
	var o Chunk
	o = Chunk{}

	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Could not randomize struct: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(chunkPrimaryKeyColumns) {
		t.Errorf("Expected args to be len %d, but got %d", len(chunkPrimaryKeyColumns), len(args))
	}

	if o.ID != args[0] {
		t.Errorf("Expected args[0] to be value of o.ID, but got %#v", args[0])
	}
}

func TestChunksSliceInPrimaryKeyArgs(t *testing.T) {
	var err error
	o := make(ChunkSlice, 3)

	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Could not randomize slice: %s", err)
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(chunkPrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(chunkPrimaryKeyColumns)*3, len(args))
	}

	for i := 0; i < len(chunkPrimaryKeyColumns)*3; i++ {

		if o[i].ID != args[i] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
	}
}

func chunkBeforeCreateHook(o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterCreateHook(o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkBeforeUpdateHook(o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterUpdateHook(o *Chunk) error {
	*o = Chunk{}
	return nil
}

func TestChunksHooks(t *testing.T) {
	var err error

	empty := &Chunk{}
	o := &Chunk{}

	if err = boil.RandomizeStruct(o, chunkDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Chunk object: %s", err)
	}

	ChunkAddHook(boil.HookBeforeCreate, chunkBeforeCreateHook)
	if err = o.doBeforeCreateHooks(); err != nil {
		t.Errorf("Unable to execute doBeforeCreateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeCreateHook function to empty object, but got: %#v", o)
	}

	chunkBeforeCreateHooks = []ChunkHook{}
	chunksDeleteAllRows(t)
}

func TestChunksInsert(t *testing.T) {
	var err error

	var errs []error
	_ = errs

	emptyTime := time.Time{}.String()
	_ = emptyTime

	nullTime := null.NewTime(time.Time{}, true)
	_ = nullTime

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(ChunkSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(o); i++ {
		j[i], err = ChunkFindG(o[i].ID)
		if err != nil {
			t.Errorf("Unable to find Chunk row: %s", err)
		}
		err = chunkCompareVals(o[i], j[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	chunksDeleteAllRows(t)

	item := &Chunk{}
	boil.RandomizeValidatedStruct(item, chunkValidatedColumns, chunkDBTypes)
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item Chunk:\n%#v\nErr: %s", item, err)
	}

	for _, c := range chunkAutoIncrementColumns {
		// Ensure the auto increment columns are returned in the object.
		if errs = boil.IsZeroValue(item, false, c); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected auto-increment columns to be greater than 0, err: %s\n", e)
			}
		}
	}

	defaultValues := []interface{}{}

	// Ensure the simple default column values are returned correctly.
	if len(chunkColumnsWithSimpleDefault) > 0 && len(defaultValues) > 0 {
		if len(chunkColumnsWithSimpleDefault) != len(defaultValues) {
			t.Fatalf("Mismatch between slice lengths: %d, %d", len(chunkColumnsWithSimpleDefault), len(defaultValues))
		}

		if errs = boil.IsValueMatch(item, chunkColumnsWithSimpleDefault, defaultValues); errs != nil {
			for _, e := range errs {
				t.Errorf("Expected default value to match column value, err: %s\n", e)
			}
		}
	}

	regularCols := []string{"file_id", "size", "hash", "position", "created_at", "updated_at"}

	// Remove the validated columns, they can never be zero values
	regularCols = boil.SetComplement(regularCols, chunkValidatedColumns)

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

	item = &Chunk{}

	wl, rc := item.generateInsertColumns()
	if !reflect.DeepEqual(rc, chunkColumnsWithDefault) {
		t.Errorf("Expected return columns to contain all columns with default values:\n\nGot: %v\nWanted: %v", rc, chunkColumnsWithDefault)
	}

	if !reflect.DeepEqual(wl, chunkColumnsWithoutDefault) {
		t.Errorf("Expected whitelist to contain all columns without default values:\n\nGot: %v\nWanted: %v", wl, chunkColumnsWithoutDefault)
	}

	if err = boil.RandomizeStruct(item, chunkDBTypes, false); err != nil {
		t.Errorf("Unable to randomize item: %s", err)
	}

	wl, rc = item.generateInsertColumns()
	if len(rc) > 0 {
		t.Errorf("Expected return columns to contain no columns:\n\nGot: %v", rc)
	}

	if !reflect.DeepEqual(wl, chunkColumns) {
		t.Errorf("Expected whitelist to contain all columns values:\n\nGot: %v\nWanted: %v", wl, chunkColumns)
	}

	chunksDeleteAllRows(t)
}



func TestChunkToOneFile_File(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var foreign File
	var local Chunk

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


func TestChunksReload(t *testing.T) {
	var err error

	o := Chunk{}
	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.InsertG(); err != nil {
		t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o, err)
	}

	// Create another copy of the object
	o1, err := ChunkFindG(o.ID)
	if err != nil {
		t.Errorf("Unable to find Chunk row.")
	}

	// Randomize the struct values again, except for the primary key values, so we can call update.
	err = boil.RandomizeStruct(&o, chunkDBTypes, true, chunkPrimaryKeyColumns...)
	if err != nil {
		t.Errorf("Unable to randomize Chunk struct members excluding primary keys: %s", err)
	}

	colsWithoutPrimKeys := boil.SetComplement(chunkColumns, chunkPrimaryKeyColumns)

	if err = o.UpdateG(colsWithoutPrimKeys...); err != nil {
		t.Errorf("Unable to update the Chunk row: %s", err)
	}

	if err = o1.ReloadG(); err != nil {
		t.Errorf("Unable to reload Chunk object: %s", err)
	}
	err = chunkCompareVals(&o, o1, true)
	if err != nil {
		t.Error(err)
	}

	chunksDeleteAllRows(t)
}

func TestChunksReloadAll(t *testing.T) {
	var err error

	o1 := make(ChunkSlice, 3)
	o2 := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o1, chunkDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o1); i++ {
		if err = o1[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o1[i], err)
		}
	}

	for i := 0; i < len(o1); i++ {
		o2[i], err = ChunkFindG(o1[i].ID)
		if err != nil {
			t.Errorf("Unable to find Chunk row.")
		}
		err = chunkCompareVals(o1[i], o2[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	// Randomize the struct values again, except for the primary key values, so we can call update.
	err = boil.RandomizeSlice(&o1, chunkDBTypes, false, chunkPrimaryKeyColumns...)
	if err != nil {
		t.Errorf("Unable to randomize Chunk slice excluding primary keys: %s", err)
	}

	colsWithoutPrimKeys := boil.SetComplement(chunkColumns, chunkPrimaryKeyColumns)

	for i := 0; i < len(o1); i++ {
		if err = o1[i].UpdateG(colsWithoutPrimKeys...); err != nil {
			t.Errorf("Unable to update the Chunk row: %s", err)
		}
	}

	if err = o2.ReloadAllG(); err != nil {
		t.Errorf("Unable to reload Chunk object: %s", err)
	}

	for i := 0; i < len(o1); i++ {
		err = chunkCompareVals(o2[i], o1[i], true)
		if err != nil {
			t.Error(err)
		}
	}

	chunksDeleteAllRows(t)
}

func TestChunksSelect(t *testing.T) {
	// Only run this test if there are ample cols to test on
	if len(chunkAutoIncrementColumns) == 0 {
		return
	}

	var err error

	x := &struct {
	}{}

	item := Chunk{}

	blacklistCols := boil.SetMerge(chunkAutoIncrementColumns, chunkPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, chunkDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert item Chunk:\n%#v\nErr: %s", item, err)
	}

	err = ChunksG(qm.Select(chunkAutoIncrementColumns...), qm.Where(`"id"=$1`, item.ID)).Bind(x)
	if err != nil {
		t.Errorf("Unable to select insert results with bind: %s", err)
	}

	chunksDeleteAllRows(t)
}

func TestChunksUpdate(t *testing.T) {
	var err error

	item := Chunk{}
	boil.RandomizeValidatedStruct(&item, chunkValidatedColumns, chunkDBTypes)
	if err = item.InsertG(); err != nil {
		t.Errorf("Unable to insert zero-value item Chunk:\n%#v\nErr: %s", item, err)
	}

	blacklistCols := boil.SetMerge(chunkAutoIncrementColumns, chunkPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, chunkDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	whitelist := boil.SetComplement(chunkColumns, chunkAutoIncrementColumns)
	if err = item.UpdateG(whitelist...); err != nil {
		t.Errorf("Unable to update Chunk: %s", err)
	}

	var j *Chunk
	j, err = ChunkFindG(item.ID)
	if err != nil {
		t.Errorf("Unable to find Chunk row: %s", err)
	}

	err = chunkCompareVals(&item, j, true)
	if err != nil {
		t.Error(err)
	}

	wl := item.generateUpdateColumns("test")
	if len(wl) != 1 && wl[0] != "test" {
		t.Errorf("Expected generateUpdateColumns whitelist to match expected whitelist")
	}

	wl = item.generateUpdateColumns()
	if len(wl) == 0 && len(chunkColumnsWithoutDefault) > 0 {
		t.Errorf("Expected generateUpdateColumns to build a whitelist for Chunk, but got 0 results")
	}

	chunksDeleteAllRows(t)
}

func TestChunksSliceUpdateAll(t *testing.T) {
	var err error

	// insert random columns to test UpdateAll
	o := make(ChunkSlice, 3)
	j := make(ChunkSlice, 3)

	if err = boil.RandomizeSlice(&o, chunkDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].InsertG(); err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	vals := M{}

	tmp := Chunk{}
	blacklist := boil.SetMerge(chunkPrimaryKeyColumns, chunkUniqueColumns)
	if err = boil.RandomizeStruct(&tmp, chunkDBTypes, false, blacklist...); err != nil {
		t.Errorf("Unable to randomize struct Chunk: %s", err)
	}

	// Build the columns and column values from the randomized struct
	tmpVal := reflect.Indirect(reflect.ValueOf(tmp))
	nonBlacklist := boil.SetComplement(chunkColumns, blacklist)
	for _, col := range nonBlacklist {
		vals[col] = tmpVal.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	err = o.UpdateAllG(vals)
	if err != nil {
		t.Errorf("Failed to update all for Chunk: %s", err)
	}

	for i := 0; i < len(o); i++ {
		j[i], err = ChunkFindG(o[i].ID)
		if err != nil {
			t.Errorf("Unable to find Chunk row: %s", err)
		}

		err = chunkCompareVals(j[i], &tmp, true, blacklist...)
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < len(o); i++ {
		// Ensure Find found the correct primary key ID's
		orig := boil.GetStructValues(o[i], chunkPrimaryKeyColumns...)
		new := boil.GetStructValues(j[i], chunkPrimaryKeyColumns...)

		if !reflect.DeepEqual(orig, new) {
			t.Errorf("object %d): primary keys do not match:\n\n%#v\n%#v", i, orig, new)
		}
	}

	chunksDeleteAllRows(t)
}

func TestChunksUpsert(t *testing.T) {
	var err error

	o := Chunk{}

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
	expectedQuery := `INSERT INTO chunks ("thing", "stuff") VALUES ($1, $2) ON CONFLICT DO NOTHING`

	if query != expectedQuery {
		t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
	}

	query = o.generateUpsertQuery(true, upsertCols)
	expectedQuery = `INSERT INTO chunks ("thing", "stuff") VALUES ($1, $2) ON CONFLICT ("key1", "key2") DO UPDATE SET "aaa" = EXCLUDED."aaa", "bbb" = EXCLUDED."bbb"`

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
	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.UpsertG(false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Chunk: %s", err)
	}

	compare, err := ChunkFindG(o.ID)
	if err != nil {
		t.Errorf("Unable to find Chunk: %s", err)
	}
	err = chunkCompareVals(&o, compare, true)
	if err != nil {
		t.Error(err)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = boil.RandomizeStruct(&o, chunkDBTypes, false, chunkPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.UpsertG(true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Chunk: %s", err)
	}

	compare, err = ChunkFindG(o.ID)
	if err != nil {
		t.Errorf("Unable to find Chunk: %s", err)
	}
	err = chunkCompareVals(&o, compare, true)
	if err != nil {
		t.Error(err)
	}

	chunksDeleteAllRows(t)
}

