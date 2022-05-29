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

func ExampleSoftMax() {
	sm, _ := stats.SoftMax([]float64{3.0, 1.0, 0.2})
	fmt.Println(sm)
	// Output: [0.8360188027814407 0.11314284146556013 0.05083835575299916]
}

func TestSoftMaxEmptyInput(t *testing.T) {
	_, err := stats.SoftMax([]float64{})
	if err != stats.EmptyInputErr {
		t.Errorf("Should have returned empty input error")
	}
}

func TestSoftMax(t *testing.T) {
	sm, err := stats.SoftMax([]float64{3.0, 1.0, 0.2})
	if err != nil {
		t.Error(err)
	}

	a := 0.8360188027814407
	if sm[0] != a {
		t.Errorf("%v != %v", sm[0], a)
	}

	a = 0.11314284146556013
	if sm[1] != a {
		t.Errorf("%v != %v", sm[1], a)
	}

	a = 0.05083835575299916
	if sm[2] != a {
		t.Errorf("%v != %v", sm[1], a)
	}
}
