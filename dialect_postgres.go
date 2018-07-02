/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"regexp"
)

type PostgresDialect struct {
	DefaultDialect
}

func NewPostgresDialect() *PostgresDialect {
	return &PostgresDialect{
		DefaultDialect{
			TypeMatch: []TypeMatch{
				{"integer", GoBaseType("int32")},
				{"bigint", GoBaseType("int64")},
				{"timestamp", TimeType},
				{"character varying", GoBaseType("string")},
				{"double precision", GoBaseType("float64")},
				{"real", GoBaseType("float32")},
				{"text", GoBaseType("string")},
				{"boolean", GoBaseType("bool")},
			},
		},
	}
}

func (dia PostgresDialect) ParserTables(db *sql.DB, filter *regexp.Regexp, ctx *Context) {
	tablesRow, err := db.Query("SELECT tablename FROM pg_tables where schemaname=$1", "public")
	if err != nil {
		panic(err)
	}
	for tablesRow.Next() {
		var tableName string
		err := tablesRow.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		model := dia.parserModel(db, tableName, ctx)
		ctx.Models = append(ctx.Models, model)
	}

}

func (dia PostgresDialect) parserModel(db *sql.DB, name string, ctx *Context) *Model {
	fieldsRow, err := db.Query("SELECT column_name,data_type,is_nullable,column_default FROM information_schema.columns "+
		"WHERE table_schema = 'public' AND table_name = $1", name)
	if err != nil {
		panic(err)
	}
	nextValRe := regexp.MustCompile("^nextval\\(.*\\)$")
	columnFieldMap := map[string]*Field{}
	model := NewModel(name)
	for fieldsRow.Next() {
		var columnName, dataType, isNull string
		var _default *string
		err := fieldsRow.Scan(&columnName, &dataType, &isNull, &_default)
		if err != nil {
			panic(err)
		}
		goType := dia.TypeConvert(dataType)
		if goType.From != "" {
			ctx.Import[goType.From] = ""
		}
		field := NewField(columnName, goType.Name)
		field.Tag.Add("toyorm", "type", dataType)
		if _default != nil {
			if nextValRe.MatchString(*_default) {
				field.Tag.Add("toyorm", "auto_increment", "")
			} else {
				field.Tag.Add("toyorm", "default", *_default)
			}
		}
		field.IsNull = isNull == "YES"

		model.Fields = append(model.Fields, field)
		columnFieldMap[columnName] = field
	}

	indexsRow, err := db.Query("SELECT i.relname as index_name,"+
		"a.attname as column_name,"+
		"ix.indisunique as isunique,"+
		"ix.indisprimary as isprimary "+
		"FROM pg_class t,pg_class i,pg_index ix,pg_namespace ns, pg_attribute a "+
		"WHERE t.oid = ix.indrelid AND "+
		"i.oid = ix.indexrelid AND "+
		"a.attrelid = t.oid AND "+
		"a.attnum = ANY(ix.indkey) AND "+
		"t.relkind = 'r' AND "+
		"t.relnamespace = ns.oid AND "+
		"ns.nspname = 'public' AND "+
		"t.relname = $1",
		name)
	if err != nil {
		panic(err)
	}
	for indexsRow.Next() {
		var idxName, columnName string
		var isUnique, isPrimary bool
		err := indexsRow.Scan(&idxName, &columnName, &isUnique, &isPrimary)
		if err != nil {
			panic(err)
		}
		if isPrimary {
			columnFieldMap[columnName].Tag.Add("toyorm", "primary key", "")
		} else if isUnique {
			columnFieldMap[columnName].Tag.Add("toyorm", "unique index", idxName)
		} else {
			columnFieldMap[columnName].Tag.Add("toyorm", "index", idxName)
		}
	}
	return model
}
