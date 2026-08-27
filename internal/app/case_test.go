package app

import (
	"context"
	"errors"
	"testing"

	"github.com/lacsar712/lyophdry/internal/config"
	"github.com/lacsar712/lyophdry/internal/model"
)

func TestCase(t *testing.T) {
	a, err := New(config.Default())
	if err != nil {
		t.Fatal(err)
	}
	err = a.ValidateMoldDrift(context.Background(), 25.0)
	if err == nil {
		t.Fatal("expected moisture drift violation")
	}
	if !errors.Is(err, model.ErrShelfDrift) {
		t.Fatalf("expected ErrShelfDrift, got %v", err)
	}
}
