package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// структура для хранения в двусвзном листе ключа и значения
// извлечение ключа из этой структуры нужно для удаления из мапы для кэша за O(1).
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
	c.mu.Lock()
	defer c.mu.Unlock()

	elem := cacheItem{key, value}

	// Если элемент уже присутствует в кэше, то
	// меняем значение, помещаем в начало списка и возвращаем true
	if litem, ok := c.items[key]; ok {
		litem.Value = elem
		c.queue.MoveToFront(litem)
		return true
	}

	// иначе добавляем элемент в начало очереди
	litem := c.queue.PushFront(elem)
	c.items[key] = litem

	// и если размер очереди превышает ёмкость, то
	// удаляем последний элемент из очереди и мапы
	if c.queue.Len() > c.capacity {
		lastListItem := c.queue.Back()
		c.queue.Remove(lastListItem)
		delete(c.items, lastListItem.Value.(cacheItem).key)
	}

	// и возвращаем false, в знак того что ключа не было в кэше
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если элемент присутствует в кэше
	if litem, ok := c.items[key]; ok {
		c.queue.MoveToFront(litem)
		return litem.Value.(cacheItem).value, true
	}

	// Если элемент отсутствует в кэше
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
