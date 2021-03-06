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

	rlSql "github.com/bhojpur/mathematics/pkg/dataframe/dbq/mysql"
)

// ScanFaster eradicates the use of the reflect package when unmarshaling.
// The ConcreteStruct pointer must implement this interface to make use of this feature.
// If you don't need to perform fancy time conversions or interpret weakly typed data (via mapstructure pkg), then
// this is the recommended approach as it is more performant.
//
// Example:
//
//  type user struct {
//     ID       int    `dbq:"id"`
//     Name     string `dbq:"name"`
//  }
//
//  func (u *user) ScanFast() []interface{} {
//     return []interface{}{&u.ID, &u.Name}
//  }
//
type ScanFaster interface {

	// ScanFast is used to directly scan the results from the query to the ConcreteStruct pointer.
	// The number of columns returned from the query must match the length of the slice returned.
	//
	// See: https://golang.org/pkg/database/sql/#Rows.Scan
	//
	// WARNING: "SELECT * FROM ..." may return more columns in the future if the table structure changes.
	ScanFast() []interface{}
}

// PostUnmarshaler allows you to further modify all results after unmarshaling.
// The ConcreteStruct pointer must implement this interface to make use of this feature.
//
// Example:
//
//  type user struct {
//     ID       int    `dbq:"id"`
//     Name     string `dbq:"name"`
//     HashedID string `dbq:"-"` // Obfuscated ID
//  }
//
//  func (u *user) PostUnmarshal(ctx context.Context, row, total int) error {
//     u.HashedID = obfuscate(u.ID)
//     return nil
//  }
//
type PostUnmarshaler interface {

	// PostUnmarshal is called for each row after all results have been fetched.
	// You can use it to further modify the values of each ConcreteStruct.
	PostUnmarshal(ctx context.Context, row, total int) error
}

// ExecContexter is for modifying the database state.
type ExecContexter interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// QueryContexter is for querying the database.
type QueryContexter interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type queryContexter2 interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*rlSql.Rows, error)
}

// SQLBasic allows for querying and executing statements.
type SQLBasic interface {
	ExecContexter
	QueryContexter
}

// BeginTxer is an object than can begin a transaction.
type BeginTxer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type beginTxer2 interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*rlSql.Tx, error)
}

type rows interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dest ...interface{}) error
}
