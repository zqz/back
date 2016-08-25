package models

import (
	"testing"
	"reflect"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/boil/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testFiles(t *testing.T) {
	t.Parallel()

	query := Files(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testFilesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = file.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFilesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Files(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFilesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileSlice{file}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFilesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := FileExists(tx, file.ID)
	if err != nil {
		t.Errorf("Unable to check if File exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected FileExistsG to return true, but got false.")
	}
}

func testFilesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	fileFound, err := FileFind(tx, file.ID)
	if err != nil {
		t.Error(err)
	}

	if fileFound == nil {
		t.Error("want a record, got nil")
	}
}

func testFilesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Files(tx).Bind(file); err != nil {
		t.Error(err)
	}
}

func testFilesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Files(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testFilesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileOne := &File{}
	fileTwo := &File{}
	if err = randomize.Struct(seed, fileOne, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}
	if err = randomize.Struct(seed, fileTwo, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Files(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testFilesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	fileOne := &File{}
	fileTwo := &File{}
	if err = randomize.Struct(seed, fileOne, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}
	if err = randomize.Struct(seed, fileTwo, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

var fileDBTypes = map[string]string{"NumChunks": "integer", "State": "integer", "Name": "text", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone", "ID": "uuid", "Size": "integer", "Slug": "text", "Hash": "text", "Type": "text"}

func testFilesInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	var o File
	o = File{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &o, fileDBTypes, true); err != nil {
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

func testFilesSliceInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	o := make(FileSlice, 3)

	seed := randomize.NewSeed()
	for i := range o {
		o[i] = &File{}
		if err = randomize.Struct(seed, o[i], fileDBTypes, true); err != nil {
			t.Errorf("Could not randomize struct: %s", err)
		}
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(filePrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(filePrimaryKeyColumns)*3, len(args))
	}

	argC := 0
	for i := 0; i < 3; i++ {

		if o[i].ID != args[argC] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
		argC++
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

func testFilesHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &File{}
	o := &File{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, fileDBTypes, false); err != nil {
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
}

func testFilesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFilesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx, fileColumns...); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileToManyChunks(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Chunk

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	seed := randomize.NewSeed()
	randomize.Struct(seed, &b, chunkDBTypes, false, "file_id")
	randomize.Struct(seed, &c, chunkDBTypes, false, "file_id")

	b.FileID = a.ID
	c.FileID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	chunk, err := a.Chunks(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range chunk {
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

	slice := FileSlice{&a}
	if err = a.Loaded.LoadChunks(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.Loaded.Chunks); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.Loaded.Chunks = nil
	if err = a.Loaded.LoadChunks(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.Loaded.Chunks); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", chunk)
	}
}

func testFileToManyThumbnails(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Thumbnail

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	seed := randomize.NewSeed()
	randomize.Struct(seed, &b, thumbnailDBTypes, false, "file_id")
	randomize.Struct(seed, &c, thumbnailDBTypes, false, "file_id")

	b.FileID = a.ID
	c.FileID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	thumbnail, err := a.Thumbnails(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range thumbnail {
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

	slice := FileSlice{&a}
	if err = a.Loaded.LoadThumbnails(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.Loaded.Thumbnails); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.Loaded.Thumbnails = nil
	if err = a.Loaded.LoadThumbnails(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.Loaded.Thumbnails); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", thumbnail)
	}
}



func testFilesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = file.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testFilesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileSlice{file}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testFilesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Files(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

func testFilesUpdate(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, file, fileDBTypes, true, filePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	// If table only contains primary key columns, we need to pass
	// them into a whitelist to get a valid test result,
	// otherwise the Update method will error because it will not be able to
	// generate a whitelist (due to it excluding primary key columns).
	if strmangle.StringSliceMatch(fileColumns, filePrimaryKeyColumns) {
		if err = file.Update(tx, filePrimaryKeyColumns...); err != nil {
			t.Error(err)
		}
	} else {
		if err = file.Update(tx); err != nil {
			t.Error(err)
		}
	}
}

func testFilesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, file, fileDBTypes, true, filePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(fileColumns, filePrimaryKeyColumns) {
		fields = fileColumns
	} else {
		fields = strmangle.SetComplement(
			fileColumns,
			filePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(file))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := FileSlice{file}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}

func testFilesUpsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	file := File{}
	if err = randomize.Struct(seed, &file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert File: %s", err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &file, fileDBTypes, false, filePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = file.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert File: %s", err)
	}

	count, err = Files(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

