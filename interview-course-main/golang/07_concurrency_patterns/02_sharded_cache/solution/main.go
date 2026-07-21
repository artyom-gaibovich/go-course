package main

import (
	"hash/fnv"
	"sync"
)

// ShardedCache - кэш с шардированием для уменьшения блокировок.
type ShardedCache struct {
	shards []*shard
	count  int
}

type shard struct {
	mu    sync.RWMutex
	items map[string]string
}

// NewShardedCache создает новый шардированный кэш.
func NewShardedCache(shardCount int) *ShardedCache {
	shards := make([]*shard, shardCount)
	for i := range shards {
		shards[i] = &shard{
			items: make(map[string]string),
		}
	}
	return &ShardedCache{
		shards: shards,
		count:  shardCount,
	}
}

// getShard возвращает шард для ключа на основе хеша.
func (c *ShardedCache) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c.shards[h.Sum32()%uint32(c.count)]
}

// Set устанавливает значение в кэш.
func (c *ShardedCache) Set(k string, v string) {
	sh := c.getShard(k)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.items[k] = v
}

// Get получает значение из кэша.
func (c *ShardedCache) Get(k string) (string, bool) {
	sh := c.getShard(k)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	val, ok := sh.items[k]
	return val, ok
}

// Объяснение:
// 1. Шардирование разделяет кэш на несколько частей (шардов).
// 2. Каждый шард имеет свой мьютекс, что уменьшает конкуренцию.
// 3. Ключ распределяется по шардам на основе хеша.
// 4. Это позволяет параллельно обрабатывать операции на разных шардах.
// 5. Улучшает производительность при высокой конкурентности.
