package lmdlx

import (
	"fmt"
	"unicode/utf8"
)

type item struct {
	Typ itemType `json:"Typ"`
	Val string   `json:"Val"`
	Id  int      `json:"Id"`
}
type lexer struct {
	input string
	start int
	state stateFn
	pos   int
	width int
	items chan item
	id    int
}
type itemType int
type stateFn func(l *lexer) stateFn

//Welcome to the character definition zone
const eof = -1
const (
	itemError itemType = iota
	itemText
	itemAsterisk
	itemUnderscore
	itemBacktick
	itemBlockq
	itemNewline
	itemTilda
	itemEOF
)

//sym pos must = items pos
var sym = []rune{'R', 'T', '*', '_', '`', '>', '\n', '~', 'E'}
var items = []itemType{itemError, itemText, itemAsterisk, itemUnderscore, itemBacktick, itemBlockq, itemNewline, itemTilda, itemEOF}

//hope you had a nice stay

func (i item) String() string {
	switch i.Typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.Val
	}
	if len(i.Val) > 10 {
		return fmt.Sprintf("%. 10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

//skip one char at the beginning
//will also skip a space if there is one present :D
func (l *lexer) BSkipOne() {
	l.start++
	l.pos = l.start
	resv := l.peek()
	if resv == ' ' {
		l.start++
	}
	l.pos++
}

//lex a item with a left and a right delimiter
//takes a right delimiter
func (l *lexer) LexLRDelim(rd rune) {
	//skip left delim
	l.BSkipOne()
	var rp rune
	for {
		_ = l.next()
		rp = l.peek()
		if rp == rd {
			for i := 0; i < len(sym); i++ {
				if rd == sym[i] {
					l.emit(items[i])
					break
				}
			}
			l.start++
			l.pos = l.start
			break
		}
	}
}

func (l *lexer) emit(t itemType) {
	if t == 7 {
	}
	l.items <- item{t, l.input[l.start:l.pos], l.id}
	l.id++
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() rune {
	var r rune
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item

}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		state: lexText,
		items: make(chan item, 2),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}
