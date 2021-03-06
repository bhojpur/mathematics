package dbq

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
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/mitchellh/mapstructure"
)

// StructorConfig is used to expose a subset of the configuration options
// provided by the mapstructure package.
//
// See: https://godoc.org/github.com/mitchellh/mapstructure#DecoderConfig
type StructorConfig struct {

	// DecodeHook, if set, will be called before any decoding and any
	// type conversion (if WeaklyTypedInput is on). This lets you modify
	// the values before they're set down onto the resulting struct.
	//
	// If an error is returned, the entire decode will fail with that
	// error.
	DecodeHook mapstructure.DecodeHookFunc

	// If WeaklyTypedInput is true, the decoder will make the following
	// "weak" conversions:
	//
	//   - bools to string (true = "1", false = "0")
	//   - numbers to string (base 10)
	//   - bools to int/uint (true = 1, false = 0)
	//   - strings to int/uint (base implied by prefix)
	//   - int to bool (true if value != 0)
	//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
	//     FALSE, false, False. Anything else is an error)
	//   - empty array = empty map and vice versa
	//   - negative numbers to overflowed uint values (base 10)
	//   - slice of maps to a merged map
	//   - single values are converted to slices if required. Each
	//     element is weakly decoded. For example: "4" can become []int{4}
	//     if the target type is an int slice.
	//
	WeaklyTypedInput bool
}

// SingleResult is a convenient option for the common case of expecting
// a single result from a query.
var SingleResult = &Options{SingleResult: true}

// Options is used to modify the default behavior.
type Options struct {

	// ConcreteStruct can be set to any concrete struct (not a pointer).
	// When set, the mapstructure package is used to convert the returned
	// results automatically from a map to a struct. The `dbq` struct tag
	// can be used to map column names to the struct's fields.
	//
	// See: https://godoc.org/github.com/mitchellh/mapstructure
	ConcreteStruct interface{}

	// DecoderConfig is used to configure the decoder used by the mapstructure
	// package. If it's not supplied, a default StructorConfig is assumed. This means
	// WeaklyTypedInput is set to true and no DecodeHook is provided. Alternatively, if you require
	// a configuration for common time-based conversions, StdTimeConversionConfig is available.
	//
	// See: https://godoc.org/github.com/mitchellh/mapstructure
	DecoderConfig *StructorConfig

	// SingleResult can be set to true if you know the query will return at most 1 result.
	// When true, a nil is returned if no result is found. Alternatively, it will return the
	// single result directly (instead of wrapped in a slice). This makes it easier to
	// type assert.
	SingleResult bool

	// PostFetch is called after all results are fetched but before PostUnmarshaler is called (if applicable).
	// It can be used to return a database connection back to the pool.
	PostFetch func(ctx context.Context) error

	// ConcurrentPostUnmarshal can be set to true if PostUnmarshal must be called concurrently.
	ConcurrentPostUnmarshal bool

	// RawResults can be set to true for results to be returned unprocessed ([]byte).
	// This option does nothing if ConcreteStruct is provided.
	RawResults bool

	// RetryPolicy can be set if you want to retry the query in the event of failure.
	//
	// Example:
	//
	//  dbq.ExponentialRetryPolicy(60 * time.Second, 3)
	//
	RetryPolicy backoff.BackOff
}

// Q is a convenience function that calls dbq.Q.
// It allows you to recycle common options.
func (o *Options) Q(ctx context.Context, db interface{}, query string, args ...interface{}) (out interface{}, rErr error) {
	return Q(ctx, db, query, o, args...)
}

// MustQ is a convenience function that calls dbq.MustQ.
// It allows you to recycle common options.
func (o *Options) MustQ(ctx context.Context, db interface{}, query string, args ...interface{}) interface{} {
	return MustQ(ctx, db, query, o, args...)
}

// Qs is a convenience function that calls dbq.Qs.
// It allows you to recycle common options.
func (o *Options) Qs(ctx context.Context, db interface{}, query string, ConcreteStruct interface{}, args ...interface{}) (out interface{}, rErr error) {
	return Qs(ctx, db, query, ConcreteStruct, o, args...)
}

// MustQs is a convenience function that calls dbq.MustQs.
// It allows you to recycle common options.
func (o *Options) MustQs(ctx context.Context, db interface{}, query string, ConcreteStruct interface{}, args ...interface{}) interface{} {
	return MustQs(ctx, db, query, ConcreteStruct, o, args...)
}
