package hashmap

import (
	"container/list"
	"crypto/sha256"
	"math/big"
)

type chainedHashMap struct {
	hashMapBase
	buckets []*list.List
}

func (c *chainedHashMap) Init(cap uint32) {
	c.hashMapBase.Init(cap)
	if cap <= 0 {
		c.buckets = nil
	} else {
		c.buckets = make([]*list.List, c.Cap, c.Cap)
	}
}

// func (c *chainedHashMap) Move(cap uint32) {
// 	oldBuckets := c.buckets
// 	c.Init(cap)

// }

func (c *chainedHashMap) hash(key interface{}) uint32 {
	hashValue := c.HashFunc(key, sha256.New())
	mb := big.NewInt(int64(c.Cap))
	hashValue.Mod(hashValue, mb)
	return uint32(hashValue.Uint64())
}

func (c *chainedHashMap) existInList(key interface{}, list *list.List) (*list.Element, bool) {
	for e := list.Front(); e != nil; e = e.Next() {
		if e.Value.(element).Key == key {
			return e, true
		}
	}
	return nil, false
}

func (c *chainedHashMap) Insert(key, value interface{}) {
	c.UpScale()
	hashKey := c.hash(key)
	if c.buckets[hashKey] == nil {
		c.buckets[hashKey] = list.New()
	}
	e := element{key, value}
	le, exist := c.existInList(key, c.buckets[hashKey])
	if exist {
		le.Value = e
	} else {
		c.buckets[hashKey].PushFront(e)
		c.Count++
	}
}

func (c *chainedHashMap) Get(key interface{}) (interface{}, bool) {
	if c.Count == 0 {
		return nil, false
	}
	hashKey := c.hash(key)
	if c.buckets[hashKey] == nil {
		return nil, false
	}
	le, exist := c.existInList(key, c.buckets[hashKey])
	if exist {
		return le.Value.(element).Value, true
	}
	return nil, false
}

func (c *chainedHashMap) Delete(key interface{}) {
	if c.Count == 0 {
		return
	}

	hashKey := c.hash(key)
	if c.buckets[hashKey] == nil {
		return
	}

	le, exist := c.existInList(key, c.buckets[hashKey])
	if exist {
		c.buckets[hashKey].Remove(le)
	}

	if c.buckets[hashKey].Len() == 0 {
		c.buckets[hashKey] = nil
		c.Count--
	}

	c.DownScale()
}

func newChainedHashMap() *chainedHashMap {
	h := new(chainedHashMap)
	h.hashMapBase.hashMap = h
	h.hashMapBase.scalableMap = h
	return h
}
