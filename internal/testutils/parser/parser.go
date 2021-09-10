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

// Package parser contains the parsing primitives needed to process the
// datadriven tests.
package parser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/irfansharif/solver/internal/testutils/parser/ast"
	"github.com/irfansharif/solver/internal/testutils/parser/lexer"
	"github.com/irfansharif/solver/internal/testutils/parser/token"
	"github.com/stretchr/testify/require"
)

// Parser exposes a set of parsing primitives to process the datadriven tests.
type Parser struct {
	lexer *lexer.Lexer
	cur   token.Token

	tb     testing.TB
	trying bool // whether we're currently under a try closure
	failed bool // whether the try closure has failed
}

// New initializes a new parser for the given input.
func New(tb testing.TB, input string) *Parser {
	p := &Parser{tb: tb, lexer: lexer.New(input)}
	p.cur = p.lexer.Next() // stage the current token
	return p
}

// ---------------------------------------------------------------- Token types.

// Digits = Digit { Digit } .
func (p *Parser) Digits() int {
	digits := p.cur.Value
	p.eat(token.DIGITS)
	n, err := strconv.Atoi(digits)
	require.Nil(p, err)
	return n
}

// Word = Letter { Letter } .
func (p *Parser) Word() string {
	word := p.cur.Value
	p.eat(token.WORD)
	return word
}

// Boolean = "true" | "false" .
func (p *Parser) Boolean() bool {
	boolean := p.cur.Value
	p.eat(token.BOOL)
	b, err := strconv.ParseBool(boolean)
	require.Nil(p, err)
	return b
}

// ----------------------------------------------------------------- Base types.

// Identifier = Word .
func (p *Parser) Identifier() string {
	return p.Word()
}

// Number = [ "-" ] Digits .
func (p *Parser) Number() int {
	negative := p.match(token.MINUS)
	if negative {
		p.eat(token.MINUS)
	}

	number := p.Digits()
	if negative {
		number *= -1
	}
	return number
}

// Domain = "[" Number "," Number "]" .
func (p *Parser) Domain() *ast.Domain {
	domain := &ast.Domain{}
	p.eat(token.LBRACKET)
	domain.LowerBound = p.Number()
	p.eat(token.COMMA)
	domain.UpperBound = p.Number()
	p.eat(token.RBRACKET)
	return domain
}

// Variable = Identifier | Letter "to" Letter .
func (p *Parser) Variable() string {
	first := p.Identifier()
	if !p.match(token.TO) {
		return first
	}
	p.eat(token.TO)
	second := p.Identifier()
	require.Lenf(p, first, 1, "expected single letter, got %s", first)
	require.Lenf(p, second, 1, "expected single letter, got %s", second)
	return fmt.Sprintf("%s to %s", first, second)
}

// Interval = Identifier "as" "[" Identifier "," Identifier "|" Identifier "]" .
func (p *Parser) Interval() *ast.Interval {
	interval := &ast.Interval{}
	interval.Name = p.Identifier()
	p.eat(token.AS, token.LBRACKET)
	interval.Start = p.Identifier()
	p.eat(token.COMMA)
	interval.End = p.Identifier()
	p.eat(token.PIPE)
	interval.Size = p.Identifier()
	p.eat(token.RBRACKET)
	return interval
}

// LinearTerm = { Digits } Identifier | Digits .
func (p *Parser) LinearTerm() *ast.LinearTerm {
	term := &ast.LinearTerm{}

	var digits int
	if p.try(func() { digits = p.Digits() }) {
		term.Coefficient = digits
	} else {
		term.Coefficient = 1
		term.Variable = p.Identifier()
		return term
	}

	var variable string
	if p.try(func() { variable = p.Identifier() }) {
		term.Variable = variable
	}

	return term
}

