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

func testChannelMembers(t *testing.T) {
	t.Parallel()

	query := ChannelMembers()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testChannelMembersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
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

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChannelMembersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ChannelMembers().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChannelMembersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ChannelMemberSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testChannelMembersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ChannelMemberExists(ctx, tx, o.ChannelID, o.UserID)
	if err != nil {
		t.Errorf("Unable to check if ChannelMember exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ChannelMemberExists to return true, but got false.")
	}
}

func testChannelMembersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	channelMemberFound, err := FindChannelMember(ctx, tx, o.ChannelID, o.UserID)
	if err != nil {
		t.Error(err)
	}

	if channelMemberFound == nil {
		t.Error("want a record, got nil")
	}
}

func testChannelMembersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ChannelMembers().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testChannelMembersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ChannelMembers().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testChannelMembersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	channelMemberOne := &ChannelMember{}
	channelMemberTwo := &ChannelMember{}
	if err = randomize.Struct(seed, channelMemberOne, channelMemberDBTypes, false, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}
	if err = randomize.Struct(seed, channelMemberTwo, channelMemberDBTypes, false, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = channelMemberOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = channelMemberTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ChannelMembers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testChannelMembersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	channelMemberOne := &ChannelMember{}
	channelMemberTwo := &ChannelMember{}
	if err = randomize.Struct(seed, channelMemberOne, channelMemberDBTypes, false, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}
	if err = randomize.Struct(seed, channelMemberTwo, channelMemberDBTypes, false, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = channelMemberOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = channelMemberTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func channelMemberBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func channelMemberAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ChannelMember) error {
	*o = ChannelMember{}
	return nil
}

func testChannelMembersHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ChannelMember{}
	o := &ChannelMember{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, channelMemberDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ChannelMember object: %s", err)
	}

	AddChannelMemberHook(boil.BeforeInsertHook, channelMemberBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	channelMemberBeforeInsertHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.AfterInsertHook, channelMemberAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	channelMemberAfterInsertHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.AfterSelectHook, channelMemberAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	channelMemberAfterSelectHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.BeforeUpdateHook, channelMemberBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	channelMemberBeforeUpdateHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.AfterUpdateHook, channelMemberAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	channelMemberAfterUpdateHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.BeforeDeleteHook, channelMemberBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	channelMemberBeforeDeleteHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.AfterDeleteHook, channelMemberAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	channelMemberAfterDeleteHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.BeforeUpsertHook, channelMemberBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	channelMemberBeforeUpsertHooks = []ChannelMemberHook{}

	AddChannelMemberHook(boil.AfterUpsertHook, channelMemberAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	channelMemberAfterUpsertHooks = []ChannelMemberHook{}
}

func testChannelMembersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testChannelMembersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(channelMemberColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testChannelMemberToOneChannelUsingChannel(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local ChannelMember
	var foreign Channel

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, channelMemberDBTypes, false, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, channelDBTypes, false, channelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Channel struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.ChannelID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Channel().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ChannelMemberSlice{&local}
	if err = local.L.LoadChannel(ctx, tx, false, (*[]*ChannelMember)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Channel == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Channel = nil
	if err = local.L.LoadChannel(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Channel == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testChannelMemberToOneUserUsingUser(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local ChannelMember
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, channelMemberDBTypes, false, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.UserID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.User().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ChannelMemberSlice{&local}
	if err = local.L.LoadUser(ctx, tx, false, (*[]*ChannelMember)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.User = nil
	if err = local.L.LoadUser(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testChannelMemberToOneSetOpChannelUsingChannel(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ChannelMember
	var b, c Channel

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, channelMemberDBTypes, false, strmangle.SetComplement(channelMemberPrimaryKeyColumns, channelMemberColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, channelDBTypes, false, strmangle.SetComplement(channelPrimaryKeyColumns, channelColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, channelDBTypes, false, strmangle.SetComplement(channelPrimaryKeyColumns, channelColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Channel{&b, &c} {
		err = a.SetChannel(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Channel != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ChannelMembers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ChannelID != x.ID {
			t.Error("foreign key was wrong value", a.ChannelID)
		}

		if exists, err := ChannelMemberExists(ctx, tx, a.ChannelID, a.UserID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testChannelMemberToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ChannelMember
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, channelMemberDBTypes, false, strmangle.SetComplement(channelMemberPrimaryKeyColumns, channelMemberColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetUser(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.User != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ChannelMembers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID != x.ID {
			t.Error("foreign key was wrong value", a.UserID)
		}

		if exists, err := ChannelMemberExists(ctx, tx, a.ChannelID, a.UserID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}

func testChannelMembersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
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

func testChannelMembersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ChannelMemberSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testChannelMembersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ChannelMembers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	channelMemberDBTypes = map[string]string{`ChannelID`: `integer`, `UserID`: `integer`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`, `DeletedAt`: `timestamp without time zone`}
	_                    = bytes.MinRead
)

func testChannelMembersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(channelMemberPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(channelMemberAllColumns) == len(channelMemberPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testChannelMembersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(channelMemberAllColumns) == len(channelMemberPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ChannelMember{}
	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, channelMemberDBTypes, true, channelMemberPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(channelMemberAllColumns, channelMemberPrimaryKeyColumns) {
		fields = channelMemberAllColumns
	} else {
		fields = strmangle.SetComplement(
			channelMemberAllColumns,
			channelMemberPrimaryKeyColumns,
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

	slice := ChannelMemberSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testChannelMembersUpsert(t *testing.T) {
	t.Parallel()

	if len(channelMemberAllColumns) == len(channelMemberPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ChannelMember{}
	if err = randomize.Struct(seed, &o, channelMemberDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ChannelMember: %s", err)
	}

	count, err := ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, channelMemberDBTypes, false, channelMemberPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ChannelMember struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ChannelMember: %s", err)
	}

	count, err = ChannelMembers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
