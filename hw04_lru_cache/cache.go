package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	sync.Mutex
	items map[Key]*ListItem
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

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	newItem := &cacheItem{key, value}
	if v, ok := c.items[key]; ok {
		v.Value = newItem
		c.queue.MoveToFront(v)
		return true
	}
	addedItem := c.queue.PushFront(newItem)
	c.items[key] = addedItem

	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		listItem := lastItem.Value.(*cacheItem)
		delete(c.items, listItem.key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		listItem := v.Value.(*cacheItem)
		return listItem.value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem)
}
