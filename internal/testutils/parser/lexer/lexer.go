package lexer

import (
	"github.com/irfansharif/solver/internal/testutils/parser/token"
)

const eof = rune(0)

// Lexer converts an input string into its constituent tokens.
type Lexer struct {
	input []rune
	idx   int
}

// New creates a lexer from the given input.
func New(input string) *Lexer {
	lexer := &Lexer{input: []rune(input)}
	return lexer
}

// TODO(irfansharif): Do we want to capture '¬' for a logical negation? We could
// also use the '!' operator.

// Next returns the next token from the input and moves the current position of
// the lexer ahead.
//
// If encountering the end of the input, token.EOF is returned.
func (l *Lexer) Next() token.Token {
	for isWhitespace(l.rune()) { // skip whitespace and position at index after
		l.move()
	}

	tok := func(tt token.Type, r rune) token.Token {
		return token.Token{Type: tt, Value: string(r)}
	}
	var t token.Token
	switch r := l.rune(); r {
	case eof:
		t = tok(token.EOF, r)
	case '+':
		t = tok(token.PLUS, r)
	case '-':
		t = tok(token.MINUS, r)
	case '*':
		t = tok(token.ASTERISK, r)
	case '/':
		t = tok(token.SLASH, r)
	case '→':
		t = tok(token.IMPL, r)
	case '%':
		t = tok(token.MOD, r)
	case '<':
		t = tok(token.LT, r)
	case '>':
		t = tok(token.GT, r)
	case '∈':
		t = tok(token.EXISTS, r)
	case '∉':
		t = tok(token.NEXISTS, r)
	case '∪':
		t = tok(token.UNION, r)
	case '.':
		t = tok(token.DOT, r)
	case ':':
		t = tok(token.COLON, r)
	case ',':
		t = tok(token.COMMA, r)
	case '|':
		t = tok(token.PIPE, r)
	case 'Σ':
		t = tok(token.SUM, r)
	case '(':
		t = tok(token.LPAREN, r)
	case ')':
		t = tok(token.RPAREN, r)
	case '[':
		t = tok(token.LBRACKET, r)
	case ']':
		t = tok(token.RBRACKET, r)
	case '=':
		if l.peek() == '=' {
			l.move() // move the cursor to the end of the token
			t = token.Token{Type: token.EQ, Value: "=="}
		} else {
			t = tok(token.ILLEGAL, r)
		}
	case '!':
		if l.peek() == '=' {
			l.move() // move the cursor to the end of the token
			t = token.Token{Type: token.NEQ, Value: "!="}
		} else {
			t = tok(token.BANG, r)
		}
	default:
		switch {
		case isLetter(r):
			word := l.word()
			t = token.LookupWordToken(word)
		case isDigit(r):
			digits := l.digits()
			t = token.Token{Type: token.DIGITS, Value: digits}
		default:
			t = tok(token.ILLEGAL, r)
		}
	}

	l.move() // move the cursor past the end of the token
	return t
}

// Index returns the current position of the lexer.
func (l *Lexer) Index() int {
	return l.idx
}

// Reposition positions the lexer to the given index.
func (l *Lexer) Reposition(idx int) {
	if idx < 0 {
		panic("moving past start of input")
	}
	if idx > len(l.input) {
		panic("moving past end of input")
	}

	l.idx = idx
}

// rune returns the rune under examination. If we're at the end of the input,
// eof is returned.
func (l *Lexer) rune() rune {
	if l.idx == len(l.input) {
		return eof
	}
	return l.input[l.idx]
}

// peek returns the next rune from the input without moving the current position
// ahead. If the next position is the end of the input, eof is returned. This is
// symmetric with Lexer.rune.
func (l *Lexer) peek() rune {
	if l.idx+1 == len(l.input) {
		return eof
	}
	return l.input[l.idx+1]
}

// move moves the current position of the lexer up by one. func (l *Lexer)
func (l *Lexer) move() {
	if l.idx+1 > len(l.input) {
		l.idx = len(l.input)
	} else {
		l.idx += 1
	}
}

// word lexes a single word, moving the cursor to the end of the word.
func (l *Lexer) word() string {
	start := l.idx
	for isLetter(l.peek()) {
		l.move()
	}

	return string(l.input[start : l.idx+1])
}

// digits lexes a sequence of digits, moving the cursor to the end of sequence.
func (l *Lexer) digits() string {
	start := l.idx
	for isDigit(l.peek()) {
		l.move()
	}

	return string(l.input[start : l.idx+1])
}

// isWhiteSpace returns true if the rune is a whitespace.
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// isLetter returns true if the rune is a letter.
func isLetter(r rune) bool {
	return isBetween(r, 'a', 'z') || isBetween(r, 'A', 'Z')
}

// isDigit returns true if the rune is a digit.
func isDigit(r rune) bool {
	return isBetween(r, '0', '9')
}

// isBetween returns true if the rune is between (inclusive) the given start and
// end runes.
func isBetween(r, start, end rune) bool {
	return start <= r && r <= end
}
