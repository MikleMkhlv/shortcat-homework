package main

import (
	"fmt"
	"testing"
)

func setupCacheeWithEntry(t *testing.T, cap int) *LRUCache {
	t.Helper()
	cache := NewLRUCache(cap)

	testEntrys := []struct {
		key   string
		value any
	}{
		{"test-1", "value_from_test1"},
		{"test-2", "value_from_test2"},
		{"test-3", "value_from_test3"},
		{"test-4", "value_from_test4"},
		{"test-5", "value_from_test5"},
	}

	for _, val := range testEntrys {
		cache.Put(val.key, val.value)
	}
	return cache
}

func TestCorrectPutNewEntryInCacheWithoutDisplacement(t *testing.T) {
	cache := setupCacheeWithEntry(t, 10)

	expectedEntry := CacheEntry{
		Key:   "test-6",
		Value: "value_from_test6",
	}

	cache.Put(expectedEntry.Key, expectedEntry.Value)

	if !cache.Contains(expectedEntry.Key) {
		t.Errorf("%v: entry id=%s not contains in cache", t.Name(), expectedEntry.Key)
	}

	if cache.statistic.Evictions != 0 {
		t.Errorf("%v: premature displacement. got %d, wait %d", t.Name(), cache.statistic.Evictions, 0)
	}
}

func TestCorrectUpdateValue(t *testing.T) {
	cache := setupCacheeWithEntry(t, 10)
	expectedEntryWithNewValue := CacheEntry{
		Key:   "test-6",
		Value: "new_value_from_test6",
	}

	cache.Put(expectedEntryWithNewValue.Key, expectedEntryWithNewValue.Value)

	updatedEntry, exists := cache.Get(expectedEntryWithNewValue.Key)
	if exists {
		if updatedEntry.(*CacheEntry).Value != expectedEntryWithNewValue.Value {
			t.Errorf("%v: failed update value in id=%s. got value=%v, wait value=%v", t.Name(), updatedEntry.(*CacheEntry).Value, expectedEntryWithNewValue.Value, expectedEntryWithNewValue.Key)
		}
	} else {
		t.Errorf("%v: entry id=%s is not exist", t.Name(), expectedEntryWithNewValue.Key)
	}
}

func TestCheckcorrectIncrementUsesCount(t *testing.T) {
	cache := setupCacheeWithEntry(t, 10)
	expectedEntry := CacheEntry{
		Key:       "test-5",
		UsesCount: 8,
	}
	countCalls := 8

	var currentCountCalls int
	for i := 0; i < countCalls; i++ {
		entry, exists := cache.Get(expectedEntry.Key)
		if !exists {
			t.Errorf("%v: entry id=%s is not exist", t.Name(), expectedEntry.Key)
		}
		currentCountCalls = entry.(*CacheEntry).UsesCount
	}

	if currentCountCalls != expectedEntry.UsesCount {
		t.Errorf("%v: The number of calls is calculated incorrectly (cache entry). got %d, wait %d", t.Name(), currentCountCalls, expectedEntry.UsesCount)
	}
}

func TestCorrectPutNewEntryInCacheWithDisplacement(t *testing.T) {
	cache := setupCacheeWithEntry(t, 5)
	expectedEntry := CacheEntry{
		Key: "test-3",
	}

	usesCountInEntryes := make(map[string]int, cache.Len())

	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("test-%d", i)
		if key == expectedEntry.Key {
			continue
		}
		usesCountInEntryes[key] = 0
		for j := 0; j < i; j++ {
			entry, exsist := cache.Get(key)
			if !exsist {
				t.Fatalf("%v: entry id=%s is not exist", t.Name(), expectedEntry.Key)
			}
			usesCountInEntryes[key] = entry.(*CacheEntry).UsesCount
		}
	}

	newEntry := CacheEntry{
		Key:   "test-6",
		Value: "value_from_test6",
	}

	cache.Put(newEntry.Key, newEntry.Value)
	if exsist := cache.Contains(newEntry.Key); !exsist {
		t.Errorf("%v: error put new entry id=%s", t.Name(), newEntry.Key)
	}

	if exsist := cache.Contains(expectedEntry.Key); exsist {
		t.Errorf("%v: error displacement old entry id=%s", t.Name(), newEntry.Key)
	}
}

func TestDeleteEntryById(t *testing.T) {
	cache := setupCacheeWithEntry(t, 5)
	expectedEntry := CacheEntry{
		Key: "test-3",
	}

	if deleted := cache.Delete(expectedEntry.Key); !deleted {
		t.Errorf("%v: error delete entry id=%s", t.Name(), expectedEntry.Key)
	}
}
