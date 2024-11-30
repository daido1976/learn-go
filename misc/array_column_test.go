package main

import (
	"reflect"
	"testing"
)

// TestArrayColumn tests the arrayColumn function.
func TestArrayColumn(t *testing.T) {
	users := []map[string]any{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
	}

	t.Run("Extract 'name' column", func(t *testing.T) {
		names := arrayColumn(users, "name")
		expected := []any{"Alice", "Bob", "Charlie"}
		if !reflect.DeepEqual(names, expected) {
			t.Errorf("got %v, want %v", names, expected)
		}
	})

	t.Run("Extract 'id' column", func(t *testing.T) {
		ids := arrayColumn(users, "id")
		expected := []any{1, 2, 3}
		if !reflect.DeepEqual(ids, expected) {
			t.Errorf("got %v, want %v", ids, expected)
		}
	})

	t.Run("Extract non-existent column", func(t *testing.T) {
		nonExistent := arrayColumn(users, "age")
		expected := []any{}
		if !reflect.DeepEqual(nonExistent, expected) {
			t.Errorf("got %v, want %v", nonExistent, expected)
		}
	})
}