// LinearExpr = [ "-" ] LinearTerm { ("+" | "-") LinearTerm } | "Σ" "(" Variables ")" .
func (p *Parser) LinearExpr() *ast.LinearExpr {
	expr := &ast.LinearExpr{}
	if p.match(token.SUM) {
		p.eat(token.SUM, token.LPAREN)
		variables := p.Variables()
		p.eat(token.RPAREN)
		for _, variable := range variables {
			expr.LinearTerms = append(expr.LinearTerms, &ast.LinearTerm{
				Coefficient: 1,
				Variable:    variable,
			})
		}
		return expr
	}

	negative := p.match(token.MINUS)
	if negative {
		p.eat(token.MINUS)
	}

	term := p.LinearTerm()
	if negative {
		term.Coefficient *= -1
	}
	expr.LinearTerms = append(expr.LinearTerms, term)

	for {
		if !p.match(token.PLUS, token.MINUS) {
			break
		}

		negative = p.match(token.MINUS)
		p.eat(p.cur.Type)

		term = p.LinearTerm()
		if negative {
			term.Coefficient *= -1
		}
		expr.LinearTerms = append(expr.LinearTerms, term)
	}

	return expr
}

// IntervalDemand = Identifier ":" Number .
func (p *Parser) IntervalDemand() *ast.IntervalDemand {
	demand := &ast.IntervalDemand{}
	demand.Name = p.Identifier()
	p.eat(token.COLON)
	demand.Demand = p.Number()
	return demand
}

// ----------------------------------------------------------------- List types.

// Booleans = Boolean { "," Boolean } .
func (p *Parser) Booleans() []bool {
	var booleans []bool
	booleans = append(booleans, p.Boolean())

	for {
		if !p.match(token.COMMA) {
			break
		}

		p.eat(token.COMMA)
		booleans = append(booleans, p.Boolean())
	}

	return booleans
}

// Numbers = Number { "," Number } .
func (p *Parser) Numbers() []int {
	var numbers []int
	numbers = append(numbers, p.Number())

	for {
		if !p.match(token.COMMA) {
			break
		}

		p.eat(token.COMMA)
		numbers = append(numbers, p.Number())
	}

	return numbers
}

// Domains = Domain { "∪" Domain } .
func (p *Parser) Domains() []*ast.Domain {
	var domains []*ast.Domain
	domains = append(domains, p.Domain())

	for {
		if !p.match(token.UNION) {
			break
		}

		p.eat(token.UNION)
		domains = append(domains, p.Domain())
	}

	return domains
}

// Variables = Variable { "," Variable } .
func (p *Parser) Variables() []string {
	var variables []string
	variables = append(variables, p.Variable())

	for {
		if !p.match(token.COMMA) {
			break
		}

		p.eat(token.COMMA)
		variables = append(variables, p.Variable())
	}

	var expanded []string
	letter := func(s string) rune { return []rune(s)[0] }
	for _, v := range variables {
		parts := strings.Split(v, " to ")
		if len(parts) == 1 {
			expanded = append(expanded, v)
			continue
		}

		start, end := letter(parts[0]), letter(parts[1])
		if end < start {
			start, end = end, start
		}
		for c := start; c <= end; c++ {
			expanded = append(expanded, string(c))
		}
	}

	return expanded
}

// Intervals = Interval { "," Interval } .
func (p *Parser) Intervals() []*ast.Interval {
	var intervals []*ast.Interval
	intervals = append(intervals, p.Interval())

	for {
		if !p.match(token.COMMA) {
			break
		}

		p.eat(token.COMMA)
		intervals = append(intervals, p.Interval())
	}

	return intervals
}

// LinearExprs = LinearExpr { "," LinearExpr } .
func (p *Parser) LinearExprs() []*ast.LinearExpr {
	var exprs []*ast.LinearExpr
	exprs = append(exprs, p.LinearExpr())

	for {
		if !p.match(token.COMMA) {
			break
		}

		p.eat(token.COMMA)
		exprs = append(exprs, p.LinearExpr())
	}

	return exprs
}

// IntervalDemands = IntervalDemand {"," IntervalDemand } .
func (p *Parser) IntervalDemands() []*ast.IntervalDemand {
	var demands []*ast.IntervalDemand
	demands = append(demands, p.IntervalDemand())

	for {
		if !p.match(token.COMMA) {
			break
		}

		p.eat(token.COMMA)
		demands = append(demands, p.IntervalDemand())
	}

	return demands
}

// --------------------------------------------------------- List of list types.

