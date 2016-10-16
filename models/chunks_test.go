package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testChunks(t *testing.T) {
	t.Parallel()

	query := Chunks(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testChunksDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = chunk.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChunksQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Chunks(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChunksSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ChunkSlice{chunk}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testChunksExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ChunkExists(tx, chunk.ID)
	if err != nil {
		t.Errorf("Unable to check if Chunk exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ChunkExistsG to return true, but got false.")
	}
}
func testChunksFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	chunkFound, err := FindChunk(tx, chunk.ID)
	if err != nil {
		t.Error(err)
	}

	if chunkFound == nil {
		t.Error("want a record, got nil")
	}
}
func testChunksBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Chunks(tx).Bind(chunk); err != nil {
		t.Error(err)
	}
}

func testChunksOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Chunks(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testChunksAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunkOne := &Chunk{}
	chunkTwo := &Chunk{}
	if err = randomize.Struct(seed, chunkOne, chunkDBTypes, false, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}
	if err = randomize.Struct(seed, chunkTwo, chunkDBTypes, false, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunkOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = chunkTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Chunks(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testChunksCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	chunkOne := &Chunk{}
	chunkTwo := &Chunk{}
	if err = randomize.Struct(seed, chunkOne, chunkDBTypes, false, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}
	if err = randomize.Struct(seed, chunkTwo, chunkDBTypes, false, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunkOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = chunkTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func chunkBeforeInsertHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterInsertHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterSelectHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkBeforeUpdateHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterUpdateHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkBeforeDeleteHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterDeleteHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkBeforeUpsertHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func chunkAfterUpsertHook(e boil.Executor, o *Chunk) error {
	*o = Chunk{}
	return nil
}

func testChunksHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Chunk{}
	o := &Chunk{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, chunkDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Chunk object: %s", err)
	}

	AddChunkHook(boil.BeforeInsertHook, chunkBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	chunkBeforeInsertHooks = []ChunkHook{}

	AddChunkHook(boil.AfterInsertHook, chunkAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	chunkAfterInsertHooks = []ChunkHook{}

	AddChunkHook(boil.AfterSelectHook, chunkAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	chunkAfterSelectHooks = []ChunkHook{}

	AddChunkHook(boil.BeforeUpdateHook, chunkBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	chunkBeforeUpdateHooks = []ChunkHook{}

	AddChunkHook(boil.AfterUpdateHook, chunkAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	chunkAfterUpdateHooks = []ChunkHook{}

	AddChunkHook(boil.BeforeDeleteHook, chunkBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	chunkBeforeDeleteHooks = []ChunkHook{}

	AddChunkHook(boil.AfterDeleteHook, chunkAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	chunkAfterDeleteHooks = []ChunkHook{}

	AddChunkHook(boil.BeforeUpsertHook, chunkBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	chunkBeforeUpsertHooks = []ChunkHook{}

	AddChunkHook(boil.AfterUpsertHook, chunkAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	chunkAfterUpsertHooks = []ChunkHook{}
}
func testChunksInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testChunksInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx, chunkColumns...); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testChunkToOneFileUsingFile(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Chunk
	var foreign File

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.FileID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.File(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ChunkSlice{&local}
	if err = local.L.LoadFile(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.File == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.File = nil
	if err = local.L.LoadFile(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.File == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testChunkToOneSetOpFileUsingFile(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Chunk
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, chunkDBTypes, false, strmangle.SetComplement(chunkPrimaryKeyColumns, chunkColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*File{&b, &c} {
		err = a.SetFile(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.File != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Chunks[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.FileID != x.ID {
			t.Error("foreign key was wrong value", a.FileID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.FileID))
		reflect.Indirect(reflect.ValueOf(&a.FileID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.FileID != x.ID {
			t.Error("foreign key was wrong value", a.FileID, x.ID)
		}
	}
}
func testChunksReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = chunk.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testChunksReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ChunkSlice{chunk}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testChunksSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Chunks(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	chunkDBTypes = map[string]string{"CreatedAt": "timestamp without time zone", "FileID": "uuid", "Hash": "text", "ID": "uuid", "Position": "integer", "Size": "integer", "UpdatedAt": "timestamp without time zone"}
	_            = bytes.MinRead
)

func testChunksUpdate(t *testing.T) {
	t.Parallel()

	if len(chunkColumns) == len(chunkPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = chunk.Update(tx); err != nil {
		t.Error(err)
	}
}

func testChunksSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(chunkColumns) == len(chunkPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	chunk := &Chunk{}
	if err = randomize.Struct(seed, chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, chunk, chunkDBTypes, true, chunkPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(chunkColumns, chunkPrimaryKeyColumns) {
		fields = chunkColumns
	} else {
		fields = strmangle.SetComplement(
			chunkColumns,
			chunkPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(chunk))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ChunkSlice{chunk}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testChunksUpsert(t *testing.T) {
	t.Parallel()

	if len(chunkColumns) == len(chunkPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	chunk := Chunk{}
	if err = randomize.Struct(seed, &chunk, chunkDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = chunk.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Chunk: %s", err)
	}

	count, err := Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &chunk, chunkDBTypes, false, chunkPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Chunk struct: %s", err)
	}

	if err = chunk.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Chunk: %s", err)
	}

	count, err = Chunks(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
