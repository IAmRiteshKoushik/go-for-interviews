package main

// Prototype interface for "node" and not a reference to index-node in Linux
type Inode interface {
	print(string)
	clone() Inode
}
