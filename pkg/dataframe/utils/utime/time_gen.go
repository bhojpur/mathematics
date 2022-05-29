package utime

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
	"time"
)

// NextTime will return the next time in the sequence. You can call it repeatedly
// to obtain a sequence.
type NextTime func() time.Time

// TimeGenerator will create a generator function, which when called will return
// a sequence of times. The sequence will begin at startTime. When reverse is true,
// the sequence will be backwards.
type TimeGenerator func(startTime time.Time, reverse bool) NextTime

// TimeIntervalGenerator is used to create a sequence of times based on an interval defined by
// timeFreq. timeFreq can be in the format: nYnMnWnD, where n is a non-negative integer and
// Y, M, W and D represent years, months, weeks and days respectively. Alternatively, timeFreq
// can be a valid positive input to time.ParseDuration.
//
// Example:
//
//  gen, _ := utime.TimeIntervalGenerator("1W1D")
//  ntg := gen(time.Now().UTC(), false)
//  for {
//     fmt.Println(ntg())
//     time.Sleep(500 * time.Millisecond)
//  }
//
// See: https://golang.org/pkg/time/#ParseDuration
func TimeIntervalGenerator(timeFreq string) (TimeGenerator, error) {

	// Prevent negative sign
	if len(timeFreq) > 0 && timeFreq[0:1] == "-" {
		return nil, fmt.Errorf("negative sign disallowed: %s", timeFreq)
	}

	var (
		d *time.Duration
		p *parsed
	)

	_d, err := time.ParseDuration(timeFreq)
	if err != nil {
		_p, err := parse(timeFreq)
		if err != nil {
			return nil, fmt.Errorf("could not parse: %s", timeFreq)
		}
		if _p.isZero() {
			return nil, fmt.Errorf("can't be zero: %s", timeFreq)
		}
		p = &_p
	} else {
		if _d == 0 {
			return nil, fmt.Errorf("can't be zero: %s", timeFreq)
		}
		d = &_d
	}

	return func(startTime time.Time, reverse bool) NextTime {
		var prevTime *time.Time

		return func() time.Time {
			var nt time.Time

			if prevTime == nil {
				nt = startTime
			} else {
				if d == nil {
					nt = (*prevTime).AddDate((*p).addDate(reverse))
				} else {
					if reverse {
						nt = (*prevTime).Add(-*d)
					} else {
						nt = (*prevTime).Add(*d)
					}
				}
			}
			prevTime = &nt
			return nt
		}
	}, nil
}
