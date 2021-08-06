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

	"github.com/irfansharif/solver/internal/testutils/parser/ast"
	"github.com/irfansharif/solver/internal/testutils/parser/lexer"
	"github.com/irfansharif/solver/internal/testutils/parser/token"
)

// Parser exposes a set of parsing primitives to process the datadriven tests.
type Parser struct {
	l   *lexer.Lexer
	cur token.Token
}

// New initializes a new parser for the given input.
func New(s string) *Parser {
	p := &Parser{l: lexer.New(s)}
	p.cur = p.l.Next() // stage cur
	return p
}

// XXX: Use testing.T throughout instead?

// ------------------------------- TOKEN TYPES.

// Digits = Digit { Digit } .
func (p *Parser) Digits() (int, error) {
	val := p.cur.Value
	if err := p.eat(token.DIGITS); err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

// Word = Letter { Letter } .
func (p *Parser) Word() (string, error) {
	result := p.cur.Value
	if err := p.eat(token.WORD); err != nil {
		return "", err
	}
	return result, nil
}

// Boolean = "true" | "false" .
func (p *Parser) Boolean() (bool, error) {
	val := p.cur.Value
	if err := p.eat(token.BOOL); err != nil {
		return false, err
	}
	return strconv.ParseBool(val)
}

// ------------------------------- BASE TYPES.

// Identifier = Word .
func (p *Parser) Identifier() (string, error) {
	return p.Word()
}

// Number = [ "-" ] Digits .
func (p *Parser) Number() (int, error) {
	negative := p.cur.Type == token.MINUS
	if negative {
		_ = p.eat(token.MINUS)
	}

	number, err := p.Digits()
	if err != nil {
		return 0, err
	}
	if negative {
		number *= -1
	}
	return number, nil
}

// Domain = "[" Number "," Number "]" .
func (p *Parser) Domain() (*ast.Domain, error) {
	domain := &ast.Domain{}
	var err error
	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	if domain.LowerBound, err = p.Number(); err != nil {
		return nil, err
	}
	if err = p.eat(token.COMMA); err != nil {
		return nil, err
	}
	if domain.UpperBound, err = p.Number(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}
	return domain, nil
}

// Variable = Identifier | Letter "to" Letter .
func (p *Parser) Variable() (string, error) {
	var out strings.Builder
	first, err := p.Identifier()
	if err != nil {
		return "", err
	}
	out.WriteString(first)

	if p.cur.Type != token.TO {
		return out.String(), nil
	}

	out.WriteString(" to ")
	_ = p.eat(token.TO)

	second, err := p.Identifier()
	if err != nil {
		return "", err
	}
	out.WriteString(second)

	if len(first) != 1 {
		return "", fmt.Errorf("expected single letter, got %s", first)
	}
	if len(second) != 1 {
		return "", fmt.Errorf("expected single letter, got %s", second)
	}
	return out.String(), nil
}

// Interval = Identifier "as" "[" Identifier "," Identifier "|" Identifier "]" .
func (p *Parser) Interval() (*ast.Interval, error) {
	interval := &ast.Interval{}
	var err error
	interval.Name, err = p.Identifier()
	if err != nil {
		return nil, err
	}

	if err = p.eat(token.AS); err != nil {
		return nil, err
	}

	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}

	interval.Start, err = p.Identifier()
	if err != nil {
		return nil, err
	}

	if err = p.eat(token.COMMA); err != nil {
		return nil, err
	}

	interval.End, err = p.Identifier()
	if err != nil {
		return nil, err
	}

	if err = p.eat(token.PIPE); err != nil {
		return nil, err
	}

	interval.Size, err = p.Identifier()
	if err != nil {
		return nil, err
	}

	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}

	return interval, nil
}

