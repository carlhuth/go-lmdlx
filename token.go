package lmdlx

import (
	"unicode/utf8"
)

func lexNewLine(l *lexer) stateFn {
	l.pos++
	l.emit(itemNewline)
	return lexText
}

/*
func lexAsterisk(l *lexer) stateFn {
	l.BSkipOne()
	var rp rune
	for {
		_ = l.next()
		rp = l.peek()
		if rp == '*' {
			l.emit(itemAsterisk)
			l.start++
			l.pos = l.start
			break
		}
	}

	return lexText
}
*/
func lexAsterisk(l *lexer) stateFn {
	l.LexLRDelim('*')
	return lexText
}
func lexUnderscore(l *lexer) stateFn {
	l.LexLRDelim('_')
	return lexText
}

func lexBacktick(l *lexer) stateFn {
	l.LexLRDelim('`')
	return lexText
}

func lexBlockq(l *lexer) stateFn {
	l.BSkipOne()
	var rp rune
	for {
		_ = l.next()
		if rp = l.peek(); rp == '\n' {
			l.emit(itemBlockq)
			break
		}
	}
	return lexText
}

func lexTilda(l *lexer) stateFn {
	l.LexLRDelim('~')
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
		case current == '~':
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexTilda
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
