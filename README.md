## SQL worker pool wrapper for [sqlx](https://github.com/jmoiron/sqlx) in Go.

[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/jimmy-go/pgwp.svg?branch=master)](https://travis-ci.org/jimmy-go/pgwp)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimmy-go/pgwp)](https://goreportcard.com/report/github.com/jimmy-go/pgwp)
[![GoDoc](http://godoc.org/github.com/jimmy-go/pgwp?status.png)](http://godoc.org/github.com/jimmy-go/pgwp)
[![Coverage Status](https://coveralls.io/repos/github/jimmy-go/pgwp/badge.svg?branch=master)](https://coveralls.io/github/jimmy-go/pgwp?branch=master)

### Install:

```
go get gopkg.in/jimmy-go/pgwp.v0
```

### Usage:

Equal to sqlx, just change Connect func by Open func.

```
// NOTE: Declare a new pool with 5 connections (*sqlx.DB).
db, err := pgwp.Open("postgres", "connection-string", 5)

var list []Item
db.SelectContext(ctx, &list, "SELECT * FROM people WHERE name=$1 LIMIT 10", "lisa" )
```

### Notes:

+ Only Context methods are supported. This ensures Go good practices.
+ Really useful working with SQLite because of [this](https://github.com/mattn/go-sqlite3/issues/50).
+ For not implemented sqlx methods use Execute method.

### License:

The MIT License (MIT)

Copyright (c) 2016 Angel del Castillo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
