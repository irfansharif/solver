// Copyright 2021 Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package token

// Type represents the type of a given token.
type Type int

// Token consists of a type and value; it's the unit output of a lexer.
type Token struct {
	Type  Type
	Value string
}

//go:generate stringer -type=Type
const (
	ILLEGAL Type = iota + 128
	EOF

	// Words and digits.
	WORD   // x, yyz, ...
	DIGITS // 42, 1343, ...

	// Operations.
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /
	IMPL     // →
	MOD      // %
	LT       // <
	GT       // >
	EXISTS   // ∈
	NEXISTS  // ∉
	UNION    // ∪
	EQ       // ==
	NEQ      // !=

	// Delimiters.
	DOT      // .
	COLON    // :
	COMMA    // ,
	PIPE     // |
	SUM      // Σ
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]

	// Keywords.
	AS   // "as"
	IF   // "if"
	IN   // "in"
	MAX  // "max"
	MIN  // "min"
	TO   // "to"
	BOOL // "true" or "false"
)

var keywords = map[string]Type{
	"as":    AS,
	"if":    IF,
	"in":    IN,
	"max":   MAX,
	"min":   MIN,
	"to":    TO,
	"true":  BOOL,
	"false": BOOL,
}

// LookupWordToken returns the token for the given word. It checks against a
// static list of reserved keywords, and if the input matches, returns the
// corresponding type. If it doesn't, an WORD token is returned.
func LookupWordToken(w string) Token {
	if tt, ok := keywords[w]; ok {
		return Token{
			Type:  tt,
			Value: w,
		}
	}

	return Token{
		Type:  WORD,
		Value: w,
	}
}
