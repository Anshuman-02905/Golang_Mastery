// PACKAGE FOR CONSISTENT HASHING
package consistenthash
// so these virtual nodes have a range of  data with index as sorted hash .
//  when queuing if the code has >= to queried data then that closest node has the data

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	Hash     Hash
	replicas int
	Keys     []int //sorted
	HashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	//intialise only 3 things
	m := &Map{
		replicas: replicas,
		Hash:     fn,
		Keys:     make([]int, 0),
		HashMap:  make(map[int]string),
	}
	if m.Hash == nil {
		m.Hash = crc32.ChecksumIEEE
	}
	return m
}

// isEmpty returns true if there  are no items not available
func (m *Map) IsEmpty() bool {
	return len(m.Keys) == 0
}

// Add adds some keys to the hash. Creates Virtual nodes in a ring like structure
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			fmt.Println(strconv.Itoa(i) + key)
			hash := int(m.Hash([]byte(strconv.Itoa(i) + key)))
			m.Keys = append(m.Keys, hash)
			m.HashMap[hash] = key
		}
	}
	sort.Ints(m.Keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.Hash([]byte(key)))

	idx := sort.Search(len(m.Keys), func(i int) bool { return m.Keys[i] >= hash })

	if idx == len(m.Keys) {
		idx = 0
	}
	return m.HashMap[m.Keys[idx]]
}
