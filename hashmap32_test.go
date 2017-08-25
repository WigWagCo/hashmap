package hashmap

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"unsafe"
)

// type Animal struct {
// 	name string
// }

func TestMapCreation32(t *testing.T) {
	m := &HashMap32{}
	if m.Len() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestOverwrite32(t *testing.T) {
	m := &HashMap32{}

	elephant := "elephant"
	monkey := "monkey"

	m.Set(1, unsafe.Pointer(&elephant))
	m.Set(1, unsafe.Pointer(&monkey))

	if m.Len() != 1 {
		t.Error("map should contain exactly one element.")
	}

	tmp, ok := m.Get(1) // Retrieve inserted element.
	if !ok {
		t.Error("ok should be true for item stored within the map.")
	}

	item := (*string)(tmp) // Type assertion.
	if *item != "monkey" {
		t.Error("wrong item returned.")
	}
}

func TestInsert32(t *testing.T) {
	m := &HashMap32{}
	c := uint32(16)
	ok := m.Insert(uint32(128), unsafe.Pointer(&c))
	if !ok {
		t.Error("insert did not succeed.")
	}
	ok = m.Insert(uint32(128), unsafe.Pointer(&c))
	if ok {
		t.Error("insert on existing item did succeed.")
	}
	_, ok = m.GetUintKey(128)
	if !ok {
		t.Error("ok should be true for item stored within the map.")
	}
	if m.Len() != 1 {
		t.Error("map should contain exactly 1 element.")
	}
}

func TestSet32(t *testing.T) {
	m := New(4)
	elephant := "elephant"
	monkey := "monkey"

	m.Set(4, unsafe.Pointer(&elephant))
	m.Set(3, unsafe.Pointer(&elephant))
	m.Set(2, unsafe.Pointer(&monkey))
	m.Set(1, unsafe.Pointer(&monkey))

	if m.Len() != 4 {
		t.Error("map should contain exactly 4 elements.")
	}
}

func TestGet32(t *testing.T) {
	m := &HashMap32{}
	elephant := "elephant"

	val, ok := m.Get("animal") // Get a missing element.
	if ok {
		t.Error("ok should be false when item is missing from map.")
	}
	if val != nil {
		t.Error("Missing values should return as nil.")
	}

	m.Set("animal", unsafe.Pointer(&elephant))

	_, ok = m.Get("human") // Get a missing element.
	if ok {
		t.Error("ok should be false when item is missing from map.")
	}

	val, ok = m.Get("animal") // Retrieve inserted element.
	if !ok {
		t.Error("ok should be true for item stored within the map.")
	}

	animal := (*string)(val)
	if *animal != elephant {
		t.Error("item was modified.")
	}
}

func TestResize32(t *testing.T) {
	m := New(2)
	itemCount := 50

	for i := 0; i < itemCount; i++ {
		m.Set(uint64(i), unsafe.Pointer(&Animal{strconv.Itoa(i)}))
	}

	if m.Len() != uint64(itemCount) {
		t.Error("Expected element count did not match.")
	}

	if m.Fillrate() < 30 {
		t.Error("Expecting >= 30 percent fillrate.")
	}

	for i := 0; i < itemCount; i++ {
		_, ok := m.Get(uint64(i))
		if !ok {
			t.Error("Getting inserted item failed.")
		}
	}
}

func TestStringer32(t *testing.T) {
	m := &HashMap32{}
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	s := m.String()
	if s != "[]" {
		t.Error("empty map as string does not match.")
	}

	m.Set(0, unsafe.Pointer(elephant))
	s = m.String()
	if s != "[148298089]" {
		t.Error("1 item map as string does not match:", s)
	}

	m.Set(1, unsafe.Pointer(monkey))
	s = m.String()
	if s != "[148298089,4089149075]" {
		t.Error("2 item map as string does not match:", s)
	}
}

func TestDelete32(t *testing.T) {
	m := &HashMap32{}
	m.Del(0)

	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}
	m.Set(1, unsafe.Pointer(elephant))
	m.Set(2, unsafe.Pointer(monkey))
	m.Del(0)
	m.Del(3)
	if m.Len() != 2 {
		t.Error("map should contain exactly two elements.")
	}

	m.Del(1)
	m.Del(1)
	m.Del(2)
	if m.Len() != 0 {
		t.Error("map should be empty.")
	}

	val, ok := m.Get(1) // Get a missing element.
	if ok {
		t.Error("ok should be false when item is missing from map.")
	}
	if val != nil {
		t.Error("Missing values should return as nil.")
	}

	m.Set(1, unsafe.Pointer(elephant))
}

