package card_test

import (
	"context"
	"testing"
	"time"

	"github.com/events-app/events-api2/internal/platform/auth"
	"github.com/events-app/events-api2/internal/card"
	"github.com/events-app/events-api2/internal/schema"
	"github.com/events-app/events-api2/internal/tests"
	"github.com/google/go-cmp/cmp"
)

func TestCards(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	newP := card.NewCard{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e", // This is just some random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

	p0, err := card.Create(ctx, db, claims, newP, now)
	if err != nil {
		t.Fatalf("creating card p0: %s", err)
	}

	p1, err := card.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting card p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	update := card.UpdateCard{
		Name: tests.StringPointer("Comics"),
		Cost: tests.IntPointer(25),
	}
	updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := card.Update(ctx, db, claims, p0.ID, update, updatedTime); err != nil {
		t.Fatalf("creating card p0: %s", err)
	}

	saved, err := card.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting card p0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original card
	// and change just the fields we expect then diff it with what was saved.
	want := *p0
	want.Name = "Comics"
	want.Cost = 25
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}

	if err := card.Delete(ctx, db, p0.ID); err != nil {
		t.Fatalf("deleting card: %v", err)
	}

	_, err = card.Retrieve(ctx, db, p0.ID)
	if err == nil {
		t.Fatalf("should not be able to retrieve deleted card")
	}
}

func TestCardList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := card.List(context.Background(), db)
	if err != nil {
		t.Fatalf("listing cards: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected card list size %v, got %v", exp, got)
	}
}