// NumbersList = "[" Numbers "]" { "∪" "[" Numbers "]" } .
func (p *Parser) NumbersList() [][]int {
	p.eat(token.LBRACKET)
	var array [][]int
	array = append(array, p.Numbers())
	p.eat(token.RBRACKET)

	for {
		if !p.match(token.UNION) {
			break
		}

		p.eat(token.UNION, token.LBRACKET)
		array = append(array, p.Numbers())
		p.eat(token.RBRACKET)
	}

	return array
}

// BooleansList  = "[" Booleans "]" { "∪" "[" Booleans "]" } .
func (p *Parser) BooleansList() [][]bool {
	p.eat(token.LBRACKET)
	var array [][]bool
	array = append(array, p.Booleans())
	p.eat(token.RBRACKET)
	for {
		if !p.match(token.UNION) {
			break
		}

		p.eat(token.UNION, token.LBRACKET)
		array = append(array, p.Booleans())
		p.eat(token.RBRACKET)
	}

	return array
}

// ------------------------------------------------------------- Argument types.

// AssignmentsArgument = "[" Variables "]" ("∈" | "∉") (NumbersList | BooleanList) .
func (p *Parser) AssignmentsArgument() ast.Argument {
	argument := &ast.AssignmentsArgument{}
	p.eat(token.LBRACKET)
	argument.Variables = p.Variables()
	p.eat(token.RBRACKET)
	if !p.match(token.EXISTS, token.NEXISTS) {
		p.Fatalf("expected either %s or %s token, got %q (%s)",
			token.EXISTS, token.NEXISTS, p.cur.Value, p.cur.Type)
	}
	argument.In = p.match(token.EXISTS)
	p.eat(p.cur.Type)

	var numbers [][]int
	if p.try(func() { numbers = p.NumbersList() }) {
		argument.AllowedIntVarAssignments = numbers
		return argument
	}

	argument.AllowedLiteralAssignments = p.BooleansList()
	return argument
}

// BinaryOpArgument = Identifier ( "/" | "%" | "*" ) Identifier "==" Identifier .
func (p *Parser) BinaryOpArgument() ast.Argument {
	argument := &ast.BinaryOpArgument{}
	argument.Left = p.Identifier()

	if !p.match(token.SLASH, token.MOD, token.ASTERISK) {
		p.Fatalf("expected one of %s, %s, or %s tokens, got %q (%s)",
			token.SLASH, token.MOD, token.ASTERISK, p.cur.Value, p.cur.Type)
	}
	argument.Op = p.cur.Value
	p.eat(p.cur.Type)
	argument.Right = p.Identifier()
	p.eat(token.EQ)
	argument.Target = p.Identifier()
	return argument
}

// ConstantsArgument = Variables "==" Number .
func (p *Parser) ConstantsArgument() ast.Argument {
	argument := &ast.ConstantsArgument{}
	argument.Variables = p.Variables()
	p.eat(token.EQ)
	argument.Constant = p.Number()
	return argument
}

// CumulativeArgument = IntervalDemands "|" Number .
func (p *Parser) CumulativeArgument() ast.Argument {
	argument := &ast.CumulativeArgument{}
	argument.IntervalDemands = p.IntervalDemands()
	p.eat(token.PIPE)
	argument.Capacity = p.Number()
	return argument
}

// DomainArgument = ( Variables | LinearExprs ) "in" Domains .
func (p *Parser) DomainArgument() ast.Argument {
	argument := &ast.DomainArgument{}

	var variables []string
	if p.try(func() { variables = p.Variables() }) {
		argument.Variables = variables
	} else {
		argument.LinearExprs = p.LinearExprs()
	}

	p.eat(token.IN)
	argument.Domains = p.Domains()
	return argument
}

// ElementArgument = Identifier "==" "[" Variables "]" "[" Identifier "]" .
func (p *Parser) ElementArgument() ast.Argument {
	argument := &ast.ElementArgument{}
	argument.Target = p.Identifier()
	p.eat(token.EQ, token.LBRACKET)
	argument.Variables = p.Variables()
	p.eat(token.RBRACKET, token.LBRACKET)
	argument.Index = p.Identifier()
	p.eat(token.RBRACKET)
	return argument
}

// ImplicationArgument = Identifier "→"  Identifier .
func (p *Parser) ImplicationArgument() ast.Argument {
	argument := &ast.ImplicationArgument{}
	argument.Left = p.Identifier()
	p.eat(token.IMPL)
	argument.Right = p.Identifier()
	return argument
}

