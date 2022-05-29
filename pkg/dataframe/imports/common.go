package imports

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

// It provides functionality to read data contained in another format to populate a DataFrame.
// It provides inverse functionality to the exports.

// GenericDataConverter is used to convert input data into a generic data type.
// This is required when importing data for a Generic Series ("SeriesGeneric").
type GenericDataConverter func(in interface{}) (interface{}, error)

// Converter is used to convert input data into a generic data type.
// This is required when importing data for a Generic Series ("dataframe.SeriesGeneric").
// As a special case, if ConcreteType is time.Time, then a SeriesTime is used.
//
// Example:
//
//  opts := imports.CSVLoadOptions{
//     DictateDataType: map[string]interface{}{
//        "Date": imports.Converter{
//           ConcreteType: time.Time{},
//           ConverterFunc: func(in interface{}) (interface{}, error) {
//              return time.Parse("2006-01-02", in.(string))
//           },
//        },
//     },
//  }
//
type Converter struct {
	ConcreteType  interface{}
	ConverterFunc GenericDataConverter
}

// parseObject converts maps within maps and moves them to the root level with
// dots.
// eg. {"A":123, "B":{"C": "D"}} => {"A":123, "B.C":"D"}
func parseObject(v map[string]interface{}, prefix string) map[string]interface{} {
	out := map[string]interface{}{}

	for k, t := range v {
		var key string
		if prefix == "" {
			key = k
		} else {
			key = prefix + "." + k
		}

		switch v := t.(type) {
		case map[string]interface{}:
			for k, t := range parseObject(v, key) {
				out[k] = t
			}
		default:
			out[key] = t
		}
	}

	return out
}
