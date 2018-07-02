/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"github.com/bigpigeon/toyorm"
	"testing"
)

var (
	TestDB     *toyorm.Toy
	TestDriver string
)

func TestMain(m *testing.M) {
	for _, sqldata := range []struct {
		Driver string
		Source string
	}{
		{"mysql", "root:@tcp(localhost:3306)/toyorm?charset=utf8&parseTime=True"},
		{"sqlite3", ":memory:"},
		{"postgres", "user=postgres dbname=toyorm sslmode=disable"},
	} {
		var err error
		TestDriver = sqldata.Driver
		TestDB, err = toyorm.Open(sqldata.Driver, sqldata.Source)
		if err != nil {
			fmt.Printf("Error: cannot open %s\n", err)
			goto Close
		}

		fmt.Printf("=========== %s ===========\n", sqldata.Driver)
		fmt.Printf("connect to %s \n\n", sqldata.Source)
		m.Run()
	Close:
		TestDB.Close()
	}
}
