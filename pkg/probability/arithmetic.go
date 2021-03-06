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

type arithmeticOperation struct {
	a, b    Uncertain
	i       int
	combine combineFunc
}

type combineFunc func(x float64, y float64) float64

func Add(a Uncertain, b Uncertain) Uncertain {
	return newArithmetic(a, b, func(x, y float64) float64 {
		return x + y
	})
}

func Sub(a Uncertain, b Uncertain) Uncertain {
	return newArithmetic(a, b, func(x, y float64) float64 {
		return x - y
	})
}

func Mul(a Uncertain, b Uncertain) Uncertain {
	return newArithmetic(a, b, func(x, y float64) float64 {
		return x * y
	})
}

func Div(a Uncertain, b Uncertain) Uncertain {
	return newArithmetic(a, b, func(x, y float64) float64 {
		return x / y
	})
}

func newArithmetic(a, b Uncertain, op combineFunc) *arithmeticOperation {
	return &arithmeticOperation{
		a:       a,
		b:       b,
		combine: op,
		i:       newID(),
	}
}

func (ar *arithmeticOperation) sampleWithTrace() *sample {
	as := ar.a.sampleWithTrace()
	bs := ar.b.sampleWithTrace()
	v := ar.combine(as.value, bs.value)
	s := as.combine(bs)
	s.value = v
	s.addTrace(ar.i, v)
	return s
}

func (ar *arithmeticOperation) sample() float64 {
	a := ar.a.sample()
	b := ar.b.sample()
	return ar.combine(a, b)
}

func (ar *arithmeticOperation) id() int {
	return ar.i
}
