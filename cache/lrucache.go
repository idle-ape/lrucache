/*
 * @Author: Bourne<kingphper@foxmail.com>
 * @Date: 2020-05-21 15:13:29
 * @LastEditTime: 2020-05-22 10:43:22
 * @LastEditors: Please set LastEditors
 * @Description: Lru Cache
 */

package cache

import (
	"sync"
	"time"
)

var l sync.RWMutex

// Entry each item store in map
type Entry struct {
	Key        string
	Value      interface{}
	Prev       *Entry
	Next       *Entry
	ExpireTime int64
}

// Cache ...
type Cache struct {
	Cache    map[string]*Entry
	Head     *Entry
	Tail     *Entry
	Capacity int
}

func NewCache(capacity int) *Cache {
	return &Cache{Cache: make(map[string]*Entry, capacity), Capacity: capacity}
}

// Get get a entry from cache
func (c *Cache) Get(key string) (value interface{}, exists bool) {

	defer l.Unlock()

	l.Lock()

	entry, exists := c.Cache[key]
	if !exists {
		return nil, false
	}

	if entry.ExpireTime > 0 && entry.ExpireTime < time.Now().Unix() {
		c.delete(entry)

		return nil, false
	}

	c.moveToHead(entry)

	return entry.Value, true
}

// Set set a entry, update the value if entry exists
func (c *Cache) Set(key string, value interface{}, expire int64) bool {

	defer l.Unlock()

	l.Lock()

	if _, exists := c.Cache[key]; !exists {
		c.put(key, value, expire)
		return true
	}

	entry := c.Cache[key]
	entry.Value = value

	c.moveToHead(entry)

	return true
}

// Setnx set a entry when it's not exists
func (c *Cache) Setnx(key string, value interface{}, expire int64) bool {

	defer l.Unlock()

	l.Lock()

	if _, exists := c.Cache[key]; exists {
		return false
	}

	c.put(key, value, expire)

	return true
}

func (c *Cache) put(key string, value interface{}, expire int64) {
	entry := &Entry{
		Key:   key,
		Value: value,
		Prev:  nil,
		Next:  c.Head,
	}
	if expire > 0 {
		entry.ExpireTime = time.Now().Unix() + expire
	}

	c.Cache[key] = entry

	if c.Head != nil {
		c.Head.Prev = entry
	}

	if c.Tail == nil {
		c.Tail = entry
	}

	c.Head = entry

	if len(c.Cache) > c.Capacity {
		c.delete(c.Tail)
	}
}

// Delete delete a entry
func (c *Cache) Delete(key string) bool {

	defer l.Unlock()

	l.Lock()

	entry, exists := c.Cache[key]

	if !exists {
		return false
	}

	c.delete(entry)

	return true
}

// delete a entry
func (c *Cache) delete(removedEntry *Entry) {
	if c.Tail == removedEntry {
		newTail := removedEntry.Prev
		if newTail != nil {
			newTail.Next = nil
		}
		c.Tail = newTail
	} else if c.Head == removedEntry {
		newHead := removedEntry.Next
		if newHead != nil {
			newHead.Prev = nil
		}
		c.Head = newHead
	} else {
		removedEntry.Prev.Next = removedEntry.Next
		removedEntry.Next.Prev = removedEntry.Prev
	}

	// delete the cache from the map
	delete(c.Cache, removedEntry.Key)
}

// move a entry to the head
func (c *Cache) moveToHead(entry *Entry) {
	if c.Head == entry {
		return
	}

	if c.Tail == entry {
		c.Tail = entry.Prev
		entry.Prev.Next = nil
	} else {
		entry.Prev.Next = entry.Next
		entry.Next.Prev = entry.Prev
	}

	c.Head.Prev = entry
	entry.Next = c.Head
	entry.Prev = nil
	c.Head = entry

	return
}
