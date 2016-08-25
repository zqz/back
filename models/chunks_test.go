package models

import (
	"testing"
	"reflect"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/boil/randomize"
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
	if e != true {
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

	chunkFound, err := ChunkFind(tx, chunk.ID)
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

var chunkDBTypes = map[string]string{"ID": "uuid", "FileID": "uuid", "Size": "integer", "Hash": "text", "Position": "integer", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone"}

func testChunksInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	var o Chunk
	o = Chunk{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &o, chunkDBTypes, true); err != nil {
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

func testChunksSliceInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	o := make(ChunkSlice, 3)

	seed := randomize.NewSeed()
	for i := range o {
		o[i] = &Chunk{}
		if err = randomize.Struct(seed, o[i], chunkDBTypes, true); err != nil {
			t.Errorf("Could not randomize struct: %s", err)
		}
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(chunkPrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(chunkPrimaryKeyColumns)*3, len(args))
	}

	argC := 0
	for i := 0; i < 3; i++ {

		if o[i].ID != args[argC] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
		argC++
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

func testChunksHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Chunk{}
	o := &Chunk{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, chunkDBTypes, false); err != nil {
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



func testChunkToOneFile_File(t *testing.T) {
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
	check, err := local.File(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ChunkSlice{&local}
	if err = local.Loaded.LoadFile(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.Loaded.File == nil {
		t.Error("struct should have been eager loaded")
	}

	local.Loaded.File = nil
	if err = local.Loaded.LoadFile(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.Loaded.File == nil {
		t.Error("struct should have been eager loaded")
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

func testChunksUpdate(t *testing.T) {
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

	// If table only contains primary key columns, we need to pass
	// them into a whitelist to get a valid test result,
	// otherwise the Update method will error because it will not be able to
	// generate a whitelist (due to it excluding primary key columns).
	if strmangle.StringSliceMatch(chunkColumns, chunkPrimaryKeyColumns) {
		if err = chunk.Update(tx, chunkPrimaryKeyColumns...); err != nil {
			t.Error(err)
		}
	} else {
		if err = chunk.Update(tx); err != nil {
			t.Error(err)
		}
	}
}

func testChunksSliceUpdateAll(t *testing.T) {
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

