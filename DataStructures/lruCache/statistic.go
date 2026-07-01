package main

type CacheStats struct {
	Hits      int
	Misses    int
	Evictions int
}

func (c *LRUCache) Stats() CacheStats {
	// TODO
	return c.statistic
}

func (c *LRUCache) HitRate() float64 {
	// TODO: вычислите процент попаданий
	total := c.statistic.Hits + c.statistic.Misses
	if total == 0 {
		return 0.0
	}
	hitRate := float64(c.statistic.Hits) / float64(total) * 100.0
	return hitRate
}
