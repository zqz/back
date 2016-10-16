package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
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
	if !e {
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

	fileFound, err := FindFile(tx, file.ID)
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
func fileBeforeInsertHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileAfterInsertHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileAfterSelectHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileBeforeUpdateHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileAfterUpdateHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileBeforeDeleteHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileAfterDeleteHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileBeforeUpsertHook(e boil.Executor, o *File) error {
	*o = File{}
	return nil
}

func fileAfterUpsertHook(e boil.Executor, o *File) error {
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

	AddFileHook(boil.BeforeInsertHook, fileBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	fileBeforeInsertHooks = []FileHook{}

	AddFileHook(boil.AfterInsertHook, fileAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	fileAfterInsertHooks = []FileHook{}

	AddFileHook(boil.AfterSelectHook, fileAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	fileAfterSelectHooks = []FileHook{}

	AddFileHook(boil.BeforeUpdateHook, fileBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	fileBeforeUpdateHooks = []FileHook{}

	AddFileHook(boil.AfterUpdateHook, fileAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	fileAfterUpdateHooks = []FileHook{}

	AddFileHook(boil.BeforeDeleteHook, fileBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	fileBeforeDeleteHooks = []FileHook{}

	AddFileHook(boil.AfterDeleteHook, fileAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	fileAfterDeleteHooks = []FileHook{}

	AddFileHook(boil.BeforeUpsertHook, fileBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	fileBeforeUpsertHooks = []FileHook{}

	AddFileHook(boil.AfterUpsertHook, fileAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	fileAfterUpsertHooks = []FileHook{}
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

func testFileToManyDownloads(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Download

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, downloadDBTypes, false, downloadColumnsWithDefault...)
	randomize.Struct(seed, &c, downloadDBTypes, false, downloadColumnsWithDefault...)
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

	download, err := a.Downloads(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range download {
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

	slice := FileSlice{&a}
	if err = a.L.LoadDownloads(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Downloads); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Downloads = nil
	if err = a.L.LoadDownloads(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Downloads); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", download)
	}
}

func testFileToManyChunks(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Chunk

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, chunkDBTypes, false, chunkColumnsWithDefault...)
	randomize.Struct(seed, &c, chunkDBTypes, false, chunkColumnsWithDefault...)

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
	if err = a.L.LoadChunks(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Chunks); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Chunks = nil
	if err = a.L.LoadChunks(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Chunks); got != 2 {
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

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, thumbnailDBTypes, false, thumbnailColumnsWithDefault...)
	randomize.Struct(seed, &c, thumbnailDBTypes, false, thumbnailColumnsWithDefault...)

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
	if err = a.L.LoadThumbnails(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Thumbnails); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Thumbnails = nil
	if err = a.L.LoadThumbnails(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Thumbnails); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", thumbnail)
	}
}

func testFileToManyAddOpDownloads(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Download

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Download{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, downloadDBTypes, false, strmangle.SetComplement(downloadPrimaryKeyColumns, downloadColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Download{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDownloads(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.FileID.String {
			t.Error("foreign key was wrong value", a.ID, first.FileID.String)
		}
		if a.ID != second.FileID.String {
			t.Error("foreign key was wrong value", a.ID, second.FileID.String)
		}

		if first.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Downloads[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Downloads[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Downloads(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testFileToManySetOpDownloads(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Download

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Download{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, downloadDBTypes, false, strmangle.SetComplement(downloadPrimaryKeyColumns, downloadColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.SetDownloads(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Downloads(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDownloads(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Downloads(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.FileID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.FileID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.FileID.String {
		t.Error("foreign key was wrong value", a.ID, d.FileID.String)
	}
	if a.ID != e.FileID.String {
		t.Error("foreign key was wrong value", a.ID, e.FileID.String)
	}

	if b.R.File != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.File != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.File != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.File != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Downloads[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Downloads[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testFileToManyRemoveOpDownloads(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Download

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Download{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, downloadDBTypes, false, strmangle.SetComplement(downloadPrimaryKeyColumns, downloadColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddDownloads(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Downloads(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDownloads(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Downloads(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.FileID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.FileID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.File != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.File != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.File != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.File != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.Downloads) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Downloads[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Downloads[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testFileToManyAddOpChunks(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Chunk

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Chunk{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, chunkDBTypes, false, strmangle.SetComplement(chunkPrimaryKeyColumns, chunkColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Chunk{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddChunks(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.FileID {
			t.Error("foreign key was wrong value", a.ID, first.FileID)
		}
		if a.ID != second.FileID {
			t.Error("foreign key was wrong value", a.ID, second.FileID)
		}

		if first.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Chunks[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Chunks[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Chunks(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testFileToManyAddOpThumbnails(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Thumbnail

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Thumbnail{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, thumbnailDBTypes, false, strmangle.SetComplement(thumbnailPrimaryKeyColumns, thumbnailColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Thumbnail{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddThumbnails(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.FileID {
			t.Error("foreign key was wrong value", a.ID, first.FileID)
		}
		if a.ID != second.FileID {
			t.Error("foreign key was wrong value", a.ID, second.FileID)
		}

		if first.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Thumbnails[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Thumbnails[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Thumbnails(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
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

var (
	fileDBTypes = map[string]string{"CreatedAt": "timestamp without time zone", "Hash": "text", "ID": "uuid", "Name": "text", "NumChunks": "integer", "Size": "integer", "Slug": "text", "State": "integer", "Type": "text", "UpdatedAt": "timestamp without time zone"}
	_           = bytes.MinRead
)

func testFilesUpdate(t *testing.T) {
	t.Parallel()

	if len(fileColumns) == len(filePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

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

	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = file.Update(tx); err != nil {
		t.Error(err)
	}
}

func testFilesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(fileColumns) == len(filePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

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

	if len(fileColumns) == len(filePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

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
