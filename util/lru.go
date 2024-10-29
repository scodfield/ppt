package util

import (
	"sync"
)

type DLNode struct {
	key       string
	value     int
	pre, next *DLNode
}

type Lru struct {
	max        int
	size       int
	cache      map[string]*DLNode
	head, tail *DLNode
}

func (lru *Lru) Get(key string) int {
	node, ok := lru.cache[key]
	if !ok {
		return -1
	}
	lru.moveToHead(node)
	return node.value
}

func (lru *Lru) Set(key string, value int) {
	node, ok := lru.cache[key]
	if !ok {
		node := &DLNode{
			key:   key,
			value: value,
		}
		lru.cache[key] = node
		lru.moveToHead(node)
		lru.size++
		if lru.size > lru.max {
			lru.removeNode(lru.tail.pre)
		}
	} else {
		node.value = value
		lru.moveToHead(node)
	}
}

func (lru *Lru) moveToHead(node *DLNode) {
	node.next = lru.head.next
	node.pre = lru.head
	lru.head.next.pre = node
	lru.head.next = node
}

func (lru *Lru) removeNode(node *DLNode) {
	node.pre.next = node.next
	node.next.pre = node.pre
	delete(lru.cache, node.key)
	lru.size--
}

func NewLru(max int) *Lru {
	lru := &Lru{
		max:   max,
		cache: make(map[string]*DLNode),
		head: &DLNode{
			key: "head",
		},
		tail: &DLNode{
			key: "tail",
		},
	}
	lru.head.next = lru.tail
	lru.tail.pre = lru.head
	return lru
}

// 参考：orcaman/concurrent-map
// 期望解决的问题：避免对整个map对象加锁，导致高并发下的性能下降
// 思路：拆分map为map切片，每次只对某一map进行加锁，目的是减小锁的颗粒度=分段锁
var SHARD_COUNT = 32

type ShardLru []*ShardLruNode

type ShardLruNode struct {
	sync.RWMutex
	lru *Lru
}

func NewShardLru(max int) ShardLru {
	lru := make(ShardLru, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		lru[i].lru = NewLru(max)
	}
	return lru
}

// 获取hash后的分片
func (shard ShardLru) GetShard(key string) *ShardLruNode {
	return shard[uint(fnv32(key))%uint(SHARD_COUNT)]
}

func (shard ShardLru) Get(key string) int {
	s := shard.GetShard(key)
	s.RLock()
	defer s.RLock()
	return s.lru.Get(key)
}

func (shard ShardLru) Set(key string, value int) {
	s := shard.GetShard(key)
	s.Lock()
	defer s.Unlock()
	s.lru.Set(key, value)
}

// 分片定位时，常用的hash算法：fnv32，bkdr
// hash中间值对素数取模，目标散列函数要均匀，尽可能少的减少冲突，若对合数取模，则该散列函数对合数的因子倍数冲突
// 的概率会增大，比如 mod = 6，对于2的倍数 2/4/6/8/10/12 的hash值是 2/4/0/2/4/0，而质数的因子只有 1 和它本身
// 因此对质数取模会得到更好的散列效果
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32        // 循环乘以素数 prime32
		hash ^= uint32(key[i]) // 按位异或
	}
	return hash
}
