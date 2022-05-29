package probability

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
	"fmt"
	"sync"
)

type Uncertain interface {
	sample() float64
	sampleWithTrace() *sample
	id() int
}

type UncertainBool interface {
	Uncertain
	sampleBool() bool
	Pr() bool
}

type sample struct {
	value float64
	trace map[int]float64
}

func newSample(val float64) *sample {
	return &sample{
		value: val,
		trace: make(map[int]float64),
	}
}

func (s *sample) addTrace(id int, val float64) {
	s.trace[id] = val
}

func (s *sample) combine(other *sample) *sample {
	out := newSample(s.value)
	for k, v := range s.trace {
		out.addTrace(k, v)
	}
	for k, v := range other.trace {
		out.addTrace(k, v)
	}
	return out
}

func (s *sample) String() string {
	return fmt.Sprintf("%0.4f : %#v", s.value, s.trace)
}

var (
	idPrinter     int
	idPrinterLock sync.Mutex
)

func newID() int {
	idPrinterLock.Lock()
	defer idPrinterLock.Unlock()
	idPrinter++
	return idPrinter
}
