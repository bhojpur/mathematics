package dataframe

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"sort"
)

// OrderedMapIntFloat64 is an ordered map[int]float64.
type OrderedMapIntFloat64 struct {
	store map[int]float64
	keys  []int // if not initialized implies that it's not ordered
}

// NewOrderedMapIntFloat64 will create a new OrderedMapIntFloat64.
// By default it will be ordered, but setting notOrdered to true will make
// it operate as a builtin map.
func NewOrderedMapIntFloat64(notOrdered ...bool) *OrderedMapIntFloat64 {
	if len(notOrdered) == 0 || !notOrdered[0] {
		// ordered
		return &OrderedMapIntFloat64{
			store: map[int]float64{},
			keys:  []int{},
		}
	}

	return &OrderedMapIntFloat64{store: map[int]float64{}}
}

// Get will return the value of the key. If it doesn't exist, it will
// return false for the second return value.
func (o *OrderedMapIntFloat64) Get(key int) (float64, bool) {
	val, exists := o.store[key]
	return val, exists
}

// Set will set a key and value pair. It will overwrite an
// existing pair if it exists already.
func (o *OrderedMapIntFloat64) Set(key int, val float64) {
	if o.keys != nil {
		if _, exists := o.store[key]; !exists {
			o.keys = append(o.keys, key)
		}
	}
	o.store[key] = val
}

// Delete will remove the key from the OrderedMapIntFloat64.
// For performance reasons, ensure the key exists beforehand when in ordered mode.
func (o *OrderedMapIntFloat64) Delete(key int) {
	if o.keys == nil {
		// unordered
		delete(o.store, key)
		return
	}

	// ordered
	delete(o.store, key)

	// Find key
	var idx *int

	for i, val := range o.keys {
		if val == key {
			idx = &[]int{i}[0]
			break
		}
	}
	if idx != nil {
		o.keys = append(o.keys[:*idx], o.keys[*idx+1:]...)
	}
}

// ValuesIterator is used to iterate through the values of OrderedMapIntFloat64.
func (o *OrderedMapIntFloat64) ValuesIterator() func() (*int, float64) {
	var keys []int

	if o.keys == nil {
		for key := range o.store {
			keys = append(keys, key)
		}
		sort.Ints(keys)
	} else {
		keys = o.keys
	}

	j := 0

	return func() (*int, float64) {
		if j > len(keys)-1 {
			return nil, 0
		}

		row := keys[j]
		j++

		return &row, o.store[row]
	}
}

// OrderedMapIntMixed is an ordered map[int]interface{}.
type OrderedMapIntMixed struct {
	store map[int]interface{}
	keys  []int
}

// NewOrderedMapIntMixed will create a new OrderedMapIntMixed.
// By default it will be ordered, but setting notOrdered to true will make
// it operate as a builtin map.
func NewOrderedMapIntMixed(notOrdered ...bool) *OrderedMapIntMixed {
	if len(notOrdered) == 0 || !notOrdered[0] {
		// ordered
		return &OrderedMapIntMixed{
			store: map[int]interface{}{},
			keys:  []int{},
		}
	}

	return &OrderedMapIntMixed{store: map[int]interface{}{}}
}

// Get will return the value of the key. If it doesn't exist, it will
// return false for the second return value.
func (o *OrderedMapIntMixed) Get(key int) (interface{}, bool) {
	val, exists := o.store[key]
	return val, exists
}

// Set will set a key and value pair. It will overwrite an
// existing pair if it exists already.
func (o *OrderedMapIntMixed) Set(key int, val interface{}) {
	if o.keys != nil {
		if _, exists := o.store[key]; !exists {
			o.keys = append(o.keys, key)
		}
	}
	o.store[key] = val
}

// Delete will remove the key from the OrderedMapIntMixed.
// For performance reasons, ensure the key exists beforehand when in ordered mode.
func (o *OrderedMapIntMixed) Delete(key int) {
	if o.keys == nil {
		// unordered
		delete(o.store, key)
		return
	}

	// ordered
	delete(o.store, key)

	// Find key
	var idx *int

	for i, val := range o.keys {
		if val == key {
			idx = &[]int{i}[0]
			break
		}
	}
	if idx != nil {
		o.keys = append(o.keys[:*idx], o.keys[*idx+1:]...)
	}
}

// ValuesIterator is used to iterate through the values of OrderedMapIntMixed.
func (o *OrderedMapIntMixed) ValuesIterator() func() (*int, interface{}) {
	var keys []int

	if o.keys == nil {
		for key := range o.store {
			keys = append(keys, key)
		}
		sort.Ints(keys)
	} else {
		keys = o.keys
	}

	j := 0

	return func() (*int, interface{}) {
		if j > len(keys)-1 {
			return nil, 0
		}

		row := keys[j]
		j++

		return &row, o.store[row]
	}
}