// IntervalsArgument = Intervals .
func (p *Parser) IntervalsArgument() ast.Argument {
	argument := &ast.IntervalsArgument{}
	argument.Intervals = p.Intervals()
	return argument
}

// KArgument = Variables "|" Digits .
func (p *Parser) KArgument() ast.Argument {
	argument := &ast.KArgument{}
	argument.Variables = p.Variables()
	p.eat(token.PIPE)
	argument.K = p.Digits()
	return argument
}

// LinearEqualityArgument = LinearExpr "==" ("max" | "min") "(" LinearExprs ")" .
func (p *Parser) LinearEqualityArgument() ast.Argument {
	argument := &ast.LinearEqualityArgument{}
	argument.Target = p.LinearExpr()
	p.eat(token.EQ)
	if !p.match(token.MAX, token.MIN) {
		p.Fatalf("expected either %s or %s token type, got %q (%s)",
			token.MAX, token.MIN, p.cur.Value, p.cur.Type)
	}
	argument.Op = p.cur.Value
	p.eat(p.cur.Type, token.LPAREN)
	argument.Exprs = p.LinearExprs()
	p.eat(token.RPAREN)
	return argument
}

// LinearExprsArgument = LinearExprs .
func (p *Parser) LinearExprsArgument() ast.Argument {
	argument := &ast.LinearExprsArgument{}
	argument.Exprs = p.LinearExprs()
	return argument
}

// NonOverlapping2DArgument = "[" Variables "]" "," "[" Variables "]" "," Boolean .
func (p *Parser) NonOverlapping2DArgument() ast.Argument {
	argument := &ast.NonOverlapping2DArgument{}
	p.eat(token.LBRACKET)
	argument.XVariables = p.Variables()
	p.eat(token.RBRACKET, token.COMMA, token.LBRACKET)
	argument.YVariables = p.Variables()
	p.eat(token.RBRACKET, token.COMMA)
	argument.BoxesWithNoAreaCanOverlap = p.Boolean()
	return argument
}

// VariableEqualityArgument = Identifier "==" ("max" | "min" ) "(" Variables ")" .
func (p *Parser) VariableEqualityArgument() ast.Argument {
	argument := &ast.VariableEqualityArgument{}
	argument.Target = p.Identifier()
	p.eat(token.EQ)
	if !p.match(token.MAX, token.MIN) {
		p.Fatalf("expected either %s or %s token type, got %q (%s)",
			token.MAX, token.MIN, p.cur.Value, p.cur.Type)
	}
	argument.Op = p.cur.Value
	p.eat(p.cur.Type, token.LPAREN)
	argument.Variables = p.Variables()
	p.eat(token.RPAREN)
	return argument
}

// VariablesArgument = Variables .
func (p *Parser) VariablesArgument() ast.Argument {
	argument := &ast.VariablesArgument{}
	argument.Variables = p.Variables()
	return argument
}

// -------------------------------------------------- Statement component types.

// Argument = AssignmentsArgument
//          | BinaryOpArgument
//          | ConstantsArgument
//          | CumulativeArgument
//          | DomainArgument
//          | ElementArgument
//          | IntervalsArgument
//          | ImplicationArgument
//          | KArgument
//          | LinearEqualityArgument
//          | LinearExprsArgument
//          | NonOverlapping2DArgument
//          | VariableEqualityArgument
//          | VariablesArgument .
func (p *Parser) Argument() ast.Argument {
	fns := []func() ast.Argument{
		p.AssignmentsArgument,
		p.BinaryOpArgument,
		p.ConstantsArgument,
		p.CumulativeArgument,
		p.DomainArgument,
		p.ElementArgument,
		p.IntervalsArgument,
		p.ImplicationArgument,
		p.KArgument,
		p.LinearEqualityArgument,
		p.NonOverlapping2DArgument,
		p.VariableEqualityArgument,

		p.VariablesArgument, // there's ambiguity; give precedence to parsing variables argument
		p.LinearExprsArgument,
	}

	for _, fn := range fns {
		var argument ast.Argument
		if p.try(func() {
			argument = fn()
			if !p.match(token.RPAREN) {
				p.Fatalf("expected %s token, got %s (value=%q)",
					token.RPAREN.String(), p.cur.Type.String(), p.cur.Value)
			}
		}) {
			return argument
		}
	}

	p.Fatal("expected to match an argument type")
	return nil
}

