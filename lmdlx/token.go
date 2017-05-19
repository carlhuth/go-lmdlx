package lmdlx

import (
	"fmt"
	"unicode/utf8"
)

func lexNewLine(l *lexer) stateFn {
	l.pos++
	l.emit(itemNewline)
	return lexText
}

func lexAsterisk(l *lexer) stateFn {
	l.pos++
	for {
		if r := l.next(); r == '\n' || r == '*' {
			l.emit(itemAsterisk)
			break
		}
	}
	return lexText
}

func lexUnderscore(l *lexer) stateFn {
	l.pos++
	for {
		if r := l.next(); r == '\n' || r == '_' {
			l.emit(itemAsterisk)
			break
		}
	}
	return lexText
}

func lexBacktick(l *lexer) stateFn {
	l.pos++
	for {
		if r := l.next(); r == '\n' || r == '`' {
			l.emit(itemBacktick)
			break
		}
	}
	return lexText
}

func lexBlockq(l *lexer) stateFn {
	l.start++
	l.pos = l.start
	resv := l.peek()
	fmt.Println("RESV", string(resv))
	if resv == ' ' {
		fmt.Println("NOTIFICATION SQUAD")
		l.start++
	}
	l.pos++
	for {
		if r := l.next(); r == '\n' || r == '`' {
			l.emit(itemBlockq)
			break
		}
	}
	return lexText
}

//should lex all text, then decide what function is next :)
func lexText(l *lexer) stateFn {
	for {
		current, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		switch {
		case current == '\n':
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.pos++
			l.emit(itemNewline)
		case current == '*':
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexAsterisk
		case current == '`':
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexBacktick
		case current == '_':
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexUnderscore
		case current == '>':
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexBlockq
		case l.next() == eof:
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.emit(itemEOF)
			return nil
		}
	}
	return nil
}
