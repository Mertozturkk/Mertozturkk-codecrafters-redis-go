package storage

import (
	"testing"
	"time"
)

func TestStore_SetAndGet(t *testing.T) {
	store := NewStore("initialKey", "initialValue")

	err := store.Set("testKey", "testValue", 10*time.Second)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	value, ok := store.Get("testKey")
	if !ok || value != "testValue" {
		t.Errorf("Expected 'testValue', got '%v'", value)
	}

}

func TestStore_SetWithEmptyKey(t *testing.T) {
	store := NewStore("initialKey", "initialValue")

	err := store.Set("", "testValue", 1*time.Second)
	if err == nil {
		t.Errorf("Expected error for empty key")
	}
}

func TestStore_GetNonExistentKey(t *testing.T) {
	store := NewStore("initialKey", "initialValue")

	_, ok := store.Get("nonExistentKey")
	if ok {
		t.Errorf("Expected false for non existent key")
	}
}
