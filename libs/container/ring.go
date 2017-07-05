/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/04/02        Yusank Kurban
 */

package container

import (
	"errors"
)

var (
	ErrRingLenNotEnough = errors.New("ring has not enough items for pop n")
	ErrRingCapNotEnough = errors.New("ring has not enough space for push n")
)

type Ring struct {
	items   	[]interface{}
	head    	int
	tail    	int
	size    	int
	maxSize 	int
}

func NewRing(maxSize int) *Ring {
	r := new(Ring)

	r.size = maxSize
	r.head = 0
	r.tail = 0
	r.maxSize = r.size + 1

	r.items = make([]interface{}, r.maxSize)

	return r
}

func (r *Ring) Len() int {
	if r.head == r.tail {
		return 0

	} else if r.tail > r.head {
		return r.tail - r.head

	} else {
		return r.tail + r.maxSize - r.head
	}
}

func (r *Ring) Cap() int {
	return r.size - r.Len()
}

func (r *Ring) MPop(n int) ([]interface{}, error) {
	if r.Len() < n {
		return nil, ErrRingLenNotEnough
	}

	items := make([]interface{}, n)

	for i := 0; i < n; i++ {
		head := (r.head + i) % r.maxSize
		items[i] = r.items[head]
		r.items[head] = nil
	}

	r.head = (r.head + n) % r.maxSize

	return items, nil
}

func (r *Ring) Pop() (interface{}, error) {
	if items, err := r.MPop(1); err != nil {
		return nil, err

	} else {
		return items[0], nil
	}
}

func (r *Ring) MPush(items []interface{}) error {
	n := len(items)

	if r.Cap() < n {
		return ErrRingCapNotEnough
	}

	for i := 0; i < n; i++ {
		tail := (r.tail + i) % r.maxSize
		r.items[tail] = items[i]
	}

	r.tail = (r.tail + n) % r.maxSize

	return nil
}

func (r *Ring) Push(item interface{}) error {
	items := []interface{}{item}

	return r.MPush(items)
}

func (r *Ring) Full() bool {
	return r.Cap() == 0
}

func (r *Ring) Empty() bool {
	return r.Len() == 0
}

func (r *Ring) Gets(n int) []interface{} {
	if r.Len() < n {
		n = r.Len()
	}
	result := make([]interface{}, n)
	for i := 0; i < n; i++ {
		result[i] = r.items[(r.head+i)%r.maxSize]
	}

	return result
}

func (r *Ring) Get() interface{} {
	if r.Empty() {
		return ErrRingLenNotEnough
	}

	return r.items[r.head]
}

func (r *Ring) GetAll() []interface{} {
	return r.Gets(r.Len())
}
