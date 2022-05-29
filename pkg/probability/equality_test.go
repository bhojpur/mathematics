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

func TestGaussianEquality(t *testing.T) {
	x := NewGaussian(1.0, 1.0)
	y := NewGaussian(4.0, 2.0)

	if GreaterThan(x, y).Pr() {
		t.Error("x > y")
	}
	if LessThan(y, x).Pr() {
		t.Error("y < x")
	}
	if !GreaterThan(y, x).Pr() {
		t.Error("!y > x")
	}
	if !LessThan(x, y).Pr() {
		t.Error("!x < y")
	}
}

func TestNotEquals(t *testing.T) {
	coinA := NewBernoulli(0.7)
	coinB := NewBernoulli(0.5)

	z := NotEquals(coinA, coinB)
	avg := Materialize(z, 10000).Average()

	// If it's not true 50% of the time...
	if !Within(avg, 0.5, epsilon) {
		t.Error("Coin bias doesn't cancel")
	}
}

func TestMontyHall(t *testing.T) {
	carInDoor := NewEvenMultinomial([]float64{1.0, 2.0, 3.0})

	chosenDoor := NewEvenMultinomial([]float64{1.0, 2.0, 3.0})
	match := Equals(carInDoor, chosenDoor)
	// Now monty opens a door
	switchWins := Not(match)

	v := Materialize(switchWins, 1000).Average()
	t.Log(ExpectedValueWithConfidence(switchWins, SampleSize(20000)))
	if !Within(v, 0.666, epsilon) {
		t.Error("Switching should win 2/3 of the time")
	}
}
