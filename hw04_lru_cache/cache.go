package hw04lrucache

import "sync"

var lock sync.Mutex

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Clear() {
	lock.Lock()
	defer lock.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	lock.Lock()
	defer lock.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)
	ci, ok := (item.Value).(cacheItem)
	if !ok {
		return nil, false
	}

	return ci.value, true
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	lock.Lock()
	defer lock.Unlock()

	ci := cacheItem{key: key, value: value}

	if item, ok := c.items[key]; ok {
		item.Value = ci
		c.queue.MoveToFront(item)

		return true
	}

	if c.queue.Len() == c.capacity {
		ci, ok := (c.queue.Back().Value).(cacheItem)
		if !ok {
			panic("Cache item is not cacheItem")
		}
		delete(c.items, ci.key)
		c.queue.Remove(c.queue.Back())
	}

	c.items[key] = c.queue.PushFront(ci)

	return false
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
