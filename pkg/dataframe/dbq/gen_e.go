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
	"database/sql"
	"reflect"
	"strings"

	"github.com/cenkalti/backoff/v4"
)

// MustE is a wrapper around the E function. It will panic upon encountering an error.
// This can erradicate boiler-plate error handing code.
func MustE(ctx context.Context, db ExecContexter, query string, options *Options, args ...interface{}) sql.Result {
	whTHct, cuAxhx := E(ctx, db, query, options, args...)
	if cuAxhx != nil {
		panic(cuAxhx)
	}
	return whTHct
}

// E is used for "Exec" queries such as insert, update and delete.
//
// args is a list of values to replace the placeholders in the query. When an arg is a slice, the values of the slice
// will automatically be flattened to a list of interface{}.
func E(ctx context.Context, db ExecContexter, query string, options *Options, args ...interface{}) (sql.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, v := range args {
		if arg := reflect.ValueOf(v); arg.Kind() == reflect.Slice {
			args = FlattenArgs(args...)
			break
		}
	}

	if options == nil || options.RetryPolicy == nil {
		return db.ExecContext(ctx, query, args...)
	}

	o := *options
	o.RetryPolicy = backoff.WithContext(o.RetryPolicy, ctx)

	var res sql.Result

	operation := func() error {
		var err error
		res, err = db.ExecContext(ctx, query, args...)
		if err != nil {
			if err == sql.ErrTxDone || err == sql.ErrConnDone || (strings.Contains(err.Error(), "sql: expected") && strings.Contains(err.Error(), "arguments, got")) {
				return &backoff.PermanentError{err}
			}
			return err
		}
		return nil
	}

	err := backoff.Retry(operation, o.RetryPolicy)
	if err != nil {
		return nil, err
	}

	return res, nil
}
