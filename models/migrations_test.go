// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testMigrations(t *testing.T) {
	t.Parallel()

	query := Migrations()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testMigrationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMigrationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Migrations().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMigrationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MigrationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMigrationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := MigrationExists(ctx, tx, o.Version)
	if err != nil {
		t.Errorf("Unable to check if Migration exists: %s", err)
	}
	if !e {
		t.Errorf("Expected MigrationExists to return true, but got false.")
	}
}

func testMigrationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	migrationFound, err := FindMigration(ctx, tx, o.Version)
	if err != nil {
		t.Error(err)
	}

	if migrationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testMigrationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Migrations().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testMigrationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Migrations().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testMigrationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	migrationOne := &Migration{}
	migrationTwo := &Migration{}
	if err = randomize.Struct(seed, migrationOne, migrationDBTypes, false, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}
	if err = randomize.Struct(seed, migrationTwo, migrationDBTypes, false, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = migrationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = migrationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Migrations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testMigrationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	migrationOne := &Migration{}
	migrationTwo := &Migration{}
	if err = randomize.Struct(seed, migrationOne, migrationDBTypes, false, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}
	if err = randomize.Struct(seed, migrationTwo, migrationDBTypes, false, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = migrationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = migrationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func migrationBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func migrationAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Migration) error {
	*o = Migration{}
	return nil
}

func testMigrationsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Migration{}
	o := &Migration{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, migrationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Migration object: %s", err)
	}

	AddMigrationHook(boil.BeforeInsertHook, migrationBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	migrationBeforeInsertHooks = []MigrationHook{}

	AddMigrationHook(boil.AfterInsertHook, migrationAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	migrationAfterInsertHooks = []MigrationHook{}

	AddMigrationHook(boil.AfterSelectHook, migrationAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	migrationAfterSelectHooks = []MigrationHook{}

	AddMigrationHook(boil.BeforeUpdateHook, migrationBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	migrationBeforeUpdateHooks = []MigrationHook{}

	AddMigrationHook(boil.AfterUpdateHook, migrationAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	migrationAfterUpdateHooks = []MigrationHook{}

	AddMigrationHook(boil.BeforeDeleteHook, migrationBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	migrationBeforeDeleteHooks = []MigrationHook{}

	AddMigrationHook(boil.AfterDeleteHook, migrationAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	migrationAfterDeleteHooks = []MigrationHook{}

	AddMigrationHook(boil.BeforeUpsertHook, migrationBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	migrationBeforeUpsertHooks = []MigrationHook{}

	AddMigrationHook(boil.AfterUpsertHook, migrationAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	migrationAfterUpsertHooks = []MigrationHook{}
}

func testMigrationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMigrationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(migrationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMigrationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMigrationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MigrationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMigrationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Migrations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	migrationDBTypes = map[string]string{`Version`: `bigint`, `Dirty`: `boolean`}
	_                = bytes.MinRead
)

func testMigrationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(migrationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(migrationAllColumns) == len(migrationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testMigrationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(migrationAllColumns) == len(migrationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Migration{}
	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, migrationDBTypes, true, migrationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(migrationAllColumns, migrationPrimaryKeyColumns) {
		fields = migrationAllColumns
	} else {
		fields = strmangle.SetComplement(
			migrationAllColumns,
			migrationPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := MigrationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testMigrationsUpsert(t *testing.T) {
	t.Parallel()

	if len(migrationAllColumns) == len(migrationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Migration{}
	if err = randomize.Struct(seed, &o, migrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Migration: %s", err)
	}

	count, err := Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, migrationDBTypes, false, migrationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Migration struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Migration: %s", err)
	}

	count, err = Migrations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