// LinearTerm = { Digits } Identifier | Digits .
func (p *Parser) LinearTerm() (*ast.LinearTerm, error) {
	term := &ast.LinearTerm{}
	var err error
	digits := p.try(func() error {
		term.Coefficient, err = p.Digits()
		if err != nil {
			return err
		}
		return nil
	})
	if !digits {
		term.Coefficient = 1
		if term.Variable, err = p.Identifier(); err != nil {
			return nil, err
		}
		return term, nil
	}

	_ = p.try(func() error {
		if term.Variable, err = p.Identifier(); err != nil {
			return err
		}
		return nil
	})

	return term, nil
}

// LinearExpr = [ "-" ] LinearTerm { ("+" | "-") LinearTerm } | "Σ" "(" Variables ")" .
func (p *Parser) LinearExpr() (*ast.LinearExpr, error) {
	expr := &ast.LinearExpr{}
	if p.cur.Type == token.SUM {
		if err := p.eat(token.SUM); err != nil {
			return nil, err
		}
		if err := p.eat(token.LPAREN); err != nil {
			return nil, err
		}
		variables, err := p.Variables()
		if err != nil {
			return nil, err
		}
		if err := p.eat(token.RPAREN); err != nil {
			return nil, err
		}

		for _, variable := range variables {
			expr.LinearTerms = append(expr.LinearTerms, &ast.LinearTerm{
				Coefficient: 1,
				Variable:    variable,
			})
		}
		return expr, nil
	}

	negative := p.cur.Type == token.MINUS
	if negative {
		_ = p.eat(token.MINUS)
	}

	term, err := p.LinearTerm()
	if err != nil {
		return nil, err
	}
	if negative {
		term.Coefficient *= -1
	}
	expr.LinearTerms = append(expr.LinearTerms, term)

	for {
		if p.cur.Type != token.PLUS && p.cur.Type != token.MINUS {
			break
		}

		negative = p.cur.Type == token.MINUS
		_ = p.eat(p.cur.Type)

		term, err = p.LinearTerm()
		if err != nil {
			return nil, err
		}
		if negative {
			term.Coefficient *= -1
		}
		expr.LinearTerms = append(expr.LinearTerms, term)
	}

	return expr, nil
}

// IntervalDemand = Identifier ":" Number .
func (p *Parser) IntervalDemand() (*ast.IntervalDemand, error) {
	demand := &ast.IntervalDemand{}
	var err error
	if demand.Name, err = p.Identifier(); err != nil {
		return nil, err
	}

	if err = p.eat(token.COLON); err != nil {
		return nil, err
	}

	if demand.Demand, err = p.Number(); err != nil {
		return nil, err
	}

	return demand, nil
}

// ------------------------------- LIST TYPES.

// Booleans = Boolean { "," Boolean } .
func (p *Parser) Booleans() ([]bool, error) {
	var booleans []bool
	boolean, err := p.Boolean()
	if err != nil {
		return nil, err
	}
	booleans = append(booleans, boolean)

	for {
		if p.cur.Type != token.COMMA {
			break
		}

		_ = p.eat(token.COMMA)
		boolean, err = p.Boolean()
		if err != nil {
			return nil, err
		}
		booleans = append(booleans, boolean)
	}

	return booleans, nil
}

// Numbers = Number { "," Number } .
func (p *Parser) Numbers() ([]int, error) {
	var numbers []int
	number, err := p.Number()
	if err != nil {
		return nil, err
	}
	numbers = append(numbers, number)

	for {
		if p.cur.Type != token.COMMA {
			break
		}

		_ = p.eat(token.COMMA)
		number, err = p.Number()
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}

	return numbers, nil
}

