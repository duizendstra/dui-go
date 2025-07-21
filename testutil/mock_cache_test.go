package testutil

import (
	"testing"
)

func TestMockCache(t *testing.T) {
	mc := NewMockCache()

	mc.Set("foo", "bar")
	val, ok := mc.Get("foo")
	if !ok || val != "bar" {
		t.Errorf("expected 'bar', got %v", val)
	}

	mc.SetAll(map[string]interface{}{"a": 1, "b": 2})
	all := mc.GetAll()
	if len(all) != 3 {
		t.Errorf("expected 3 items, got %d", len(all))
	}

	mc.Flush()
	all = mc.GetAll()
	if len(all) != 0 {
		t.Error("expected no items after flush")
	}
}
