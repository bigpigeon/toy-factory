/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"regexp"
	"strings"
)

type SqliteDialect struct {
	DefaultDialect
}

func NewSqlite3Dialect() *SqliteDialect {
	return &SqliteDialect{
		DefaultDialect{
			TypeMatch: []TypeMatch{
				{"integer", GoBaseType("int32")},
				{"real", GoBaseType("float32")},
				{"datetime", TimeType},
				{"timestamp", TimeType},
				{"varchar\\(\\d+\\)", GoBaseType("string")},
				{"text", GoBaseType("string")},
				{"bool", GoBaseType("bool")},
			},
		},
	}
}

func (dia SqliteDialect) ParserTables(db *sql.DB, filter *regexp.Regexp, ctx *Context) {
	rows, err := db.Query("SELECT name,sql FROM sqlite_master WHERE type='table'")
	if err != nil {
		panic(err)
	}
	createTableRe := regexp.MustCompile("^CREATE TABLE [^(]*\\((.*)\\)$")
	primaryKeysRe := regexp.MustCompile("^PRIMARY KEY[ ]*\\((.*)\\)$")
	fieldNameTrimFunc := func(r rune) bool {
		switch r {
		case ' ', '\'', '"', '`':
			return true
		}
		return false
	}
	for rows.Next() {
		var tableName, sqlStr string
		err := rows.Scan(&tableName, &sqlStr)
		// bult-in table ignore
		if tableName == "sqlite_sequence" {
			continue
		}
		if err != nil {
			panic(err)
		}
		model := NewModel(tableName)
		nameFieldMap := map[string]*Field{}
		fieldData := createTableRe.FindStringSubmatch(sqlStr)[1]
		for _, attr := range SplitWithExternalComma(fieldData) {
			attr = strings.TrimSpace(attr)
			if strings.HasPrefix(attr, "PRIMARY KEY") {
				fieldNames := primaryKeysRe.FindStringSubmatch(attr)[1]
				for _, fieldName := range SplitWithExternalComma(fieldNames) {
					fieldName = strings.TrimFunc(fieldName, fieldNameTrimFunc)
					nameFieldMap[fieldName].IsPrimaryKey = true
				}
			} else {
				fieldAttrs := strings.Split(attr, " ")
				fieldName := strings.TrimFunc(fieldAttrs[0], fieldNameTrimFunc)
				var fieldType string
				if len(fieldAttrs) > 1 {
					fieldType = fieldAttrs[1]
				} else {
					fieldType = "text"
				}
				goType := dia.TypeConvert(fieldType)
				if goType.From != "" {
					ctx.Import[goType.From] = ""
				}
				field := NewField(fieldName, goType.Name)

				field.Tag.Add("toyorm", "type", fieldType)
				if len(fieldAttrs) > 2 {
					dia.ParserFieldAttr(fieldAttrs[2:], field)
				}
				model.Fields = append(model.Fields, field)
				nameFieldMap[fieldName] = field
			}
		}
		ctx.Models = append(ctx.Models, model)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func (dia SqliteDialect) ParserFieldAttr(attrs []string, field *Field) {
	for i := 0; i < len(attrs); i++ {
		switch attrs[i] {
		case "DEFAULT":
			i++
			field.Tag.Add("toyorm", "default", attrs[i])
		case "NOT":
			i++
			field.Tag.Add("toyorm", "NOT "+attrs[i], "")
		default:
			field.Tag.Add("toyorm", attrs[i], "")
		}
	}
}
