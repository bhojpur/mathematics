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
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"cloud.google.com/go/civil"
	"github.com/cenkalti/backoff/v4"
	"github.com/mitchellh/mapstructure"
)

// MustQ is a wrapper around the Q function. It will panic upon encountering an error.
// This can erradicate boiler-plate error handing code.
func MustQ(ctx context.Context, db interface{}, query string, options *Options, args ...interface{}) interface{} {
	aRzLNT, XYeUCW := Q(ctx, db, query, options, args...)
	if XYeUCW != nil {
		panic(XYeUCW)
	}
	return aRzLNT
}

// Q is used for querying a SQL database. A []map[string]interface{} is ordinarily returned.
// Each returned row (an item in the slice) contains a map where the keys are the columns, and
// the values are the data for each column.
// However, when a ConcreteStruct is provided via the options, the mapstructure package is used to automatically
// return []*struct instead. To bypass the mapstructure package, ScanFaster interface can be implemented.
//
// args is a list of values to replace the placeholders in the query. When an arg is a slice, the values of the slice
// will automatically be flattened to a list of interface{}.
//
// NOTE: sql.ErrNoRows is never returned as an error: A slice is always returned, unless the
// behavior is modified by the SingleResult Option.
func Q(ctx context.Context, db interface{}, query string, options *Options, args ...interface{}) (out interface{}, rErr error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var o Options
	if options != nil {
		o = *options

		if o.RetryPolicy != nil {
			o.RetryPolicy = backoff.WithContext(o.RetryPolicy, ctx)
		}
	}

	defer func() {
		if rErr == nil && o.SingleResult {
			rows := reflect.ValueOf(out)
			if rows.Len() == 0 {
				out = nil
			} else {
				row := rows.Index(0)
				out = row.Interface()
			}
		}
	}()

	for _, v := range args {
		if arg := reflect.ValueOf(v); arg.Kind() == reflect.Slice {
			args = FlattenArgs(args...)
			break
		}
	}

	var (
		outStruct     interface{}
		outMap        = []map[string]interface{}{}
		scanFast      bool
		postUnmarshal bool
	)

	if o.ConcreteStruct != nil {

		csTyp := reflect.New(reflect.TypeOf(o.ConcreteStruct)).Interface()
		_, scanFast = csTyp.(ScanFaster)
		_, postUnmarshal = csTyp.(PostUnmarshaler)

		typ := reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(o.ConcreteStruct)))
		outStruct = reflect.MakeSlice(typ, 0, 0)
	}

	var (
		rows      rows
		err       error
		operation func() error
	)

	if o.RetryPolicy == nil {
		switch db := db.(type) {
		case QueryContexter:
			rows, err = db.QueryContext(ctx, query, args...)
		case queryContexter2:
			rows, err = db.QueryContext(ctx, query, args...)
		default:
			panic(fmt.Sprintf("interface conversion: %T is not dbq.QueryContexter: missing method: QueryContext", db))
		}
	} else {
		switch db := db.(type) {
		case QueryContexter:
			operation = func() error {
				rows, err = db.QueryContext(ctx, query, args...)
				if err != nil {
					if err == sql.ErrTxDone || err == sql.ErrConnDone || (strings.Contains(err.Error(), "sql: expected") && strings.Contains(err.Error(), "arguments, got")) {
						return &backoff.PermanentError{err}
					}
					return err
				}
				return nil
			}
		case queryContexter2:
			operation = func() error {
				rows, err = db.QueryContext(ctx, query, args...)
				if err != nil {
					if err == sql.ErrTxDone || err == sql.ErrConnDone || (strings.Contains(err.Error(), "sql: expected") && strings.Contains(err.Error(), "arguments, got")) {
						return &backoff.PermanentError{err}
					}
					return err
				}
				return nil
			}
		default:
			panic(fmt.Sprintf("interface conversion: %T is not dbq.QueryContexter: missing method: QueryContext", db))
		}

		err = backoff.Retry(operation, o.RetryPolicy)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	totalColumns := len(cols)

	for rows.Next() {
		var rowData []interface{}

		if scanFast {
			res := reflect.New(reflect.TypeOf(o.ConcreteStruct)).Interface()
			if err := rows.Scan(res.(ScanFaster).ScanFast()...); err != nil {
				return nil, err
			}
			outStruct = reflect.Append(outStruct.(reflect.Value), reflect.ValueOf(res))
			continue
		} else {
			rowData = make([]interface{}, totalColumns)
			for i := range rowData {
				rowData[i] = &sql.RawBytes{}
			}
			if err := rows.Scan(rowData...); err != nil {
				return nil, err
			}
		}

		vals := map[string]interface{}{}
		if o.ConcreteStruct != nil {
			for colID, elem := range rowData {
				fieldName := cols[colID].Name()
				raw := elem.(*sql.RawBytes)
				if *raw == nil {
					vals[fieldName] = nil
				} else {
					vals[fieldName] = string(*raw)
				}
			}

			res := reflect.New(reflect.TypeOf(o.ConcreteStruct)).Interface()
			if o.DecoderConfig != nil {
				dc := &mapstructure.DecoderConfig{
					DecodeHook:       o.DecoderConfig.DecodeHook,
					ZeroFields:       true,
					TagName:          "dbq",
					WeaklyTypedInput: o.DecoderConfig.WeaklyTypedInput,
					Result:           res,
				}
				decoder, err := mapstructure.NewDecoder(dc)
				if err != nil {
					return nil, err
				}
				err = decoder.Decode(vals)
				if err != nil {
					return nil, err
				}
			} else {
				dc := &mapstructure.DecoderConfig{
					ZeroFields:       true,
					TagName:          "dbq",
					WeaklyTypedInput: true,
					Result:           res,
				}
				decoder, err := mapstructure.NewDecoder(dc)
				if err != nil {
					return nil, err
				}
				err = decoder.Decode(vals)
				if err != nil {
					return nil, err
				}
			}
			outStruct = reflect.Append(outStruct.(reflect.Value), reflect.ValueOf(res))
			continue
		}

		for colID, elem := range rowData {
			fieldName := cols[colID].Name()
			raw := elem.(*sql.RawBytes)

			if o.RawResults {
				cpy := make([]byte, len(*raw))
				copy(cpy, []byte(*raw))
				vals[fieldName] = cpy
				continue
			}

			colType := cols[colID].DatabaseTypeName()
			nullable, hasNullableInfo := cols[colID].Nullable()

			var val *string

			if *raw != nil {
				val = &[]string{string(*raw)}[0]
			}

			switch colType {
			case "NULL":
				vals[fieldName] = nil
			case "CHAR", "VARCHAR", "TEXT", "NVARCHAR", "MEDIUMTEXT", "LONGTEXT":
				if nullable || !hasNullableInfo {
					vals[fieldName] = val
				} else {
					if hasNullableInfo {

						vals[fieldName] = *val
					}
				}
			case "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC", "FLOAT4", "FLOAT8":
				if nullable || !hasNullableInfo {
					if val == nil {
						vals[fieldName] = (*float64)(nil)
					} else {
						f, _ := strconv.ParseFloat(*val, 64)
						vals[fieldName] = &f
					}
				} else {
					if hasNullableInfo {

						f, _ := strconv.ParseFloat(*val, 64)
						vals[fieldName] = f
					}
				}
			case "INT", "TINYINT", "INT2", "INT4", "INT8", "MEDIUMINT", "SMALLINT", "BIGINT":

				switch cols[colID].ScanType().Kind() {
				case reflect.Uint:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*uint)(nil)
						} else {
							vals[fieldName] = parseUintP(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseUint(*val)
						}
					}
				case reflect.Uint8:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*uint8)(nil)
						} else {
							vals[fieldName] = parseUint8P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseUint8(*val)
						}
					}
				case reflect.Uint16:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*uint16)(nil)
						} else {
							vals[fieldName] = parseUint16P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseUint16(*val)
						}
					}
				case reflect.Uint32:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*uint32)(nil)
						} else {
							vals[fieldName] = parseUint32P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseUint32(*val)
						}
					}
				case reflect.Uint64:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*uint64)(nil)
						} else {
							vals[fieldName] = parseUint64P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseUint64(*val)
						}
					}
				case reflect.Int:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*int)(nil)
						} else {
							vals[fieldName] = parseIntP(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseInt(*val)
						}
					}
				case reflect.Int8:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*int8)(nil)
						} else {
							vals[fieldName] = parseInt8P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseInt8(*val)
						}
					}
				case reflect.Int16:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*int16)(nil)
						} else {
							vals[fieldName] = parseInt16P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseInt16(*val)
						}
					}
				case reflect.Int32:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*int32)(nil)
						} else {
							vals[fieldName] = parseInt32P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseInt32(*val)
						}
					}
				case reflect.Int64:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*int64)(nil)
						} else {
							vals[fieldName] = parseInt64P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseInt64(*val)
						}
					}
				default:
					if nullable || !hasNullableInfo {
						if val == nil {
							vals[fieldName] = (*int64)(nil)
						} else {
							vals[fieldName] = parseInt64P(*val)
						}
					} else {
						if hasNullableInfo {

							vals[fieldName] = parseInt64(*val)
						}
					}
				}
			case "BOOL":
				if nullable || !hasNullableInfo {
					if val == nil {
						vals[fieldName] = (*bool)(nil)
					} else {
						if *val == "true" || *val == "TRUE" || *val == "1" {
							vals[fieldName] = &[]bool{true}[0]
						} else {
							vals[fieldName] = &[]bool{false}[0]
						}
					}
				} else {
					if hasNullableInfo {

						if *val == "true" || *val == "TRUE" || *val == "1" {
							vals[fieldName] = true
						} else {
							vals[fieldName] = false
						}
					}
				}
			case "DATETIME", "TIMESTAMP", "TIMESTAMPTZ":
				if nullable || !hasNullableInfo {
					if val == nil {
						vals[fieldName] = (*time.Time)(nil)
					} else {
						t, err := time.Parse("2006-01-02 15:04:05", *val)
						if err != nil {
							t, _ = time.Parse(time.RFC3339, *val)
						}
						vals[fieldName] = &t
					}
				} else {
					if hasNullableInfo {

						t, err := time.Parse("2006-01-02 15:04:05", *val)
						if err != nil {
							t, _ = time.Parse(time.RFC3339, *val)
						}
						vals[fieldName] = &t
					}
				}
			case "JSON", "JSONB":
				if val == nil {
					vals[fieldName] = nil
				} else {
					var jData interface{}
					json.Unmarshal(*raw, &jData)
					vals[fieldName] = jData
				}
			case "DATE":
				if nullable || !hasNullableInfo {
					if val == nil {
						vals[fieldName] = (*civil.Date)(nil)
					} else {
						d, err := civil.ParseDate(*val)
						if err != nil {
							t, _ := time.Parse(time.RFC3339, *val)
							d = civil.Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}
						}
						vals[fieldName] = &d
					}
				} else {
					if hasNullableInfo {

						d, err := civil.ParseDate(*val)
						if err != nil {
							t, _ := time.Parse(time.RFC3339, *val)
							d = civil.Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}
						}
						vals[fieldName] = d
					}
				}
			case "TIME":
				if nullable || !hasNullableInfo {
					if val == nil {
						vals[fieldName] = (*civil.Time)(nil)
					} else {
						t, _ := civil.ParseTime(*val)
						vals[fieldName] = &t
					}
				} else {
					if hasNullableInfo {

						t, _ := civil.ParseTime(*val)
						vals[fieldName] = t
					}
				}

			default:

				if nullable || !hasNullableInfo {
					vals[fieldName] = val
				} else {
					if hasNullableInfo {

						vals[fieldName] = *val
					}
				}
			}
		}
		outMap = append(outMap, vals)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if o.PostFetch != nil {
		err := o.PostFetch(ctx)
		if err != nil {
			return nil, err
		}
	}

	if o.ConcreteStruct != nil {
		rows := outStruct.(reflect.Value)
		count := rows.Len()
		if count > 0 {
			if postUnmarshal {
				if o.ConcurrentPostUnmarshal && runtime.GOMAXPROCS(0) > 1 {
					g, newCtx := errgroup.WithContext(ctx)

					for i := 0; i < count; i++ {
						i := i
						g.Go(func() error {
							if err := newCtx.Err(); err != nil {
								return err
							}

							row := rows.Index(i).Interface()
							err := row.(PostUnmarshaler).PostUnmarshal(newCtx, i, count)
							if err != nil {
								return xerrors.Errorf("dbq.PostUnmarshal @ row %d: %w", i, err)
							}
							return nil
						})
					}

					if err := g.Wait(); err != nil {
						return nil, err
					}
				} else {
					for i := 0; i < count; i++ {
						if err := ctx.Err(); err != nil {
							return nil, err
						}

						row := rows.Index(i).Interface()
						err := row.(PostUnmarshaler).PostUnmarshal(ctx, i, count)
						if err != nil {
							return nil, xerrors.Errorf("dbq.PostUnmarshal @ row %d: %w", i, err)
						}
					}
				}
			}
		}
		return outStruct.(reflect.Value).Interface(), nil
	}

	return outMap, nil
}
