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

type Gaussian struct {
	mean   float64
	stddev float64
	i      int
}

var _ Uncertain = &Gaussian{}

func NewGaussian(mean, stddev float64) *Gaussian {
	return &Gaussian{
		mean:   mean,
		stddev: stddev,
		i:      newID(),
	}
}

func NewNormal(mean, stddev float64) *Gaussian {
	return NewGaussian(mean, stddev)
}

func (g *Gaussian) sample() float64 {
	r := randNormalFloat64()
	return (r * g.stddev) + g.mean
}

func (g *Gaussian) sampleWithTrace() *sample {
	val := g.sample()
	s := newSample(val)
	s.addTrace(g.i, val)
	return s
}

func (g *Gaussian) id() int {
	return g.i
}
