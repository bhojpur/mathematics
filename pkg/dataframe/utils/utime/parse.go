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
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile(`^(\d+Y)?(\d+M)?(\d+W)?(\d+D)?$`)

type parsed struct {
	years  int
	months int
	weeks  int
	days   int
}

func (p parsed) isZero() bool {
	if p.years == 0 && p.months == 0 && p.weeks == 0 && p.days == 0 {
		return true
	}
	return false
}

func (p parsed) String() string {

	if p.isZero() {
		return ""
	}

	var out string

	// Convert days to weeks
	rem := p.days % 7
	p.weeks = p.weeks + p.days/7
	p.days = rem

	if p.years != 0 {
		out = fmt.Sprintf("%dY", p.years)
	}

	if p.months != 0 {
		out = out + fmt.Sprintf("%dM", p.months)
	}

	if p.weeks != 0 {
		out = out + fmt.Sprintf("%dW", p.weeks)
	}

	if p.days != 0 {
		out = out + fmt.Sprintf("%dD", p.days)
	}

	return out
}

func (p parsed) addDate(reverse bool) (int, int, int) {
	if reverse {
		return -p.years, -p.months, -7*p.weeks - p.days
	}
	return p.years, p.months, 7*p.weeks + p.days
}

func parse(s string) (parsed, error) {
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return parsed{}, errors.New("could not parse")
	}
	return parsed{
		years:  parseComponent(matches[1]),
		months: parseComponent(matches[2]),
		weeks:  parseComponent(matches[3]),
		days:   parseComponent(matches[4]),
	}, nil
}

func parseComponent(s string) int {
	if s == "" {
		return 0
	}
	s = s[0 : len(s)-1] // Remove last letter
	n, _ := strconv.Atoi(s)
	return n
}
