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

type OptionType int

const (
	sampleSizeOpt OptionType = iota
	zScoreOpt
	percentErrorOpt
)

type Option struct {
	optionType OptionType
	intVal     int
	floatVal   float64
}

func SampleSize(n int) Option {
	return Option{
		optionType: sampleSizeOpt,
		intVal:     n,
	}
}

func getSampleSize(opts []Option, def int) int {
	for _, v := range opts {
		if v.optionType == sampleSizeOpt {
			return v.intVal
		}
	}
	return def
}

const zScore95 = 1.96

func ZScore95() Option {
	return ZScore(zScore95)
}

const zScore99 = 2.58

func ZScore99() Option {
	return ZScore(zScore99)
}

const zScore90 = 1.64

func ZScore90() Option {
	return ZScore(zScore90)
}

func ZScore(score float64) Option {
	return Option{
		optionType: zScoreOpt,
		floatVal:   score,
	}
}

func getZScore(opts []Option, def float64) float64 {
	for _, v := range opts {
		if v.optionType == zScoreOpt {
			return v.floatVal
		}
	}
	return def
}

func PercentError(v float64) Option {
	return Option{
		optionType: percentErrorOpt,
		floatVal:   v,
	}
}

func getPercentError(opts []Option, def float64) float64 {
	for _, v := range opts {
		if v.optionType == percentErrorOpt {
			return v.floatVal
		}
	}
	return def
}