// Receiver = Identifier .
func (p *Parser) Receiver() string {
	return p.Identifier()
}

// Method = Identifier { "-" | Identifier | Digits } .
func (p *Parser) Method() ast.Method {
	var out strings.Builder
	identifier := p.Identifier()
	out.WriteString(identifier)

	for p.match(token.MINUS, token.WORD, token.DIGITS) {
		out.WriteString(p.cur.Value)
		p.eat(p.cur.Type)
	}

	methodStr := out.String()
	method, ok := ast.LookupMethod(methodStr)
	require.Truef(p, ok, "unrecognized method: %s", methodStr)
	return method
}

// Enforcement = "if" Variables .
func (p *Parser) Enforcement() *ast.Enforcement {
	p.eat(token.IF)
	enforcement := &ast.Enforcement{}
	enforcement.Variables = p.Variables()
	return enforcement
}

// Statement = Receiver "." Method "(" [ Argument ] ")" [ Enforcement ] .
func (p *Parser) Statement() *ast.Statement {
	stmt := &ast.Statement{}
	stmt.Receiver = p.Receiver()
	p.eat(token.DOT)
	stmt.Method = p.Method()
	p.eat(token.LPAREN)
	if !p.match(token.RPAREN) {
		stmt.Argument = p.Argument()
	}
	p.eat(token.RPAREN)
	if !p.match(token.EOF) {
		stmt.Enforcement = p.Enforcement()
	}
	p.eat(token.EOF)
	return stmt
}

// EOF returns true if we're at the end of the input.
func (p *Parser) EOF() bool {
	return p.match(token.EOF)
}

// try attempts to parse using the given closure; if it fails, it resets the
// underlying cursors to their original states. The returned boolean indicates
// whether the parsing attempt was successful. It's invalid to use state from
// the closure if the attempt failed.
func (p *Parser) try(parse func()) (success bool) {
	idx, cur := p.lexer.Index(), p.cur
	trying, failed := p.trying, p.failed
	defer func() {
		p.trying, p.failed = trying, failed
		if !success {
			p.lexer.Reposition(idx)
			p.cur = cur
		}
	}()

	p.trying, p.failed = true, false
	parse()
	return !p.failed
}

// eat asserts that the current and subsequent tokens are of the given types,
// consuming them as it does. It moves the cursor over past the last token.
func (p *Parser) eat(ts ...token.Type) {
	for _, t := range ts {
		require.Truef(p, p.match(t), "expected %s token, got %s (value=%s)", t.String(), p.cur.Type.String(), p.cur.Value)
		p.cur = p.lexer.Next()
	}
}

// match returns whether the current token is one of the given types.
func (p *Parser) match(ts ...token.Type) bool {
	match := false
	for _, t := range ts {
		match = match || p.cur.Type == t
	}
	return match
}

// testingT is n wrapper around testing.T; it's used by the parser to collect
// intercept parsing errors with arbitrary look-ahead.
type testingT interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fail()
	FailNow()
}

var _ testingT = &Parser{}

// Errorf is parting of the testingT interface.
func (p *Parser) Errorf(format string, args ...interface{}) {
	if !p.trying {
		p.tb.Logf(format, args...)
	}
	p.Fail()
}

// Fatalf is parting of the testingT interface.
func (p *Parser) Fatalf(format string, args ...interface{}) {
	if !p.trying {
		p.tb.Logf(format, args...)
	}
	p.FailNow()
}

// Fatal is parting of the testingT interface.
func (p *Parser) Fatal(args ...interface{}) {
	if !p.trying {
		p.tb.Log(args...)
	}
	p.FailNow()
}

// Fail is parting of the testingT interface.
func (p *Parser) Fail() {
	if p.trying {
		p.failed = true
		return
	}

	p.tb.Fail()
}

// FailNow is parting of the testingT interface.
func (p *Parser) FailNow() {
	if p.trying {
		p.failed = true
		return
	}
	p.tb.FailNow()
}
