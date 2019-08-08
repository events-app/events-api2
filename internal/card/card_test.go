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

	newC := card.NewCard{
		Name:     "New Card",
		Content:     "This is testing card.",
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e", // This is just some random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

	c0, err := card.Create(ctx, db, claims, newC, now)
	if err != nil {
		t.Fatalf("creating card c0: %s", err)
	}

	c1, err := card.Retrieve(ctx, db, c0.ID)
	if err != nil {
		t.Fatalf("getting card c0: %s", err)
	}

	if diff := cmp.Diff(c1, c0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	update := card.UpdateCard{
		Name: tests.StringPointer("New Card"),
		Content: tests.StringPointer("Some new content of a card"),
	}
	updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := card.Update(ctx, db, claims, c0.ID, update, updatedTime); err != nil {
		t.Fatalf("creating card c0: %s", err)
	}

	saved, err := card.Retrieve(ctx, db, c0.ID)
	if err != nil {
		t.Fatalf("getting card c0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original card
	// and change just the fields we expect then diff it with what was saved.
	want := *c0
	want.Name = "New Card"
	want.Cost = 25
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}

	if err := card.Delete(ctx, db, c0.ID); err != nil {
		t.Fatalf("deleting card: %v", err)
	}

	_, err = card.Retrieve(ctx, db, c0.ID)
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
