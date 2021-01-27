package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/lwch/hashmap"
)

type node struct {
	k        string
	v        string
	deadline time.Time
}

type stringSlice struct {
	data []node
	size uint64
}

func (s *stringSlice) Make(size uint64) {
	s.data = make([]node, size)
}

func (s *stringSlice) Resize(size uint64) {
	data := make([]node, size)
	copy(data, s.data)
	s.data = data
}

func (s *stringSlice) Size() uint64 {
	return s.size
}

func (s *stringSlice) Cap() uint64 {
	return uint64(len(s.data))
}

func (s *stringSlice) Hash(key interface{}) uint64 {
	sum := md5.Sum([]byte(key.(string)))
	a := binary.BigEndian.Uint64(sum[:])
	b := binary.BigEndian.Uint64(sum[8:])
	return a + b
}

func (s *stringSlice) KeyEqual(idx uint64, key interface{}) bool {
	data := s.data[int(idx)%len(s.data)]
	return data.k == key.(string)
}

func (s *stringSlice) Empty(idx uint64) bool {
	data := s.data[int(idx)%len(s.data)]
	return len(data.k) == 0 && len(data.v) == 0
}

func (s *stringSlice) Set(idx uint64, key, value interface{}, deadline time.Time, update bool) {
	data := &s.data[int(idx)%len(s.data)]
	data.k = key.(string)
	data.v = value.(string)
	data.deadline = deadline
	if !update {
		s.size++
	}
}

func (s *stringSlice) Get(idx uint64) interface{} {
	data := s.data[int(idx)%len(s.data)]
	return data.v
}

func (s *stringSlice) Reset(idx uint64) {
	data := &s.data[int(idx)%len(s.data)]
	data.k = ""
	data.v = ""
	data.deadline = time.Unix(0, 0)
	s.size--
}

func (s *stringSlice) Timeout(idx uint64) bool {
	return time.Now().After(s.data[int(idx)%len(s.data)].deadline)
}

func main() {
	s := &stringSlice{}
	mp := hashmap.New(s, 1000, 5, 1000, 10*time.Second)
	for i := 0; i < 100; i++ {
		str := fmt.Sprintf("%d", i)
		mp.Set(str, str)
	}
	fmt.Printf("5=%v\n", mp.Get("5"))
	fmt.Printf("size=%d\n", mp.Size())
	time.Sleep(time.Minute)
	fmt.Printf("size=%d\n", mp.Size())
}
