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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/cenkalti/backoff/v4"
)

// Database is used to set the Database.
// Different databases have different syntax for placeholders etc.
type Database int

const (
	// MySQL database
	MySQL Database = 0
	// PostgreSQL database
	PostgreSQL Database = 1
)

// INSERTStmt will generate an INSERT statement. It can be used for bulk inserts.
//
// NOTE: You may have to escape the column names. For MySQL, use backticks. Databases also have a limit
// to the number of query placeholders you can have. This will limit the number of rows you can insert.
func INSERTStmt(tableName string, columns []string, rows int, dbtype ...Database) string {
	return fmt.Sprintf("INSERT INTO %s ( %s ) VALUES %s", tableName, strings.Join(columns, ","), Ph(len(columns), rows, 0, dbtype...))
}

// INSERT is the legacy equivalent of INSERTStmt.
//
// Deprecated: It will be removed in v3. Use INSERTStmt instead.
func INSERT(tableName string, columns []string, rows int, dbtype ...Database) string {
	return INSERTStmt(tableName, columns, rows, dbtype...)
}

// Ph generates the placeholders for SQL queries.
// For a bulk insert operation, nRows is the number of rows you intend
// to insert, and nCols is the number of fields per row.
// For the IN function, set nRows to 1.
// For PostgreSQL, you can use incr to increment the placeholder starting count.
//
// NOTE: The function panics if either nCols or nRows is 0.
//
// Example:
//
//  dbq.Ph(3, 1, 0)
//  // Output: ( ?,?,? )
//
//  dbq.Ph(3, 2, 0)
//  // Output: ( ?,?,? ),( ?,?,? )
//
//  dbq.Ph(3, 2, 6, dbq.PostgreSQL)
//  // Output: ($7,$8,$9),($10,$11,$12)
//
func Ph(nCols, nRows int, incr int, dbtype ...Database) string {

	var typ Database
	if len(dbtype) > 0 {
		typ = dbtype[0]
	}

	if nCols == 0 {
		panic(errors.New("nCols must not be 0"))
	}

	if nRows == 0 {
		panic(errors.New("nRows must not be 0"))
	}

	if typ == MySQL {
		inner := "( " + strings.TrimSuffix(strings.Repeat("?,", nCols), ",") + " ),"
		return strings.TrimSuffix(strings.Repeat(inner, nRows), ",")
	}

	var singleValuesStr string

	varCount := 1 + incr
	for i := 1; i <= nRows; i++ {
		singleValuesStr = singleValuesStr + "("
		for j := 1; j <= nCols; j++ {
			singleValuesStr = singleValuesStr + fmt.Sprintf("$%d,", varCount)
			varCount++
		}
		singleValuesStr = strings.TrimSuffix(singleValuesStr, ",") + "),"
	}

	return strings.TrimSuffix(singleValuesStr, ",")
}

// FlattenArgs will accept a list of values and flatten any slices encountered.
//
// Example:
//
//  args1 := []string{"A", "B", "C"}
//  args2 := []interface{}{2, "D"}
//  args3 := dbq.Struct(Row{"Brad Pitt", 45, time.Now()})
//
//  dbq.FlattenArgs(args1, args2, args3)
//  // Output: []interface{}{"A", "B", "C", 2, "D", "Brad Pitt", 45, time.Now()}
func FlattenArgs(args ...interface{}) []interface{} {
	out := make([]interface{}, 0, len(args))

	var sliceConv func(reflect.Value)
	sliceConv = func(arg reflect.Value) {
		if arg.Kind() == reflect.Slice {
			for i := 0; i < arg.Len(); i++ {
				sliceConv(reflect.ValueOf(arg.Index(i).Interface()))
			}
		} else if !arg.IsValid() {
			out = append(out, nil)
		} else {
			out = append(out, arg.Interface())
		}
	}

	for i := range args {
		arg := args[i]
		if rarg := reflect.ValueOf(arg); rarg.Kind() == reflect.Slice {
			sliceConv(rarg)
		} else {
			out = append(out, arg)
		}
	}

	return out
}

// ExponentialRetryPolicy is a retry policy with exponentially increasing intervals between
// each retry attempt. If maxElapsedTime is 0, it will retry forever unless restricted by retryAttempts.
//
// See: https://godoc.org/gopkg.in/cenkalti/backoff.v4#ExponentialBackOff
func ExponentialRetryPolicy(maxElapsedTime time.Duration, retryAttempts ...uint64) backoff.BackOff {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = maxElapsedTime

	if len(retryAttempts) > 0 && retryAttempts[0] != 0 {
		return backoff.WithMaxRetries(bo, retryAttempts[0])
	}

	return bo
}

// ConstantDelayRetryPolicy is a retry policy with constant intervals between
// each retry attempt. It will retry forever unless restricted by retryAttempts.
//
// See: https://godoc.org/gopkg.in/cenkalti/backoff.v4#ConstantBackOff
func ConstantDelayRetryPolicy(interval time.Duration, retryAttempts ...uint64) backoff.BackOff {
	bo := backoff.NewConstantBackOff(interval)

	if len(retryAttempts) > 0 && retryAttempts[0] != 0 {
		return backoff.WithMaxRetries(bo, retryAttempts[0])
	}

	return bo
}

