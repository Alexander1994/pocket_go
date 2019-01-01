package main

import (
	"strconv"
	"unicode"
)

func parseExpr() (obj *Object) {
	for {
		c, isEOF := NextRune()
		if isEOF {
			return nilObj
		}
		if unicode.IsSpace(c) {
			continue
		}
		if c == '\'' {
			return parseQuote(c)
		}
		if c == ';' && Peek() == ';' {
			for ; c != '\n' && !isEOF; c, isEOF = NextRune() {
			}
		}
		if unicode.IsLetter(c) || isStartOfFunc(c) {
			return parsePhrase(c)
		}
		if unicode.IsDigit(c) || (c == '-' && unicode.IsDigit(Peek())) {
			return parseNum(c)
		}
		if c == '(' {
			return parseList()
		}
		if c == ')' {
			return closeParenObj
		}
	}
}

func parseQuote(r rune) *Object {
	list := []*Object{Primitve("quote"), parseExpr()}
	return List(list)
}

func parseNum(r rune) *Object {
	atom := parseAtom(r)
	f, err := strconv.ParseFloat(atom, 32)
	if err != nil {
		panic("invalid num in parsing")
	}
	return Num(float32(f))
}

func parsePhrase(r rune) *Object {
	atom := parseAtom(r)
	if atom == "chan" {
		return Channel()
	}
	if _, ok := Functs[atom]; ok {
		return Primitve(atom)
	}
	return Symbol(atom)
}

func parseList() *Object {
	evalList := make([]*Object, 0)
	for {
		obj := parseExpr()
		if obj == closeParenObj {
			return List(evalList)
		}
		evalList = append(evalList, obj)
	}
}

func parseAtom(r rune) string {
	start := i
	end := 0
	for {
		c := Peek()
		if unicode.IsSpace(c) || c == '(' || c == ')' {
			return string(r) + SubBuffer(start, end)
		}
		end++
		NextRune()
	}
}
