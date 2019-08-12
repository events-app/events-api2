package menu_test

import (
	"context"
	"testing"
	"time"

	"github.com/events-app/events-api2/internal/menu"
	"github.com/events-app/events-api2/internal/platform/auth"
	"github.com/events-app/events-api2/internal/schema"
	"github.com/events-app/events-api2/internal/tests"
	"github.com/google/go-cmp/cmp"
)

func TestMenus(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	newC := menu.NewMenu{
		Name:    "New Menu",
		Content: "This is testing menu.",
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e", // This is just some random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

	c0, err := menu.Create(ctx, db, claims, newC, now)
	if err != nil {
		t.Fatalf("creating menu c0: %s", err)
	}

	c1, err := menu.Retrieve(ctx, db, c0.ID)
	if err != nil {
		t.Fatalf("getting menu c0: %s", err)
	}

	if diff := cmp.Diff(c1, c0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	update := menu.UpdateMenu{
		Name:    tests.StringPointer("New Menu"),
		Content: tests.StringPointer("Some new content of a menu"),
	}
	updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := menu.Update(ctx, db, claims, c0.ID, update, updatedTime); err != nil {
		t.Fatalf("creating menu c0: %s", err)
	}

	saved, err := menu.Retrieve(ctx, db, c0.ID)
	if err != nil {
		t.Fatalf("getting menu c0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original menu
	// and change just the fields we expect then diff it with what was saved.
	want := *c0
	want.Name = "New Menu"
	want.Cost = 25
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}

	if err := menu.Delete(ctx, db, c0.ID); err != nil {
		t.Fatalf("deleting menu: %v", err)
	}

	_, err = menu.Retrieve(ctx, db, c0.ID)
	if err == nil {
		t.Fatalf("should not be able to retrieve deleted menu")
	}
}

func TestMenuList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := menu.List(context.Background(), db)
	if err != nil {
		t.Fatalf("listing menus: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected menu list size %v, got %v", exp, got)
	}
}
