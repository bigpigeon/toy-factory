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

type Exec struct {
	Exec string
	Args []interface{}
}

type TypeMatch struct {
	Match   string
	Convert *GoType
}

type Dialect interface {
	ParserTables(db *sql.DB, filter *regexp.Regexp, ctx *Context)
}

type DefaultDialect struct {
	TypeMatch []TypeMatch
}

func (dia DefaultDialect) TypeConvert(s string) *GoType {
	for _, match := range dia.TypeMatch {
		if ok := regexp.MustCompile(match.Match).MatchString(strings.ToLower(s)); ok {
			return match.Convert
		}
	}
	panic("not match name type " + s)
}

func NewDialect(name string) Dialect {
	switch name {
	case "mysql":
		return NewMysqlDialect()
	case "sqlite3", "sqlite":
		return NewSqlite3Dialect()
	case "postgres":
		return NewPostgresDialect()
	default:
		panic("not match dialect from" + name)
	}

}
