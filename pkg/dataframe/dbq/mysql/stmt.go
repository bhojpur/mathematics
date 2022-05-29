package sql

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
	stdSql "database/sql"
	"time"
)

// Stmt is a prepared statement.
// A Stmt is safe for concurrent use by multiple goroutines.
type Stmt struct {
	stmt         *stdSql.Stmt
	killerPool   StdSQLDB
	connectionID string
	kto          time.Duration
}

// Unleak will release the reference to the killerPool
// in order to prevent a memory leak.
func (s *Stmt) Unleak() {
	s.killerPool = nil
	s.connectionID = ""
	s.kto = 0
}

// Close closes the statement.
func (s *Stmt) Close() error {
	err := s.stmt.Close()
	if err != nil {
		return err
	}
	s.Unleak() // Should this be called in a defer to guarantee it gets called?
	return nil
}

// Exec executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (s *Stmt) Exec(args ...interface{}) (stdSql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

// ExecContext executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (stdSql.Result, error) {

	// Create a context that is used to cancel ExecContext()
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	outChan := make(chan stdSql.Result)
	errChan := make(chan error)
	returnedChan := make(chan struct{}) // Used to indicate that this function has returned

	defer func() {
		returnedChan <- struct{}{}
	}()

	go func() {
		select {
		case <-ctx.Done():
			// context has been canceled
			kill(s.killerPool, s.connectionID, s.kto)
			errChan <- ctx.Err()
		case <-returnedChan:
		}
	}()

	go func() {
		res, err := s.stmt.ExecContext(cancelCtx, args...)
		if err != nil {
			errChan <- err
			return
		}
		outChan <- res
	}()

	select {
	case err := <-errChan:
		return nil, err
	case out := <-outChan:
		return out, nil
	}
}

// Query executes a prepared query statement with the given arguments
// and returns the query results as a *Rows.
func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

// QueryContext executes a prepared query statement with the given arguments
// and returns the query results as a *Rows.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error) {

	// We can't use the same approach used in ExecContext because defer cancelFunc()
	// cancels rows.Scan.
	defer func() {
		if ctx.Err() != nil {
			kill(s.killerPool, s.connectionID, s.kto)
		}
	}()

	rows, err := s.stmt.QueryContext(ctx, args...)
	return &Rows{ctx: ctx, rows: rows, killerPool: s.killerPool, connectionID: s.connectionID}, err
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
//
// Example usage:
//
//  var name string
//  err := nameByUseridStmt.QueryRow(id).Scan(&name)
func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}

// QueryRowContext executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {

	// Since sql.Row does not export err field, this is the best we can do:
	defer func() {
		if ctx.Err() != nil {
			kill(s.killerPool, s.connectionID, s.kto)
		}
	}()

	row := s.stmt.QueryRowContext(ctx, args...)
	return &Row{ctx: ctx, row: row, killerPool: s.killerPool, connectionID: s.connectionID}
}