// Domains = Domain { "∪" Domain } .
func (p *Parser) Domains() ([]*ast.Domain, error) {
	var domains []*ast.Domain
	domain, err := p.Domain()
	if err != nil {
		return nil, err
	}
	domains = append(domains, domain)

	for {
		if p.cur.Type != token.UNION {
			break
		}

		_ = p.eat(token.UNION)
		domain, err := p.Domain()
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

// Variables = Variable { "," Variable } .
func (p *Parser) Variables() ([]string, error) {
	var variables []string
	variable, err := p.Variable()
	if err != nil {
		return nil, err
	}
	variables = append(variables, variable)

	for {
		if p.cur.Type != token.COMMA {
			break
		}

		_ = p.eat(token.COMMA)
		variable, err = p.Variable()
		if err != nil {
			return nil, err
		}
		variables = append(variables, variable)
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

	return expanded, nil
}

// Intervals = Interval { "," Interval } .
func (p *Parser) Intervals() ([]*ast.Interval, error) {
	var intervals []*ast.Interval
	interval, err := p.Interval()
	if err != nil {
		return nil, err
	}
	intervals = append(intervals, interval)

	for {
		if p.cur.Type != token.COMMA {
			break
		}

		_ = p.eat(token.COMMA)
		interval, err = p.Interval()
		if err != nil {
			return nil, err
		}
		intervals = append(intervals, interval)
	}

	return intervals, nil
}

// LinearExprs = LinearExpr { "," LinearExpr } .
func (p *Parser) LinearExprs() ([]*ast.LinearExpr, error) {
	var exprs []*ast.LinearExpr
	expr, err := p.LinearExpr()
	if err != nil {
		return nil, err
	}
	exprs = append(exprs, expr)

	for {
		if p.cur.Type != token.COMMA {
			break
		}

		_ = p.eat(token.COMMA)
		expr, err = p.LinearExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

// IntervalDemands = IntervalDemand {"," IntervalDemand } .
func (p *Parser) IntervalDemands() ([]*ast.IntervalDemand, error) {
	var demands []*ast.IntervalDemand
	demand, err := p.IntervalDemand()
	if err != nil {
		return nil, err
	}
	demands = append(demands, demand)

	for {
		if p.cur.Type != token.COMMA {
			break
		}

		_ = p.eat(token.COMMA)
		demand, err := p.IntervalDemand()
		if err != nil {
			return nil, err
		}
		demands = append(demands, demand)
	}

	return demands, nil
}

// ------------------------------- LIST TYPES.

// NumbersList = "[" Numbers "]" { "∪" "[" Numbers "]" } .
func (p *Parser) NumbersList() ([][]int, error) {
	if err := p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	var array [][]int
	numbers, err := p.Numbers()
	if err != nil {
		return nil, err
	}
	array = append(array, numbers)
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}

	for {
		if p.cur.Type != token.UNION {
			break
		}

		_ = p.eat(token.UNION)

		if err = p.eat(token.LBRACKET); err != nil {
			return nil, err
		}
		numbers, err = p.Numbers()
		if err != nil {
			return nil, err
		}
		array = append(array, numbers)
		if err = p.eat(token.RBRACKET); err != nil {
			return nil, err
		}
	}

	return array, nil
}

// BooleansList  = "[" Booleans "]" { "∪" "[" Booleans "]" } .
func (p *Parser) BooleansList() ([][]bool, error) {
	if err := p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	var array [][]bool
	booleans, err := p.Booleans()
	if err != nil {
		return nil, err
	}
	array = append(array, booleans)
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}

	for {
		if p.cur.Type != token.UNION {
			break
		}

		_ = p.eat(token.UNION)

		if err = p.eat(token.LBRACKET); err != nil {
			return nil, err
		}
		booleans, err = p.Booleans()
		if err != nil {
			return nil, err
		}
		array = append(array, booleans)
		if err = p.eat(token.RBRACKET); err != nil {
			return nil, err
		}
	}

	return array, nil
}

// ------------------------------- ARGUMENT TYPES.

// AssignmentsArgument = "[" Variables "]" ("∈" | "∉") (NumbersList | BooleanList) .
func (p *Parser) AssignmentsArgument() (ast.Argument, error) {
	argument := &ast.AssignmentsArgument{}
	var err error
	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}

	if argument.Variables, err = p.Variables(); err != nil {
		return nil, err
	}

	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}

	if p.cur.Type != token.EXISTS && p.cur.Type != token.NEXISTS {
		return nil, fmt.Errorf("expected either %s or %s token, got %q (%s)",
			token.EXISTS, token.NEXISTS, p.cur.Value, p.cur.Type)
	}
	argument.In = p.cur.Type == token.EXISTS
	_ = p.eat(p.cur.Type)

	numbers := p.try(func() error {
		argument.NumbersList, err = p.NumbersList()
		if err != nil {
			return err
		}
		return nil
	})
	if numbers {
		return argument, nil
	}

	if argument.BooleansList, err = p.BooleansList(); err != nil {
		return nil, err
	}
	return argument, nil
}

// BinaryOpArgument = Identifier ( "/" | "%" | "*" ) Identifier "==" Identifier .
func (p *Parser) BinaryOpArgument() (ast.Argument, error) {
	argument := &ast.BinaryOpArgument{}
	var err error
	if argument.Left, err = p.Identifier(); err != nil {
		return nil, err
	}

	if p.cur.Type != token.SLASH && p.cur.Type != token.MOD && p.cur.Type != token.ASTERISK {
		return nil, fmt.Errorf("expected one of %s, %s, or %s tokens, got %q (%s)",
			token.SLASH, token.MOD, token.ASTERISK, p.cur.Value, p.cur.Type)
	}
	argument.Op = p.cur.Value
	_ = p.eat(p.cur.Type)

	if argument.Right, err = p.Identifier(); err != nil {
		return nil, err
	}

	if err = p.eat(token.EQ); err != nil {
		return nil, err
	}
	if argument.Target, err = p.Identifier(); err != nil {
		return nil, err
	}
	return argument, err
}

// ConstantsArgument = Variables "==" Number .
func (p *Parser) ConstantsArgument() (ast.Argument, error) {
	argument := &ast.ConstantsArgument{}
	var err error
	if argument.Variables, err = p.Variables(); err != nil {
		return nil, err
	}

	if err = p.eat(token.EQ); err != nil {
		return nil, err
	}

	if argument.Constant, err = p.Number(); err != nil {
		return nil, err
	}
	return argument, err
}

// CumulativeArgument = IntervalDemands "|" Number .
func (p *Parser) CumulativeArgument() (ast.Argument, error) {
	argument := &ast.CumulativeArgument{}
	var err error
	if argument.IntervalDemands, err = p.IntervalDemands(); err != nil {
		return nil, err
	}

	if err = p.eat(token.PIPE); err != nil {
		return nil, err
	}

	if argument.Capacity, err = p.Number(); err != nil {
		return nil, err
	}
	return argument, err
}

// DomainArgument = ( Variables | LinearExprs ) "in" Domains .
func (p *Parser) DomainArgument() (ast.Argument, error) {
	argument := &ast.DomainArgument{}
	var err error
	variables := p.try(func() error {
		if argument.Variables, err = p.Variables(); err != nil {
			return err
		}
		return nil
	})
	if !variables {
		if argument.LinearExprs, err = p.LinearExprs(); err != nil {
			return nil, err
		}
	}

	if err = p.eat(token.IN); err != nil {
		return nil, err
	}

	if argument.Domains, err = p.Domains(); err != nil {
		return nil, err
	}
	return argument, err
}

// ElementArgument = Identifier "==" "[" Variables "]" "[" Identifier "]" .
func (p *Parser) ElementArgument() (ast.Argument, error) {
	argument := &ast.ElementArgument{}
	var err error
	if argument.Target, err = p.Identifier(); err != nil {
		return nil, err
	}
	if err = p.eat(token.EQ); err != nil {
		return nil, err
	}
	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	if argument.Variables, err = p.Variables(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}
	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	if argument.Index, err = p.Identifier(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}
	return argument, err
}

// ImplicationArgument = Identifier "→"  Identifier .
func (p *Parser) ImplicationArgument() (ast.Argument, error) {
	argument := &ast.ImplicationArgument{}
	var err error
	if argument.Left, err = p.Identifier(); err != nil {
		return nil, err
	}

	if err = p.eat(token.IMPL); err != nil {
		return nil, err
	}

	if argument.Right, err = p.Identifier(); err != nil {
		return nil, err
	}
	return argument, err
}

// IntervalsArgument = Intervals .
func (p *Parser) IntervalsArgument() (ast.Argument, error) {
	argument := &ast.IntervalsArgument{}
	var err error
	if argument.Intervals, err = p.Intervals(); err != nil {
		return nil, err
	}
	return argument, nil
}

// KArgument = Variables "|" Digits .
func (p *Parser) KArgument() (ast.Argument, error) {
	argument := &ast.KArgument{}
	var err error
	if argument.Variables, err = p.Variables(); err != nil {
		return nil, err
	}

	if err = p.eat(token.PIPE); err != nil {
		return nil, err
	}

	if argument.K, err = p.Digits(); err != nil {
		return nil, err
	}
	return argument, err
}

// LinearEqualityArgument = LinearExpr "==" ("max" | "min") "(" LinearExprs ")" .
func (p *Parser) LinearEqualityArgument() (ast.Argument, error) {
	argument := &ast.LinearEqualityArgument{}
	var err error
	if argument.Target, err = p.LinearExpr(); err != nil {
		return nil, err
	}

	if err = p.eat(token.EQ); err != nil {
		return nil, err
	}

	if p.cur.Type != token.MAX && p.cur.Type != token.MIN {
		return nil, fmt.Errorf("expected either %s or %s token type, got %q (%s)",
			token.MAX, token.MIN, p.cur.Value, p.cur.Type)
	}
	argument.Op = p.cur.Value
	_ = p.eat(p.cur.Type)

	if err = p.eat(token.LPAREN); err != nil {
		return nil, err
	}
	if argument.Exprs, err = p.LinearExprs(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RPAREN); err != nil {
		return nil, err
	}
	return argument, nil
}

// LinearExprsArgument = LinearExprs .
func (p *Parser) LinearExprsArgument() (ast.Argument, error) {
	argument := &ast.LinearExprsArgument{}
	var err error
	if argument.Exprs, err = p.LinearExprs(); err != nil {
		return nil, err
	}
	return argument, nil
}

// NonOverlapping2DArgument = "[" Variables "]" "," "[" Variables "]" "," Boolean .
func (p *Parser) NonOverlapping2DArgument() (ast.Argument, error) {
	argument := &ast.NonOverlapping2DArgument{}
	var err error
	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	if argument.XVariables, err = p.Variables(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}
	if err = p.eat(token.COMMA); err != nil {
		return nil, err
	}
	if err = p.eat(token.LBRACKET); err != nil {
		return nil, err
	}
	if argument.YVariables, err = p.Variables(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RBRACKET); err != nil {
		return nil, err
	}
	if err = p.eat(token.COMMA); err != nil {
		return nil, err
	}
	if argument.BoxesWithNoAreaCanOverlap, err = p.Boolean(); err != nil {
		return nil, err
	}
	return argument, nil
}

// VariableEqualityArgument = Identifier "==" ("max" | "min" ) "(" Variables ")" .
func (p *Parser) VariableEqualityArgument() (ast.Argument, error) {
	argument := &ast.VariableEqualityArgument{}
	var err error
	if argument.Target, err = p.Identifier(); err != nil {
		return nil, err
	}

	if err = p.eat(token.EQ); err != nil {
		return nil, err
	}

	if p.cur.Type != token.MAX && p.cur.Type != token.MIN {
		return nil, fmt.Errorf("expected either %s or %s token type, got %q (%s)",
			token.MAX, token.MIN, p.cur.Value, p.cur.Type)
	}
	argument.Op = p.cur.Value
	_ = p.eat(p.cur.Type)

	if err = p.eat(token.LPAREN); err != nil {
		return nil, err
	}
	if argument.Variables, err = p.Variables(); err != nil {
		return nil, err
	}
	if err = p.eat(token.RPAREN); err != nil {
		return nil, err
	}
	return argument, nil
}

// VariablesArgument = Variables .
func (p *Parser) VariablesArgument() (ast.Argument, error) {
	argument := &ast.VariablesArgument{}
	var err error
	if argument.Variables, err = p.Variables(); err != nil {
		return nil, err
	}
	return argument, nil
}

// ------------------------------- STATEMENT COMPONENTS.

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
func (p *Parser) Argument() (ast.Argument, error) {
	fns := []func() (ast.Argument, error){
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

		p.VariablesArgument, // there's ambiguity: give precedence to parsing variables argument
		p.LinearExprsArgument,
	}

	for _, fn := range fns {
		var argument ast.Argument
		var err error
		if found := p.try(func() error {
			argument, err = fn()
			if err != nil {
				return err
			}

			if p.cur.Type != token.RPAREN {
				return fmt.Errorf("expected %s token, got %q (%s)", token.RPAREN, p.cur.Value, p.cur.Type)
			}
			return nil
		}); found {
			return argument, nil
		}
	}

	return nil, fmt.Errorf("expected to match an argument type")
}

// Receiver = Identifier .
func (p *Parser) Receiver() (string, error) {
	return p.Identifier()
}

// Method = Identifier { "-" | Identifier | Digits } .
func (p *Parser) Method() (ast.Method, error) {
	var out strings.Builder
	identifier, err := p.Identifier()
	if err != nil {
		return ast.Unrecognized, err
	}
	out.WriteString(identifier)

	for p.cur.Type == token.MINUS || p.cur.Type == token.WORD || p.cur.Type == token.DIGITS {
		out.WriteString(p.cur.Value)
		_ = p.eat(p.cur.Type)
	}

	method, ok := ast.LookupMethod(out.String())
	if !ok {
		return ast.Unrecognized, fmt.Errorf("unrecognized method: %s", out.String())
	}
	return method, nil
}

// Enforcement = "if" Variables .
func (p *Parser) Enforcement() (*ast.Enforcement, error) {
	var err error
	if err = p.eat(token.IF); err != nil {
		return nil, err
	}

	enforcement := &ast.Enforcement{}
	if enforcement.Variables, err = p.Variables(); err != nil {
		return nil, err
	}

	return enforcement, nil
}

// Statement = Receiver "." Method "(" [ Argument ] ")" [ Enforcement ] .
func (p *Parser) Statement() (*ast.Statement, error) {
	stmt := &ast.Statement{}

	var err error
	if stmt.Receiver, err = p.Receiver(); err != nil {
		return nil, err
	}

	if err = p.eat(token.DOT); err != nil {
		return nil, err
	}

	if stmt.Method, err = p.Method(); err != nil {
		return nil, err
	}

	if err = p.eat(token.LPAREN); err != nil {
		return nil, err
	}

	if p.cur.Type != token.RPAREN {
		if stmt.Argument, err = p.Argument(); err != nil {
			return nil, err
		}
	}

	if err = p.eat(token.RPAREN); err != nil {
		return nil, err
	}

	if p.cur.Type != token.EOF {
		if stmt.Enforcement, err = p.Enforcement(); err != nil {
			return nil, err
		}
	}

	if err = p.eat(token.EOF); err != nil {
		return nil, err
	}
	return stmt, nil
}

// try executes the given closure, and if an error is returned, resets the
// underlying cursors to their original states.
func (p *Parser) try(f func() error) bool {
	idx, cur := p.l.Index(), p.cur
	if err := f(); err != nil {
		p.l.Reposition(idx)
		p.cur = cur
		return false
	}
	return true
}

// EOF returns true if we're at the end of the input.
func (p *Parser) EOF() bool {
	return p.cur.Type == token.EOF
}

// eat asserts that the current token is of the given type and moves and cursor
// over to the next token.
func (p *Parser) eat(t token.Type) error {
	if p.cur.Type != t {
		return fmt.Errorf("expected %s token, got %q (%s)", t, p.cur.Value, p.cur.Type)
	}

	p.cur = p.l.Next()
	return nil
}
