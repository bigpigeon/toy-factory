/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	filter     = flag.String("filter", "", "filter table name by regular expression")
	dbname     = flag.String("name", "", "database type name only support mysql/sqlite3/postgres")
	dburl      = flag.String("url", "", "database connect url")
	packageDir = flag.String("dir", ".", "default go directory")
)

func Main() {
	db, err := sql.Open(*dbname, *dburl)
	if err != nil {
		fmt.Println(err)
		return
	}
	generator := NewGenerator(*packageDir, db, NewDialect(*dbname))
	err = generator.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	flag.Parse()
	Main()
}
