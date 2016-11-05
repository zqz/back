package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testDownloads(t *testing.T) {
	t.Parallel()

	query := Downloads(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testDownloadsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = download.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDownloadsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Downloads(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDownloadsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := DownloadSlice{download}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testDownloadsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := DownloadExists(tx, download.ID)
	if err != nil {
		t.Errorf("Unable to check if Download exists: %s", err)
	}
	if !e {
		t.Errorf("Expected DownloadExistsG to return true, but got false.")
	}
}
func testDownloadsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	downloadFound, err := FindDownload(tx, download.ID)
	if err != nil {
		t.Error(err)
	}

	if downloadFound == nil {
		t.Error("want a record, got nil")
	}
}
func testDownloadsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Downloads(tx).Bind(download); err != nil {
		t.Error(err)
	}
}

func testDownloadsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Downloads(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testDownloadsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	downloadOne := &Download{}
	downloadTwo := &Download{}
	if err = randomize.Struct(seed, downloadOne, downloadDBTypes, false, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}
	if err = randomize.Struct(seed, downloadTwo, downloadDBTypes, false, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = downloadOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = downloadTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Downloads(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testDownloadsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	downloadOne := &Download{}
	downloadTwo := &Download{}
	if err = randomize.Struct(seed, downloadOne, downloadDBTypes, false, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}
	if err = randomize.Struct(seed, downloadTwo, downloadDBTypes, false, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = downloadOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = downloadTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func downloadBeforeInsertHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadAfterInsertHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadAfterSelectHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadBeforeUpdateHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadAfterUpdateHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadBeforeDeleteHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadAfterDeleteHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadBeforeUpsertHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func downloadAfterUpsertHook(e boil.Executor, o *Download) error {
	*o = Download{}
	return nil
}

func testDownloadsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Download{}
	o := &Download{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, downloadDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Download object: %s", err)
	}

	AddDownloadHook(boil.BeforeInsertHook, downloadBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	downloadBeforeInsertHooks = []DownloadHook{}

	AddDownloadHook(boil.AfterInsertHook, downloadAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	downloadAfterInsertHooks = []DownloadHook{}

	AddDownloadHook(boil.AfterSelectHook, downloadAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	downloadAfterSelectHooks = []DownloadHook{}

	AddDownloadHook(boil.BeforeUpdateHook, downloadBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	downloadBeforeUpdateHooks = []DownloadHook{}

	AddDownloadHook(boil.AfterUpdateHook, downloadAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	downloadAfterUpdateHooks = []DownloadHook{}

	AddDownloadHook(boil.BeforeDeleteHook, downloadBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	downloadBeforeDeleteHooks = []DownloadHook{}

	AddDownloadHook(boil.AfterDeleteHook, downloadAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	downloadAfterDeleteHooks = []DownloadHook{}

	AddDownloadHook(boil.BeforeUpsertHook, downloadBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	downloadBeforeUpsertHooks = []DownloadHook{}

	AddDownloadHook(boil.AfterUpsertHook, downloadAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	downloadAfterUpsertHooks = []DownloadHook{}
}
func testDownloadsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDownloadsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx, downloadColumns...); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDownloadToOneFileUsingFile(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Download
	var foreign File

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	local.FileID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.FileID.String = foreign.ID
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

	slice := DownloadSlice{&local}
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

func testDownloadToOneSetOpFileUsingFile(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Download
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, downloadDBTypes, false, strmangle.SetComplement(downloadPrimaryKeyColumns, downloadColumnsWithoutDefault)...); err != nil {
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

		if x.R.Downloads[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.FileID.String != x.ID {
			t.Error("foreign key was wrong value", a.FileID.String)
		}

		zero := reflect.Zero(reflect.TypeOf(a.FileID.String))
		reflect.Indirect(reflect.ValueOf(&a.FileID.String)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.FileID.String != x.ID {
			t.Error("foreign key was wrong value", a.FileID.String, x.ID)
		}
	}
}

func testDownloadToOneRemoveOpFileUsingFile(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Download
	var b File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, downloadDBTypes, false, strmangle.SetComplement(downloadPrimaryKeyColumns, downloadColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetFile(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveFile(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.File(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.File != nil {
		t.Error("R struct entry should be nil")
	}

	if a.FileID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Downloads) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testDownloadsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = download.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testDownloadsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := DownloadSlice{download}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testDownloadsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Downloads(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	downloadDBTypes = map[string]string{"CacheHit": "boolean", "CreatedAt": "timestamp without time zone", "FileID": "uuid", "ID": "integer", "Ip": "inet"}
	_               = bytes.MinRead
)

func testDownloadsUpdate(t *testing.T) {
	t.Parallel()

	if len(downloadColumns) == len(downloadPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	if err = download.Update(tx); err != nil {
		t.Error(err)
	}
}

func testDownloadsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(downloadColumns) == len(downloadPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	download := &Download{}
	if err = randomize.Struct(seed, download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, download, downloadDBTypes, true, downloadPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(downloadColumns, downloadPrimaryKeyColumns) {
		fields = downloadColumns
	} else {
		fields = strmangle.SetComplement(
			downloadColumns,
			downloadPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(download))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := DownloadSlice{download}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testDownloadsUpsert(t *testing.T) {
	t.Parallel()

	if len(downloadColumns) == len(downloadPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	download := Download{}
	if err = randomize.Struct(seed, &download, downloadDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = download.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Download: %s", err)
	}

	count, err := Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &download, downloadDBTypes, false, downloadPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Download struct: %s", err)
	}

	if err = download.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Download: %s", err)
	}

	count, err = Downloads(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
