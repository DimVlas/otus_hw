package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша.
	Clear()                              // Очистить кэш.
}

type cacheItem struct {
	Value interface{}
	Key   Key
}

type lruCache struct {
	capacity int
	mtx      sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

// Создать новый LRU-кэш.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		mtx:      sync.Mutex{},
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Добавить значение в кэш.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	itm, ok := c.items[key]
	if ok { // ключ есть в кэше
		cacheVal := itm.Value.(*cacheItem)
		cacheVal.Value = value
		itm.Value = cacheVal
		c.queue.MoveToFront(itm)
	} else { // ключа нет в кэше
		if c.capacity == c.queue.Len() { // кэш полностью заполнен
			excessVal := c.queue.Back().Value.(*cacheItem)
			delete(c.items, excessVal.Key)

			c.queue.Remove(c.queue.Back())
		}

		cacheVal := &cacheItem{Key: key, Value: value}
		itm = c.queue.PushFront(cacheVal)
		c.items[key] = itm
	}

	return ok
}

// Получить значение из кэша.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.mtx.Unlock()
	c.mtx.Lock()

	itm, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(itm)

		return itm.Value.(*cacheItem).Value, true
	}

	return nil, false
}

// Очистить кэш.
func (c *lruCache) Clear() {
	c.mtx.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.mtx.Unlock()
}
