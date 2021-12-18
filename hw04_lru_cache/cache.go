package hw04lrucache

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

// type cacheItem struct {
// 	key   Key
// 	value interface{}
// }

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if v, ok := c.items[key]; ok {
		c.items[key].Value = value
		c.queue.MoveToFront(v)
		return true
	}
	newItem := c.queue.PushFront(value)
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		return v.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem)
}
