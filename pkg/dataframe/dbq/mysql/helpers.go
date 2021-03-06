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
	"time"
)

// kill is used to kill a running query.
// It is advised that db be another pool that the
// connection was NOT derived from.
func kill(db StdSQLDB, connectionID string, kto time.Duration) error {

	if connectionID == "" {
		return nil
	}

	stmt := `KILL QUERY ?`

	if kto == 0 {
		_, err := db.Exec(stmt, connectionID)
		if err != nil {
			return err
		}
	} else {
		ctx, cancelFunc := context.WithTimeout(context.Background(), kto)
		defer cancelFunc()
		_, err := db.ExecContext(ctx, stmt, connectionID)
		if err != nil {
			return err
		}
	}

	return nil
}
