package cache

import "testing"

func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache()

	// Test Set and Get
	cache.Set("foo", "bar")
	val, ok := cache.Get("foo")
	if !ok || val != "bar" {
		t.Fatalf("expected 'bar', got %v", val)
	}

	// Test SetAll
	cache.SetAll(map[string]interface{}{
		"a": 1,
		"b": 2,
	})
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("expected 1 for 'a', got %v", val)
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("expected 2 for 'b', got %v", val)
	}

	// Test GetAll
	all := cache.GetAll()
	if len(all) != 3 {
		t.Errorf("expected 3 items, got %d", len(all))
	}
	if all["foo"] != "bar" || all["a"] != 1 || all["b"] != 2 {
		t.Errorf("GetAll returned unexpected data: %v", all)
	}

	// Test Flush
	cache.Flush()
	if _, ok := cache.Get("foo"); ok {
		t.Error("expected no values after Flush, but got one for 'foo'")
	}
	if _, ok := cache.Get("a"); ok {
		t.Error("expected no values after Flush, but got one for 'a'")
	}
}
