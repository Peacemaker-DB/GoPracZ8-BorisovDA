package notes_test

import (
	"context"
	"testing"

	"example.com/gopracz8-borisovda/internal/db"
	"example.com/gopracz8-borisovda/internal/notes"
)

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()
	deps, err := db.ConnectMongo(ctx, "mongodb://root:secret@localhost:27017/?authSource=admin", "pz8_test1")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})
	r, err := notes.NewRepo(deps.Database)
	if err != nil {
		t.Fatal(err)
	}

	created, err := r.Create(ctx, "T1", "C1")
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "T1" {
		t.Fatalf("want T1 got %s", got.Title)
	}

	err = r.Delete(ctx, created.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Error("Запись удалена")
	}
}