func TestIterator32(t *testing.T) {
	m := &HashMap32{}
	itemCount := 16

	for i := itemCount; i > 0; i-- {
		m.Set(uint64(i), unsafe.Pointer(&Animal{strconv.Itoa(i)}))
	}

	counter := 0
	for item := range m.Iter() {
		val := item.Value
		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != itemCount {
		t.Error("Returned item count did not match.")
	}
}

func TestHashedKey32(t *testing.T) {
	m := &HashMap32{}
	itemCount := 16
	log := log2_32(uint32(itemCount))

	for i := 0; i < itemCount; i++ {
		m.SetHashedKey(uint32(i)<<(32-log), unsafe.Pointer(&Animal{strconv.Itoa(i)}))
	}

	if m.Len() != uint32(itemCount) {
		t.Error("Expected element count did not match.")
	}

	for i := 0; i < itemCount; i++ {
		_, ok := m.GetHashedKey(uint32(i) << (32 - log))
		if !ok {
			t.Error("Getting inserted item failed.")
		}
	}

	for i := 0; i < itemCount; i++ {
		m.DelHashedKey(uint32(i) << (32 - log))
	}

	if m.Len() != uint32(0) {
		t.Error("Map is not empty.")
	}
}

func TestCompareAndSwapHashedKey32(t *testing.T) {
	m := &HashMap32{}
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	m.SetHashedKey(1<<30, unsafe.Pointer(elephant))
	if m.Len() != 1 {
		t.Error("map should contain exactly one element.")
	}
	if !m.CasHashedKey(1<<30, unsafe.Pointer(elephant), unsafe.Pointer(monkey)) {
		t.Error("Cas should success if expectation met")
	}
	if m.CasHashedKey(1<<30, unsafe.Pointer(elephant), unsafe.Pointer(monkey)) {
		t.Error("Cas should fail if expectation didn't meet")
	}
	tmp, ok := m.GetHashedKey(1 << 30)
	if !ok {
		t.Error("ok should be true for item stored within the map.")
	}
	item := (*Animal)(tmp)
	if item != monkey {
		t.Error("wrong item returned.")
	}
}

func TestCompareAndSwap32(t *testing.T) {
	m := &HashMap32{}
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	m.Set("animal", unsafe.Pointer(elephant))
	if m.Len() != 1 {
		t.Error("map should contain exactly one element.")
	}
	if !m.Cas("animal", unsafe.Pointer(elephant), unsafe.Pointer(monkey)) {
		t.Error("Cas should success if expectation met")
	}
	if m.Cas("animal", unsafe.Pointer(elephant), unsafe.Pointer(monkey)) {
		t.Error("Cas should fail if expectation didn't meet")
	}
	tmp, ok := m.Get("animal")
	if !ok {
		t.Error("ok should be true for item stored within the map.")
	}
	item := (*Animal)(tmp)
	if item != monkey {
		t.Error("wrong item returned.")
	}
}

// TestAPICounter shows how to use the hashmap to count REST server API calls
func TestAPICounter32(t *testing.T) {
	m := &HashMap32{}

	for i := 0; i < 100; i++ {
		s := fmt.Sprintf("/api%d/", i%4)

		for {
			val, ok := m.GetStringKey(s)
			if !ok { // item does not exist yet
				c := int64(1)
				if !m.Insert(s, unsafe.Pointer(&c)) {
					continue // item was inserted concurrently, try to read it again
				}
				break
			}

			c := (*int64)(val)
			atomic.AddInt64(c, 1)
			break
		}
	}

	s := fmt.Sprintf("/api%d/", 0)
	val, _ := m.GetStringKey(s)
	c := (*int64)(val)
	if *c != 25 {
		t.Error("wrong API call count.")
	}
}

func TestExample32(t *testing.T) {
	m := &HashMap32{}
	i := 123
	m.Set("amount", unsafe.Pointer(&i))

	val, ok := m.Get("amount")
	if !ok {
		t.Fail()
	}

	j := *(*int)(val)
	if i != j {
		t.Fail()
	}
}

func TestGetOrInsert32(t *testing.T) {
	m := &HashMap32{}

	var i, j int64
	actual, loaded := m.GetOrInsert("api1", unsafe.Pointer(&i))
	if loaded {
		t.Error("item should have been inserted.")
	}

	counter := (*int64)(actual)
	if *counter != 0 {
		t.Error("item should be 0.")
	}

	atomic.AddInt64(counter, 1) // increase counter

	actual, loaded = m.GetOrInsert("api1", unsafe.Pointer(&j))
	if !loaded {
		t.Error("item should have been loaded.")
	}

	counter = (*int64)(actual)
	if *counter != 1 {
		t.Error("item should be 1.")
	}
}