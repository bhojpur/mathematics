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

import "testing"

func TestConstantSample(t *testing.T) {
	c := NewConstant(5.0)
	m := Materialize(c, 10)
	if m.First() != 5.0 {
		t.Fatal("Couldn't retrieve constant after sampling")
	}
}

func TestConstantBNNSample(t *testing.T) {
	a := NewConstant(5.0)
	b := NewConstant(6.0)
	c := Add(a, b)

	m := Materialize(c, 10)
	if m.First() != 11.0 {
		t.Fatal("Addition on constants failed")
	}
}

func TestConstantBernoulliSample(t *testing.T) {
	a := NewConstant(5.0)
	b := NewConstant(6.0)
	c := LessThan(b, a)

	m := Materialize(c, 10)
	if m.First() != 0.0 {
		t.Fatal("LessThan on constants failed")
	}
}

func TestBernoulliConditional(t *testing.T) {
	x := NewConstant(5.0)
	y := NewConstant(6.0)

	if LessThan(y, x).Pr() {
		t.Fatal("6 tests less than 5")
	}
	if !LessThan(x, y).Pr() {
		t.Fatal("5 tests not less than 6")
	}
}

func TestBernoulliEqual(t *testing.T) {
	x := NewConstant(5.0)
	y := NewConstant(5.0)

	if LessThan(x, y).Pr() {
		t.Fatal("x is less than y, incorrectly (they should be equal)")
	}
	if LessThan(y, x).Pr() {
		t.Fatal("y is less than x, incorrectly (they should be equal)")
	}
}
