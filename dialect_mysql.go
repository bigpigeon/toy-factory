/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"fmt"
	"regexp"
)

type MysqlDialect struct {
	DefaultDialect
}

func NewMysqlDialect() *MysqlDialect {
	return &MysqlDialect{
		DefaultDialect{
			TypeMatch: []TypeMatch{
				{"int\\(\\d+\\)", GoBaseType("int32")},
				{"'int\\(\\d+\\) unsigned", GoBaseType("uint32")},
				{"timestamp", TimeType},
				{"varchar\\(\\d+\\)", GoBaseType("string")},
				{"double", GoBaseType("float64")},
				{"float", GoBaseType("float32")},
				{"text", GoBaseType("string")},
				{"boolean", GoBaseType("bool")},
			},
		},
	}
}

func (dia MysqlDialect) ParserTables(db *sql.DB, filter *regexp.Regexp, ctx *Context) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tab string
		rows.Scan(&tab)
		if filter == nil || filter.MatchString(tab) {
			tables = append(tables, tab)
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	for _, name := range tables {
		model := dia.parserCreateTable(db, name, ctx)
		ctx.Models = append(ctx.Models, model)
	}
}

func (dia MysqlDialect) parserCreateTable(db *sql.DB, name string, ctx *Context) *Model {
	rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`", name))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	columnFieldMap := map[string]*Field{}
	model := NewModel(name)
	for rows.Next() {
		var name, _type, extra, null string
		var key, _default *string
		err := rows.Scan(&name, &_type, &null, &key, &_default, &extra)
		if err != nil {
			panic(err)
		}
		goType := dia.TypeConvert(_type)
		if goType.From != "" {
			ctx.Import[goType.From] = ""
		}
		field := NewField(name, goType.Name)
		field.Tag.Add("toyorm", "type", _type)
		if _default != nil {
			field.Tag.Add("toyorm", "default", *_default)
		}
		if extra != "" {
			field.Tag.Add("toyorm", extra, "")
		}

		if null == "YES" {
			field.IsNull = true
		} else {
			field.IsNull = false
		}
		model.Fields = append(model.Fields, field)
		columnFieldMap[name] = field
	}

	idxRows, err := db.Query("SELECT COLUMN_NAME,INDEX_NAME,NON_UNIQUE FROM information_schema.statistics "+
		"WHERE TABLE_SCHEMA = (SELECT database()) AND TABLE_NAME = ?", name)
	if err != nil {

		panic(err)
	}
	defer idxRows.Close()
	for idxRows.Next() {
		var column, indexName string
		var nonUnique bool
		err := idxRows.Scan(&column, &indexName, &nonUnique)
		if err != nil {
			panic(err)
		}
		if indexName == "PRIMARY" {
			columnFieldMap[column].IsPrimaryKey = true
		} else if nonUnique {
			columnFieldMap[column].Tag.Add("toyorm", "index", indexName)
		} else {
			columnFieldMap[column].Tag.Add("toyorm", "unique index", indexName)
		}
	}

	return model
}
