/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitWithExternalComma(t *testing.T) {
	data := SplitWithExternalComma("id int, data varchar(255), PRIMARY KEY(name, data)")
	assert.Equal(t, data, []string{"id int", " data varchar(255)", " PRIMARY KEY(name, data)"})
	data = SplitWithExternalComma("id int, data varchar(255), PRIMARY KEY(((name), (data)))")
	assert.Equal(t, data, []string{"id int", " data varchar(255)", " PRIMARY KEY(((name), (data)))"})

	data = SplitWithExternalComma("id int, `data,num` varchar(255), PRIMARY KEY(((name), (data)))")
	assert.Equal(t, data, []string{"id int", " `data,num` varchar(255)", " PRIMARY KEY(((name), (data)))"})

	data = SplitWithExternalComma("id int, 'data,num' varchar(255), PRIMARY KEY(((name), (data)))")
	assert.Equal(t, data, []string{"id int", " 'data,num' varchar(255)", " PRIMARY KEY(((name), (data)))"})

	data = SplitWithExternalComma("id int, \"data,num\" varchar(255), PRIMARY KEY(((name), (data)))")
	assert.Equal(t, data, []string{"id int", " \"data,num\" varchar(255)", " PRIMARY KEY(((name), (data)))"})

	data = SplitWithExternalComma("id int, 'data\\',num' varchar(255), PRIMARY KEY(((name), (data)))")
	assert.Equal(t, data, []string{"id int", " 'data\\',num' varchar(255)", " PRIMARY KEY(((name), (data)))"})
}
