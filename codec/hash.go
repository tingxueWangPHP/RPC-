package codec

import (
	"fmt"
	"sort"
	"hash/crc32"
)

type UniformityHash struct {
	serversHash 	[]int
	serversRelation map[int]string
	virtualNum 		int
	hashFunc 		func ([]byte) uint32
}

func NewHash(virtualNum int, hashFunc func ([]byte) uint32) *UniformityHash {
	obj := &UniformityHash{
		//serversRelation:make(map[int]string),
		virtualNum:virtualNum,
		hashFunc:hashFunc,
	}

	if hashFunc == nil {
		obj.hashFunc = crc32.ChecksumIEEE
	}

	return obj
}

func (u *UniformityHash) initHash() {
	u.serversHash = []int{}
	u.serversRelation = map[int]string{}
}

func (u *UniformityHash) addHash(servers ...string) {
	u.initHash()
	var hash int
	for _, v := range servers {
		for i:=0;i<u.virtualNum;i++ {
			hash = int(u.hashFunc([]byte(fmt.Sprintf("%s%d", v, i))))
			u.serversHash = append(u.serversHash, hash)
			u.serversRelation[hash] = v
		}
	}

	sort.Ints(u.serversHash)
}

func (u *UniformityHash) getHash(key string) string {
	hash := int(u.hashFunc([]byte(key)))

	l := len(u.serversHash)

	idx := sort.Search(l, func(i int) bool {
		return u.serversHash[i] >= hash
	})

	return u.serversRelation[u.serversHash[idx%l]]
}