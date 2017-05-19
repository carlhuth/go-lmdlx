package mdlx

import (
	"unicode"
	"unicode/utf8"
)

func lexNewLine(l *lexer) stateFn {
	l.pos++
	l.emit(itemNewline)
	return lexText
}

func lexUl(l *lexer) stateFn {
	l.pos++
	for {
		current, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if current != '\n' {
			l.pos++
		} else {
			break
		}
	}
	l.emit(itemUl)
	return lexText
}

func lexTabUl(l *lexer) stateFn {
	l.pos++
	for {
		current, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if current != '\n' {
			l.pos++
		} else {
			break
		}
	}
	l.emit(itemTabUl)
	return lexText
}

func lexHeader(l *lexer) stateFn {
	var header int
	l.pos++
	for {
		current, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if current == '#' {
			l.pos++
			header++
		} else {
			break
		}
	}
	for {
		current, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if unicode.IsSymbol(current) || current == '\n' {
			switch header {
			case 0:
				l.emit(itemHeader)
			case 1:
				l.emit(itemHeader2)
			case 2:
				l.emit(itemHeader3)
			case 3:
				l.emit(itemHeader4)
			}

			return lexText
		}
		if l.next() == eof {
			l.emit(itemHeader)
			return lexText
		}
	}
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

//should lex all text, then decide what function is next :)
func lexText(l *lexer) stateFn {
	var last rune = '\n'
	for {
		current, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		switch {
		case current == '#':
			if last == '\n' {
				if l.pos > l.start {
					l.emit(itemText)
				}
				last = '#'
				return lexHeader
			} else {
				l.pos++
				last = 't'
			}
		case current == '\n':
			if l.pos > l.start {
				l.emit(itemText)
			}
			last = '\n'
			l.pos++
			l.emit(itemNewline)
		case current == '*':
			if l.pos > l.start {
				l.emit(itemText)
			}
			last = '*'
			return lexAsterisk
		case current == '`':
			if l.pos > l.start {
				l.emit(itemText)
			}
			last = '`'
			return lexBacktick
		case current == '_':
			if l.pos > l.start {
				l.emit(itemText)
			}
			last = '_'
			return lexUnderscore
		case current == '+':
			loclast, _ := utf8.DecodeRuneInString(l.input[l.pos-1:])
			if loclast != '\n' {
				last = '+'
				return lexTabUl
			}
			if l.pos > l.start {
				l.emit(itemText)
			}
			last = '+'
			return lexUl
		case l.next() == eof:
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.emit(itemEOF)
			return nil
		}
		last = current
	}
	return nil
}
