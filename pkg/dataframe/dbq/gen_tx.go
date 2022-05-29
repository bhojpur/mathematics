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
	"fmt"
	"time"

	rlSql "github.com/bhojpur/mathematics/pkg/dataframe/dbq/mysql"
	"github.com/cenkalti/backoff/v4"
)

type txer interface {
	Commit() error
	Rollback() error
}

// QFn is shorthand for Q. It will automatically use the appropriate transaction.
type QFn func(ctx context.Context, query string, options *Options, args ...interface{}) (interface{}, error)

// EFn is shorthand for E. It will automatically use the appropriate transaction.
type EFn func(ctx context.Context, query string, options *Options, args ...interface{}) (sql.Result, error)

// TxCommit will commit the transaction.
type TxCommit func() error

// Tx is used to perform an arbitrarily complex operation and not have to worry about rolling back a transaction.
// The transaction is automatically rolled back unless explicitly committed by calling txCommit.
// tx is only exposed for performance purposes. Do not use it to commit or rollback.
//
// NOTE: Until this note is removed, this function is not necessarily backward compatible.
//
// Example:
//
//  ctx := context.Background()
//  pool, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db")
//
//  dbq.Tx(ctx, pool, func(tx interface{}, Q dbq.QFn, E dbq.EFn, txCommit dbq.TxCommit) {
//     stmt := dbq.INSERTStmt("table", []string{"name", "age", "created_at"}, 1)
//     res, err := E(ctx, stmt, nil, "test name", 34, time.Now())
//     if err != nil {
//        return // Automatic rollback
//     }
//     txCommit()
//  })
//
func Tx(ctx context.Context, db interface{}, fn func(tx interface{}, Q QFn, E EFn, txCommit TxCommit), retryPolicy ...backoff.BackOff) error {
	if ctx == nil {
		ctx = context.Background()
	}
	var (
		alreadyTx bool
		tx        interface{}
		err       error
	)

	switch db := db.(type) {
	case BeginTxer:
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
	case beginTxer2:
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
	case *sql.Tx, *rlSql.Tx:
		tx = db
		alreadyTx = true
	default:
		panic(fmt.Sprintf("interface conversion: %T is not dbq.BeginTxer: missing method: BeginTx", db))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.(txer).Rollback()
			panic(r)
		}
	}()

	qFn := func(ctx context.Context, query string, options *Options, args ...interface{}) (interface{}, error) {
		res, err := Q(ctx, tx, query, options, args...)
		if err == sql.ErrTxDone && !alreadyTx {
			return Q(ctx, db, query, options, args...)
		}
		return res, err
	}

	eFn := func(ctx context.Context, query string, options *Options, args ...interface{}) (sql.Result, error) {
		return E(ctx, tx.(ExecContexter), query, options, args...)
	}

	completed := false
	txCommit := func() error {
		err := tx.(txer).Commit()
		if err == nil || err == sql.ErrTxDone {
			completed = true
			return nil
		}
		return err
	}

	operation := func() error {
		fn(tx, qFn, eFn, txCommit)
		if completed {
			return nil
		}

		op2 := func() error {
			err = tx.(txer).Rollback()
			if err == sql.ErrTxDone {
				return nil
			}
			return err
		}

		exp := ExponentialRetryPolicy(120 * time.Second)
		err := backoff.Retry(op2, backoff.WithContext(exp, ctx))
		if err != nil {
			return &backoff.PermanentError{err}
		}
		return nil
	}

	if !(len(retryPolicy) > 0 && retryPolicy[0] != nil) {

		return operation()
	}

	return backoff.Retry(operation, backoff.WithContext(retryPolicy[0], ctx))
}
