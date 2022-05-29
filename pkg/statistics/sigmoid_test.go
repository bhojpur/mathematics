package statistics_test

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
	"testing"

	stats "github.com/bhojpur/mathematics/pkg/statistics"
)

func ExampleSigmoid() {
	s, _ := stats.Sigmoid([]float64{3.0, 1.0, 2.1})
	fmt.Println(s)
	// Output: [0.9525741268224334 0.7310585786300049 0.8909031788043871]
}

func TestSigmoidEmptyInput(t *testing.T) {
	_, err := stats.Sigmoid([]float64{})
	if err != stats.EmptyInputErr {
		t.Errorf("Should have returned empty input error")
	}
}

func TestSigmoid(t *testing.T) {
	sm, err := stats.Sigmoid([]float64{-0.54761371, 17.04850603, 4.86054302})
	if err != nil {
		t.Error(err)
	}

	a := 0.3664182235138545
	if sm[0] != a {
		t.Errorf("%v != %v", sm[0], a)
	}

	a = 0.9999999605608187
	if sm[1] != a {
		t.Errorf("%v != %v", sm[1], a)
	}

	a = 0.9923132671908277
	if sm[2] != a {
		t.Errorf("%v != %v", sm[2], a)
	}
}
