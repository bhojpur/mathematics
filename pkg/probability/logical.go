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

type notOperation struct {
	b UncertainBool
	i int
}

func Not(booldist UncertainBool) UncertainBool {
	return &notOperation{booldist, newID()}
}

func (not *notOperation) Pr() bool {
	return Pr(not)
}

func (not *notOperation) sampleBool() bool {
	return !not.b.sampleBool()
}

func (not *notOperation) sample() float64 {
	return convertBoolSampleToFloat(not.sampleBool())
}

func (not *notOperation) sampleWithTrace() *sample {
	s := not.b.sampleWithTrace()
	s.value = 1.0 - s.value
	s.trace[not.i] = s.value
	return s
}

func (not *notOperation) id() int {
	return not.i
}

type logicOperation struct {
	a, b UncertainBool
	i    int
	op   func(a, b bool) bool
}

func Or(a, b UncertainBool) UncertainBool {
	return newLogicOperation(a, b, func(a, b bool) bool {
		return a || b
	})
}

func And(a, b UncertainBool) UncertainBool {
	return newLogicOperation(a, b, func(a, b bool) bool {
		return a && b
	})
}

func newLogicOperation(a, b UncertainBool, f func(a, b bool) bool) *logicOperation {
	return &logicOperation{
		a:  a,
		b:  b,
		i:  newID(),
		op: f,
	}
}

func (l *logicOperation) sampleBool() bool {
	return convertFloatSampleToBool(l.sample())
}

func (l *logicOperation) sample() float64 {
	return l.sampleWithTrace().value
}

func (l *logicOperation) sampleWithTrace() *sample {
	atrace := l.a.sampleWithTrace()
	if v, ok := atrace.trace[l.b.id()]; ok {
		s := l.op(
			convertFloatSampleToBool(atrace.value),
			convertFloatSampleToBool(v),
		)
		atrace.value = convertBoolSampleToFloat(s)
		atrace.addTrace(l.i, atrace.value)
		return atrace
	}
	btrace := l.b.sampleWithTrace()
	if _, ok := btrace.trace[l.a.id()]; ok {
		// We're dependent in the other direction. Let's optimize
		l.a, l.b = l.b, l.a
		return l.sampleWithTrace()
	}
	combined := atrace.combine(btrace)
	s := l.op(
		convertFloatSampleToBool(atrace.value),
		convertFloatSampleToBool(btrace.value),
	)
	combined.value = convertBoolSampleToFloat(s)
	combined.addTrace(l.i, combined.value)
	return combined
}

func (l *logicOperation) Pr() bool {
	return Pr(l)
}

func (l *logicOperation) id() int {
	return l.i
}
