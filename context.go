/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

type Context struct {
	Import map[string]string // Import[package]alias
	Models []*Model
	Pkg    string
	error  error
	Debug  bool
	// marking used name
	Used map[string]bool
}

func NewContext() *Context {
	return &Context{
		Import: map[string]string{},
		Used:   map[string]bool{},
	}
}

func (c *Context) Error(err error) {
	if c.Debug {
		panic(err)
	}
	c.error = err
}
