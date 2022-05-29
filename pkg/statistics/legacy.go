package statistics

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

// VarP is a shortcut to PopulationVariance
func VarP(input Float64Data) (sdev float64, err error) {
	return PopulationVariance(input)
}

// VarS is a shortcut to SampleVariance
func VarS(input Float64Data) (sdev float64, err error) {
	return SampleVariance(input)
}

// StdDevP is a shortcut to StandardDeviationPopulation
func StdDevP(input Float64Data) (sdev float64, err error) {
	return StandardDeviationPopulation(input)
}

// StdDevS is a shortcut to StandardDeviationSample
func StdDevS(input Float64Data) (sdev float64, err error) {
	return StandardDeviationSample(input)
}

// LinReg is a shortcut to LinearRegression
func LinReg(s []Coordinate) (regressions []Coordinate, err error) {
	return LinearRegression(s)
}

// ExpReg is a shortcut to ExponentialRegression
func ExpReg(s []Coordinate) (regressions []Coordinate, err error) {
	return ExponentialRegression(s)
}

// LogReg is a shortcut to LogarithmicRegression
func LogReg(s []Coordinate) (regressions []Coordinate, err error) {
	return LogarithmicRegression(s)
}

// Legacy error names that didn't start with Err
var (
	EmptyInputErr = ErrEmptyInput
	NaNErr        = ErrNaN
	NegativeErr   = ErrNegative
	ZeroErr       = ErrZero
	BoundsErr     = ErrBounds
	SizeErr       = ErrSize
	InfValue      = ErrInfValue
	YCoordErr     = ErrYCoord
	EmptyInput    = ErrEmptyInput
)
