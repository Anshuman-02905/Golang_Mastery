package main

import (
	"fmt"
)

const SIZE = 5

// // var m map[string]*Node

// // var MAX_LENGTH int
// // var DLL Queue

// // func CreateNode(val string) {
// // 	ND := &Node{
// // 		data: val,
// // 		next: nil,
// // 		prev: nil,
// // 	}
// // 	m[val] = ND
// // }
// // func RemoveNode() {
// // 	if DLL.tail == nil {
// // 		return
// // 	}
// // 	delete(m, DLL.tail.data)

// // 	if DLL.tail == DLL.head {
// // 		DLL.tail, DLL.head = nil, nil
// // 	} else {
// // 		DLL.tail = DLL.tail.prev
// // 		DLL.tail.next = nil
// // 	}

// // 	DLL.length--
// // }

// // func InsertIntoCache(val string) {
// // 	if DLL.length == MAX_LENGTH {
// // 		RemoveNode()
// // 	}

// // 	if _, exists := m[val]; !exists {
// // 		CreateNode(val)
// // 	}

// // 	ND := m[val]
// // 	if DLL.head == nil {
// // 		DLL.head = ND
// // 		DLL.tail = ND
// // 	} else {
// // 		ND.next = DLL.head
// // 		DLL.head.prev = ND
// // 		DLL.head = ND
// // 	}
// // 	DLL.length = DLL.length + 1

// // }

// // func PrintCache() {

// // 	var i int
// // 	var ND *Node
// // 	ND = DLL.head
// // 	for i = 0; i < DLL.length; i++ {
// // 		fmt.Print("-->", ND.data, "-->")
// // 		ND = ND.next
// // 	}
// // 	fmt.Println()

// // }

// func main() {
// 	MAX_LENGTH = 5
// 	m = make(map[string]*Node)

// 	DLL = Queue{
// 		head:   nil,
// 		tail:   nil,
// 		length: 0,
// 	}

// 	for {
// 		var val string
// 		fmt.Println("Enter the Value: ")
// 		fmt.Scanln(&val)
// 		if val == "exit" {
// 			break
// 		}
// 		InsertIntoCache(val)
// 		PrintCache()
// 	}

// }

type Node struct {
	data string
	next *Node
	prev *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}
type Hash map[string]*Node

type Cache struct {
	Queue Queue
	Hash  Hash
}

func NewCache() Cache {
	return Cache{
		Queue: NewQueue(),
		Hash:  Hash{},
	}
}

func NewQueue() Queue {
	head := Node{}
	tail := Node{}
	head.next = &tail
	tail.prev = &head
	return Queue{Head: &head, Tail: &tail}
}
func (c *Cache) CreateNode(val string) *Node {
	ND := &Node{
		data: val,
		next: nil,
		prev: nil,
	}
	c.Hash[val] = ND
	return ND
}

func (c *Cache) Check(str string) { //Shadowing seekha earlier str was val
	node := &Node{}

	if val, ok := c.Hash[str]; ok {

		node = c.Remove(val)

	} else {
		node = c.CreateNode(str)
	}
	c.Add(node)

}

func (c *Cache) Remove(n *Node) *Node {
	fmt.Println("Rmove : %s \n", n.data)
	prev := n.prev
	next := n.next
	prev.next = next
	next.prev = prev

	delete(c.Hash, n.data)
	return n

}
func (c *Cache) Add(n *Node) {
	fmt.Println("add: %s \n", n.data)
	tmp := c.Queue.Head.next

	c.Queue.Head.next = n
	n.prev = c.Queue.Head
	n.next = tmp
	tmp.prev = n

	c.Queue.Length++
	if c.Queue.Length > SIZE {
		c.Remove(c.Queue.Tail.prev)
	}

}

func (q *Queue) Dispaly() {
	node := q.Head.next
	fmt.Printf("%d - [ ", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.data)
		if i < q.Length-1 {
			fmt.Printf("<-->")
		}
		node = node.next
	}
	fmt.Println("]")
}
func (c *Cache) Dispaly() {
	c.Queue.Dispaly()
}

func main() {
	fmt.Println("Start Cache")
	cache := NewCache()

	for {
		var val string
		fmt.Println("Enter the Value ")
		fmt.Scanln(&val)

		if val == "exit" {
			break
		}
		cache.Check(val)
		cache.Dispaly()
	}
}
