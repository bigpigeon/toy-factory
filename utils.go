/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

func GoNameConvert(s string) string {
	var nextUpperCase bool
	sByte := []byte(s)
	goName := make([]byte, 0, len(sByte))
	if len(sByte) > 0 {
		b := sByte[0]
		if b >= 'a' && b <= 'z' {
			goName = append(goName, b-'a'+'A')
			sByte = sByte[1:]
		}
	}
	for _, b := range sByte {
		if b >= 'a' && b <= 'z' {
			if nextUpperCase {
				goName = append(goName, b-'a'+'A')
			} else {
				goName = append(goName, b)
			}
			nextUpperCase = false
		} else if b >= 'A' && b <= 'Z' {
			goName = append(goName, b)
			nextUpperCase = false
		} else if b >= '0' && b <= '9' {
			goName = append(goName, b)
			nextUpperCase = true
		} else {
			nextUpperCase = true
		}
	}
	return string(goName)
}

func SplitWithExternalComma(s string) []string {
	type Status int
	const (
		External = Status(iota)
		InDoubleQuotes
		InQuotes
		InBackQuotes
		InBrackets
	)

	var status Status
	var bracketNum int
	var data []string
	pre, i := 0, 0
	for ; i < len(s); i++ {
		switch status {
		case External:
			switch s[i] {
			case '"':
				status = InDoubleQuotes
			case '\'':
				status = InQuotes
			case '`':
				status = InBackQuotes
			case '(':
				bracketNum++
				status = InBrackets
			case ',':
				data = append(data, s[pre:i])
				pre = i + 1
			}
		case InDoubleQuotes:
			switch s[i] {
			case '"':
				status = External
			case '\\':
				i++
			}
		case InQuotes:
			switch s[i] {
			case '\'':
				status = External
			case '\\':
				i++
			}
		case InBackQuotes:
			switch s[i] {
			case '`':
				status = External
			case '\\':
				i++
			}
		case InBrackets:
			if s[i] == '(' {
				bracketNum++
			} else if s[i] == ')' {
				bracketNum--
				if bracketNum == 0 {
					status = External
				}
			}
		}
	}
	if pre < len(s) {
		data = append(data, s[pre:i])
	}

	return data

}
