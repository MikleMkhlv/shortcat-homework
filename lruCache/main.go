package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"sync"
)

type CacheEntry struct {
	Key       string
	Value     interface{}
	UsesCount int
}

type LRUCache struct {
	capacity  int
	cache     map[string]*CacheEntry
	statistic CacheStats
	mx        sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:  capacity,
		cache:     make(map[string]*CacheEntry, capacity+1),
		statistic: CacheStats{},
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mx.Lock()
	if val, ok := c.cache[key]; ok {
		val.UsesCount++
		c.statistic.Hits++
		return val, true
	}
	c.mx.Unlock()
	c.statistic.Misses++

	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	// TODO: добавить/обновить элемент
	// TODO: при переполнении удалить наименее использованный
	c.mx.Lock()
	defer c.mx.Unlock()
	if val, ok := c.cache[key]; ok {
		val.Value = value
		return
	} else {
		if c.Len() >= c.capacity {
			entries := slices.Collect(maps.Values(c.cache))
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].UsesCount < entries[j].UsesCount
			})
			keyDeletedCache := entries[0].Key
			delete(c.cache, keyDeletedCache)
			c.capacity--
			c.statistic.Evictions++

			c.cache[key] = &CacheEntry{
				Key:       key,
				Value:     value,
				UsesCount: 0,
			}
		}
	}
}

func (c *LRUCache) Len() int {
	return len(c.cache)
}

func (c *LRUCache) Keys() []string {
	c.mx.Lock()
	defer c.mx.Unlock()
	keys := slices.Collect(maps.Keys(c.cache))
	return keys
}

func (c *LRUCache) Clear() {
	c.cache = make(map[string]*CacheEntry)
}

func (c *LRUCache) Delete(key string) bool {
	delete(c.cache, key)
	if _, ok := c.cache[key]; !ok {
		return true
	}
	return false
}

func (c *LRUCache) Contains(key string) bool {
	// TODO: проверить без обновления статистики использования
	if _, ok := c.cache[key]; ok {
		return true
	}
	return false
}

func (c *LRUCache) Clone() *LRUCache {
	copy := &LRUCache{
		capacity:  c.capacity,
		cache:     maps.Clone(c.cache),
		statistic: c.statistic,
	}
	return copy
}
func main() {
	cache := NewLRUCache(3)

	cache.Put("a", "one")
	cache.Put("b", "two")
	cache.Put("c", "three")

	fmt.Println(cache.Get("a")) // "one", true - теперь "a" недавно использован

	cache.Put("d", "four") // Должно вытеснить "b" (наименее использованный)

	fmt.Println(cache.Get("b")) // nil, false
	fmt.Println(cache.Get("a")) // "one", true
}
