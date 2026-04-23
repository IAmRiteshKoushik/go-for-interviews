package main

import "fmt"

type node struct {
	key   int
	value int
	prev  *node
	next  *node
}

type lrucache struct {
	capacity int
	items    map[int]*node
	head     *node
	tail     *node
}

func new(capacity int) *lrucache {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head

	return &lrucache{
		capacity: capacity,
		items:    make(map[int]*node),
		head:     head,
		tail:     tail,
	}
}

func (c *lrucache) Get(key int) int {
	n, exist := c.items[key]
	if !exist {
		return -1
	}

	c.moveToFront(n)
	return n.value
}

func (c *lrucache) Put(key int, value int) {
	if n, exists := c.items[key]; exists {
		n.value = value
		c.moveToFront(n)
		return
	}

	newNode := &node{
		key:   key,
		value: value,
	}
	c.items[key] = newNode
	c.addToFront(newNode)

	if len(c.items) > c.capacity {
		lru := c.removeLast()
		delete(c.items, lru.key)
	}
}

func (c *lrucache) moveToFront(n *node) {
	c.removeNode(n)
	c.addToFront(n)
}

func (c *lrucache) addToFront(n *node) {
	// empty head scenario
	n.prev = c.head
	n.next = c.head.next

	// actual head pointing to new node
	c.head.next.prev = n
	c.head.next = n
}

func (c *lrucache) removeNode(n *node) {
	prevNode := n.prev
	nextNode := n.next

	prevNode.next = nextNode
	nextNode.prev = prevNode
}

func (c *lrucache) removeLast() *node {
	last := c.tail.prev
	c.removeNode(last)
	return last
}

func main() {
	cache := new(2)

	cache.Put(1, 10)
	cache.Put(2, 20)
	fmt.Println(cache.Get(1))

	cache.Put(3, 30) // eviction happens
	fmt.Println(cache.Get(2))
}