// StdTimeConversionConfig provides a standard configuration for unmarshaling to
// time-related fields in a struct. It properly converts timestamps and datetime columns into
// time.Time objects. It assumes a MySQL database as default.
func StdTimeConversionConfig(dbtype ...Database) *StructorConfig {

	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	if len(dbtype) > 0 && dbtype[0] == PostgreSQL {

		layouts[0], layouts[1] = layouts[1], layouts[0]
	}

	return &StructorConfig{
		WeaklyTypedInput: true,
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}

			switch t {
			case reflect.TypeOf(civil.Date{}):
				return civil.ParseDate(data.(string))
			case reflect.TypeOf(civil.DateTime{}):
				t, err := time.Parse(layouts[0], data.(string))
				if err != nil {
					t, err = time.Parse(layouts[1], data.(string))
					if err != nil {
						return nil, err
					}
				}
				return civil.DateTime{
					Date: civil.DateOf(t),
					Time: civil.TimeOf(t),
				}, nil
			case reflect.TypeOf(civil.Time{}):
				return civil.ParseTime(data.(string))
			case reflect.TypeOf(time.Time{}):
				t, err := time.Parse(layouts[0], data.(string))
				if err != nil {
					t, err := time.Parse(layouts[1], data.(string))
					if err != nil {
						return nil, err
					}
					return t, nil
				}
				return t, nil
			default:
				return data, nil
			}

			return data, nil
		},
	}
}

// Struct converts the fields of the struct into a slice of values.
// You can use it to convert a struct into the placeholder arguments required by
// the Q and E function. tagName is used to indicate the struct tag (default is "dbq").
// The function panics if strct is not an actual struct.
func Struct(strct interface{}, tagName ...string) []interface{} {

	tg := "dbq"

	if len(tagName) > 0 {
		tg = tagName[0]
	}

	out := []interface{}{}

	if strct == nil {
		panic(errors.New("strct must be a struct"))
	}

	s := reflect.ValueOf(strct)

	if s.Kind() == reflect.Ptr {
		s = reflect.Indirect(s)
	}
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := typeOfT.Field(i)

		if f.PkgPath != "" {

			continue
		}

		fieldTag := f.Tag.Get(tg)
		fieldValRaw := s.Field(i)
		fieldVal := fieldValRaw.Interface()

		if fieldValRaw.Kind() == reflect.Map {
			continue
		}

		if fieldTag == "-" || (strings.HasSuffix(fieldTag, ",omitempty") && reflect.DeepEqual(fieldVal, reflect.Zero(reflect.TypeOf(fieldVal)).Interface())) {
			continue
		}

		if fieldValRaw.Kind() == reflect.Slice {
			out = append(out, FlattenArgs(fieldVal)...)
			continue
		}

		out = append(out, fieldVal)
	}

	return out
}

// Qs operates the same as Q except it requires you to provide a ConcreteStruct as an argument.
// This allows you to recycle common options and conveniently provide a different ConcreteStruct.
func Qs(ctx context.Context, db interface{}, query string, ConcreteStruct interface{}, options *Options, args ...interface{}) (out interface{}, rErr error) {
	if ConcreteStruct == nil {
		panic("ConcreteStruct required")
	}
	var o Options
	if options == nil {
		o.ConcreteStruct = ConcreteStruct
	} else {
		o = *options
		o.ConcreteStruct = ConcreteStruct
	}
	return Q(ctx, db, query, &o, args...)
}

// MustQs is a wrapper around the Qs function. It will panic upon encountering an error.
// This can erradicate boiler-plate error handing code.
func MustQs(ctx context.Context, db interface{}, query string, ConcreteStruct interface{}, options *Options, args ...interface{}) interface{} {
	RjxAwn, wekrBE := Qs(ctx, db, query, ConcreteStruct, options, args...)
	if wekrBE != nil {
		panic(wekrBE)
	}
	return RjxAwn
}

func parseUintP(s string) *uint {
	n, _ := strconv.ParseUint(s, 10, 0)
	return &[]uint{uint(n)}[0]
}

func parseUint8P(s string) *uint8 {
	n, _ := strconv.ParseUint(s, 10, 8)
	return &[]uint8{uint8(n)}[0]
}

func parseUint16P(s string) *uint16 {
	n, _ := strconv.ParseUint(s, 10, 16)
	return &[]uint16{uint16(n)}[0]
}

func parseUint32P(s string) *uint32 {
	n, _ := strconv.ParseUint(s, 10, 32)
	return &[]uint32{uint32(n)}[0]
}

func parseUint64P(s string) *uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return &[]uint64{uint64(n)}[0]
}

func parseIntP(s string) *int {
	n, _ := strconv.ParseInt(s, 10, 0)
	return &[]int{int(n)}[0]
}

func parseInt8P(s string) *int8 {
	n, _ := strconv.ParseInt(s, 10, 8)
	return &[]int8{int8(n)}[0]
}

func parseInt16P(s string) *int16 {
	n, _ := strconv.ParseInt(s, 10, 16)
	return &[]int16{int16(n)}[0]
}

func parseInt32P(s string) *int32 {
	n, _ := strconv.ParseInt(s, 10, 32)
	return &[]int32{int32(n)}[0]
}

func parseInt64P(s string) *int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return &[]int64{int64(n)}[0]
}

func parseUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 0)
	return uint(n)
}

func parseUint8(s string) uint8 {
	n, _ := strconv.ParseUint(s, 10, 8)
	return uint8(n)
}

func parseUint16(s string) uint16 {
	n, _ := strconv.ParseUint(s, 10, 16)
	return uint16(n)
}

func parseUint32(s string) uint32 {
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint32(n)
}

func parseUint64(s string) uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return n
}

func parseInt(s string) int {
	n, _ := strconv.ParseInt(s, 10, 0)
	return int(n)
}

func parseInt8(s string) int8 {
	n, _ := strconv.ParseInt(s, 10, 8)
	return int8(n)
}

func parseInt16(s string) int16 {
	n, _ := strconv.ParseInt(s, 10, 16)
	return int16(n)
}

func parseInt32(s string) int32 {
	n, _ := strconv.ParseInt(s, 10, 32)
	return int32(n)
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
