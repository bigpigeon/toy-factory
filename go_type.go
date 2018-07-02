/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

type GoType struct {
	Name string
	From string
}

func (g *GoType) String() string {
	return g.Name
}

func GoBaseType(name string) *GoType {
	return &GoType{Name: name}
}

var TimeType = &GoType{"time.Time", "time"}
