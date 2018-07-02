/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"github.com/bigpigeon/toyorm"
	"sort"
	"strings"
)

type Tag map[string]map[string]string

func (t Tag) Add(scope string, name string, value string) {
	if t[scope] == nil {
		t[scope] = map[string]string{}
	}
	t[scope][name] = value
}

func (t Tag) Del(scope string, name string) {
	if t[scope] == nil {
		return
	}
	delete(t[scope], name)
}

func (t Tag) String() string {
	tagScopes := make([]string, 0, len(t))
	for key := range t {
		tagScopes = append(tagScopes, key)
	}
	sort.Strings(tagScopes)
	keyValueList := make([]string, 0, len(t))
	for _, tagKey := range tagScopes {
		keys := make([]string, 0, len(t[tagKey]))
		for k := range t[tagKey] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		keyValuePair := make([]string, 0, len(keys))

		for _, k := range keys {
			value := t[tagKey][k]
			if value != "" {
				keyValuePair = append(keyValuePair, fmt.Sprintf("%s:%s", k, value))
			} else {
				keyValuePair = append(keyValuePair, k)
			}
		}
		keyValueList = append(keyValueList, fmt.Sprintf("%s:\"%s\"", tagKey, strings.Join(keyValuePair, ";")))
	}
	return strings.Join(keyValueList, " ")
}

type Field struct {
	Name         string
	Type         string
	Tag          Tag
	IsPrimaryKey bool
	IsNull       bool
}

func NewField(name, _type string) *Field {
	field := &Field{
		Name: GoNameConvert(name),
		Type: _type,
		Tag:  Tag{},
	}
	// if toyorm SqlNameConvert not match , add column tag attribute
	if toyorm.SqlNameConvert(field.Name) != name {
		field.Tag.Add("toyorm", "column", name)
	}
	return field
}

func (f *Field) String() string {
	if f.IsPrimaryKey {
		f.Tag.Add("toyorm", "primary key", "")
	}
	if f.IsNull == false {
		f.Tag.Del("toyorm", "NULL")
		f.Tag.Add("toyorm", "NOT NULL", "")
		return fmt.Sprintf("%s %s `%s`", f.Name, f.Type, f.Tag)
	} else {
		f.Tag.Del("toyorm", "NOT NULL")
		f.Tag.Add("toyorm", "NULL", "")
		return fmt.Sprintf("%s *%s `%s`", f.Name, f.Type, f.Tag)
	}
}

type Method interface {
	String() string
}

type TableNameMethod struct {
	Model      *Model
	TablerName string
}

func (tabler *TableNameMethod) String() string {
	return fmt.Sprintf(`
func (t *%[1]s) TableName() string {
	return "%[2]s"
}
`, tabler.Model.Name, tabler.TablerName)
}

type Model struct {
	Name    string
	Fields  []*Field
	Methods []Method
}

func NewModel(name string) *Model {
	model := &Model{
		Name: GoNameConvert(name),
	}
	// if SqlNameConvert not match add Tabler method
	if toyorm.SqlNameConvert(model.Name) != name {
		model.Methods = append(model.Methods, &TableNameMethod{model, name})
	}
	return model
}

func (m *Model) String() string {
	var fieldStrList []string
	for _, field := range m.Fields {
		fieldStrList = append(fieldStrList, field.String())
	}
	structStr := fmt.Sprintf(`
// [toyorm.Model]
type %s struct {
%s
}
`, m.Name, strings.Join(fieldStrList, "\n"))
	for _, method := range m.Methods {
		structStr += method.String()
	}
	return structStr
}
