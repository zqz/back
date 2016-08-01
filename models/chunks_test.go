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

func TestChunks(t *testing.T) {
	var err error

	o := make(ChunkSlice, 2)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert two random objects to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Delete all rows to give a clean slate
	err = Chunks().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Chunks: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	var c int64
	c, err = Chunks().Count()

	if c != 0 {
		t.Errorf("Expected chunks table to be empty, but got %d rows", c)
	}

	o = make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Ensure Count is valid
	c, err = Chunks().Count()
	if c != 3 {
		t.Errorf("Expected chunks table to have 3 rows, but got %d", c)
	}

	// Attempt to retrieve all objects
	res, err := Chunks().All()
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
	err := Chunks().DeleteAll()
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
	c, err = Chunks().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after ObjDeleteAllRows() call, but got %d rows", c)
	}

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(o); i++ {
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// Test DeleteAll() query function
	err = Chunks().DeleteAll()
	if err != nil {
		t.Errorf("Unable to delete all from Chunks: %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = Chunks().Count()

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
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	// test DeleteAll slice function
	if err = o.DeleteAll(); err != nil {
		t.Errorf("Unable to objSlice.DeleteAll(): %s", err)
	}

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = Chunks().Count()

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
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	o[0].Delete()

	// Check number of rows in table to ensure DeleteAll was successful
	c, err = Chunks().Count()

	if c != 2 {
		t.Errorf("Expected 2 rows after obj.Delete() call, but got %d rows", c)
	}

	o[1].Delete()
	o[2].Delete()

	// Check number of rows in table to ensure Delete worked for all rows
	c, err = Chunks().Count()

	if c != 0 {
		t.Errorf("Expected 0 rows after all obj.Delete() calls, but got %d rows", c)
	}
}

func TestChunksFind(t *testing.T) {
	var err error

	o := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	for i := 0; i < len(o); i++ {
		if err = o[i].Insert(); err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(ChunkSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = ChunkFind(o[i].ID)
		chunkCompareVals(o[i], j[i], t)
	}

	f, err := ChunkFind(o[0].ID, chunkPrimaryKeyColumns...)

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

	if err = o.Insert(); err != nil {
		t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o, err)
	}

	j := Chunk{}

	err = Chunks(qm.Where(`"id"=$1`, o.ID)).Bind(&j)
	if err != nil {
		t.Errorf("Unable to call Bind on Chunk single object: %s", err)
	}

	chunkCompareVals(&o, &j, t)

	// insert 3 rows, attempt to bind into slice
	chunksDeleteAllRows(t)

	y := make(ChunkSlice, 3)
	if err = boil.RandomizeSlice(&y, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk slice: %s", err)
	}

	// insert random columns to test DeleteAll
	for i := 0; i < len(y); i++ {
		err = y[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", y[i], err)
		}
	}

	k := ChunkSlice{}
	err = Chunks().Bind(&k)
	if err != nil {
		t.Errorf("Unable to call Bind on Chunk slice of objects: %s", err)
	}

	if len(k) != 3 {
		t.Errorf("Expected 3 results, got %d", len(k))
	}

	for i := 0; i < len(y); i++ {
		chunkCompareVals(y[i], k[i], t)
	}

	chunksDeleteAllRows(t)
}

func TestChunksOne(t *testing.T) {
	var err error

	o := Chunk{}
	if err = boil.RandomizeStruct(&o, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = o.Insert(); err != nil {
		t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o, err)
	}

	j, err := Chunks().One()
	if err != nil {
		t.Errorf("Unable to fetch One Chunk result:\n#%v\nErr: %s", j, err)
	}

	chunkCompareVals(&o, j, t)

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
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	j, err := Chunks().All()
	if err != nil {
		t.Errorf("Unable to fetch All Chunk results: %s", err)
	}

	if len(j) != 3 {
		t.Errorf("Expected 3 results, got %d", len(j))
	}

	for i := 0; i < len(o); i++ {
		chunkCompareVals(o[i], j[i], t)
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
		err = o[i].Insert()
		if err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	c, err := Chunks().Count()
	if err != nil {
		t.Errorf("Unable to count query Chunk: %s", err)
	}

	if c != 3 {
		t.Errorf("Expected 3 results from count Chunk, got %d", c)
	}

	chunksDeleteAllRows(t)
}

var chunkDBTypes = map[string]string{"Position": "integer", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone", "ID": "uuid", "FileID": "uuid", "Size": "integer", "Hash": "text"}

func chunkCompareVals(o *Chunk, j *Chunk, t *testing.T) {
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

	if j.Position != o.Position {
		t.Errorf("Expected position columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.Position, j.Position)
	}

	if o.CreatedAt.Format("02/01/2006") != j.CreatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time created_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.CreatedAt.Format("02/01/2006"), j.CreatedAt.Format("02/01/2006"))
	}

	if o.UpdatedAt.Format("02/01/2006") != j.UpdatedAt.Format("02/01/2006") {
		t.Errorf("Expected Time updated_at column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.UpdatedAt.Format("02/01/2006"), j.UpdatedAt.Format("02/01/2006"))
	}
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
		if err = o[i].Insert(); err != nil {
			t.Errorf("Unable to insert Chunk:\n%#v\nErr: %s", o[i], err)
		}
	}

	j := make(ChunkSlice, 3)
	// Perform all Find queries and assign result objects to slice for comparison
	for i := 0; i < len(j); i++ {
		j[i], err = ChunkFind(o[i].ID)
		chunkCompareVals(o[i], j[i], t)
	}

	chunksDeleteAllRows(t)

	item := &Chunk{}
	if err = item.Insert(); err != nil {
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
	local.FileID.Valid = true

	if err := foreign.InsertX(tx); err != nil {
		t.Fatal(err)
	}

	local.FileID.String = foreign.ID
	if err := local.InsertX(tx); err != nil {
		t.Fatal(err)
	}
	check, err := local.FileX(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}
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

	if err = item.Insert(); err != nil {
		t.Errorf("Unable to insert item Chunk:\n%#v\nErr: %s", item, err)
	}

	err = Chunks(qm.Select(chunkAutoIncrementColumns...), qm.Where(`"id"=$1`, item.ID)).Bind(x)
	if err != nil {
		t.Errorf("Unable to select insert results with bind: %s", err)
	}

	chunksDeleteAllRows(t)
}

func TestChunksUpdate(t *testing.T) {
	var err error

	item := Chunk{}
	if err = item.Insert(); err != nil {
		t.Errorf("Unable to insert zero-value item Chunk:\n%#v\nErr: %s", item, err)
	}

	blacklistCols := boil.SetMerge(chunkAutoIncrementColumns, chunkPrimaryKeyColumns)
	if err = boil.RandomizeStruct(&item, chunkDBTypes, false, blacklistCols...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	whitelist := boil.SetComplement(chunkColumns, chunkAutoIncrementColumns)
	if err = item.Update(whitelist...); err != nil {
		t.Errorf("Unable to update Chunk: %s", err)
	}

	var j *Chunk
	j, err = ChunkFind(item.ID)
	if err != nil {
		t.Errorf("Unable to find Chunk row: %s", err)
	}

	chunkCompareVals(&item, j, t)

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

