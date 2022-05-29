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

import "fmt"

type Samples struct {
	Samples []*sample
	n       int
	i       int
}

func FromSamples(samples []*sample) *Samples {
	return &Samples{
		Samples: samples,
		n:       0,
		i:       newID(),
	}
}

func (s *Samples) sampleWithTrace() *sample {
	if len(s.Samples) == 0 {
		panic("Must have at least some samples in a sampling distribution")
	}
	if s.n >= len(s.Samples) {
		s.n = 0
	}
	out := s.Samples[s.n]
	s.n += 1
	out.addTrace(s.i, out.value)
	return out
}

func (s *Samples) sample() float64 {
	return s.sampleWithTrace().value
}

func (s *Samples) addSample(sample *sample) {
	s.Samples = append(s.Samples, sample)
}

func (s *Samples) Average() float64 {
	total := 0.0
	for _, v := range s.Samples {
		total += v.value
	}
	n := float64(len(s.Samples))
	return total / n
}

func (s *Samples) First() float64 {
	if len(s.Samples) == 0 {
		panic("No samples in the Sampling distribution")
	}
	return s.Samples[0].value
}

func (s *Samples) id() int {
	return s.i
}

func (s *Samples) String() string {
	var out string
	for i, v := range s.Samples {
		if i == 0 {
			out = v.String()
			continue
		}
		out = fmt.Sprintf("%s\n%s", out, v.String())
	}
	return out
}

func Materialize(u Uncertain, n int) *Samples {
	out := &Samples{
		i: newID(),
	}
	for i := 0; i < n; i++ {
		out.addSample(u.sampleWithTrace())
	}
	return out
}
