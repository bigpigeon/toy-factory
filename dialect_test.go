/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"github.com/bigpigeon/toyorm"
	"regexp"
	"testing"
)

type TestFactoryTable struct {
	toyorm.ModelDefault
	Name string `toyorm:"unique index"`
	Data string `toyorm:"type:text"`
}

func TestDialectTest(t *testing.T) {
	dialect := NewDialect(TestDriver)
	brick := TestDB.Model(&TestFactoryTable{})
	_, err := brick.DropTableIfExist()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = brick.CreateTableIfNotExist()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ctx := NewContext()
	dialect.ParserTables(TestDB.DB(), regexp.MustCompile(brick.Model.Name), ctx)
	if ctx.error != nil {
		t.Error(ctx.error)
	}
	for _, model := range ctx.Models {
		t.Log(model)
	}
}
