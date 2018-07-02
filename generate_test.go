/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import "testing"

func TestGenerateRun(t *testing.T) {
	generator := NewGenerator("testdata/"+TestDriver, TestDB.DB(), NewDialect(TestDriver))
	generator.ctx.Debug = true
	err := generator.Run()
	if err != nil {
		t.Error(err)
	}
}
