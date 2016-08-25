package models

import (
	"testing"
	"reflect"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/boil/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testThumbnails(t *testing.T) {
	t.Parallel()

	query := Thumbnails(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testThumbnailsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = thumbnail.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testThumbnailsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Thumbnails(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testThumbnailsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ThumbnailSlice{thumbnail}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testThumbnailsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ThumbnailExists(tx, thumbnail.ID)
	if err != nil {
		t.Errorf("Unable to check if Thumbnail exists: %s", err)
	}
	if e != true {
		t.Errorf("Expected ThumbnailExistsG to return true, but got false.")
	}
}

func testThumbnailsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	thumbnailFound, err := ThumbnailFind(tx, thumbnail.ID)
	if err != nil {
		t.Error(err)
	}

	if thumbnailFound == nil {
		t.Error("want a record, got nil")
	}
}

func testThumbnailsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Thumbnails(tx).Bind(thumbnail); err != nil {
		t.Error(err)
	}
}

func testThumbnailsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Thumbnails(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testThumbnailsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnailOne := &Thumbnail{}
	thumbnailTwo := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnailOne, thumbnailDBTypes, false, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}
	if err = randomize.Struct(seed, thumbnailTwo, thumbnailDBTypes, false, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnailOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = thumbnailTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Thumbnails(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testThumbnailsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	thumbnailOne := &Thumbnail{}
	thumbnailTwo := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnailOne, thumbnailDBTypes, false, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}
	if err = randomize.Struct(seed, thumbnailTwo, thumbnailDBTypes, false, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnailOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = thumbnailTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

var thumbnailDBTypes = map[string]string{"ID": "uuid", "FileID": "uuid", "Size": "integer", "Hash": "text", "CreatedAt": "timestamp without time zone", "UpdatedAt": "timestamp without time zone"}

func testThumbnailsInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	var o Thumbnail
	o = Thumbnail{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &o, thumbnailDBTypes, true); err != nil {
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

func testThumbnailsSliceInPrimaryKeyArgs(t *testing.T) {
	t.Parallel()

	var err error
	o := make(ThumbnailSlice, 3)

	seed := randomize.NewSeed()
	for i := range o {
		o[i] = &Thumbnail{}
		if err = randomize.Struct(seed, o[i], thumbnailDBTypes, true); err != nil {
			t.Errorf("Could not randomize struct: %s", err)
		}
	}

	args := o.inPrimaryKeyArgs()

	if len(args) != len(thumbnailPrimaryKeyColumns)*3 {
		t.Errorf("Expected args to be len %d, but got %d", len(thumbnailPrimaryKeyColumns)*3, len(args))
	}

	argC := 0
	for i := 0; i < 3; i++ {

		if o[i].ID != args[argC] {
			t.Errorf("Expected args[%d] to be value of o.ID, but got %#v", i, args[i])
		}
		argC++
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

func testThumbnailsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Thumbnail{}
	o := &Thumbnail{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, thumbnailDBTypes, false); err != nil {
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
}

func testThumbnailsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testThumbnailsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx, thumbnailColumns...); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}



func testThumbnailToOneFile_File(t *testing.T) {
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
	check, err := local.File(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ThumbnailSlice{&local}
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


func testThumbnailsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = thumbnail.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testThumbnailsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ThumbnailSlice{thumbnail}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testThumbnailsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Thumbnails(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

func testThumbnailsUpdate(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	// If table only contains primary key columns, we need to pass
	// them into a whitelist to get a valid test result,
	// otherwise the Update method will error because it will not be able to
	// generate a whitelist (due to it excluding primary key columns).
	if strmangle.StringSliceMatch(thumbnailColumns, thumbnailPrimaryKeyColumns) {
		if err = thumbnail.Update(tx, thumbnailPrimaryKeyColumns...); err != nil {
			t.Error(err)
		}
	} else {
		if err = thumbnail.Update(tx); err != nil {
			t.Error(err)
		}
	}
}

func testThumbnailsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	thumbnail := &Thumbnail{}
	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, thumbnail, thumbnailDBTypes, true, thumbnailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(thumbnailColumns, thumbnailPrimaryKeyColumns) {
		fields = thumbnailColumns
	} else {
		fields = strmangle.SetComplement(
			thumbnailColumns,
			thumbnailPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(thumbnail))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ThumbnailSlice{thumbnail}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}

func testThumbnailsUpsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	thumbnail := Thumbnail{}
	if err = randomize.Struct(seed, &thumbnail, thumbnailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = thumbnail.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Thumbnail: %s", err)
	}

	count, err := Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &thumbnail, thumbnailDBTypes, false, thumbnailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Thumbnail struct: %s", err)
	}

	if err = thumbnail.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Thumbnail: %s", err)
	}

	count, err = Thumbnails(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

