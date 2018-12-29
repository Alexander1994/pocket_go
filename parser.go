package main

import (
	"strconv"
	"unicode"
)

func readExpr() (obj *Object) {
	for {
		c, isEOF := NextRune()
		if isEOF {
			return nilObj
		}
		if unicode.IsSpace(c) {
			continue
		}
		if unicode.IsLetter(c) || isStartOfFunc(c) {
			return readPhrase(c)
		}
		if unicode.IsDigit(c) || (c == '-' && unicode.IsDigit(Peek())) {
			return readNum(c)
		}
		if c == '(' {
			return readList()
		}
		if c == ')' {
			return closeParenObj
		}
	}
}

func readNum(r rune) *Object {
	atom := readAtom(r)
	f, err := strconv.ParseFloat(atom, 32)
	if err != nil {
		panic("invalid num in parsing")
	}
	return Num(float32(f))
}

func readPhrase(r rune) *Object {
	atom := readAtom(r)
	if atom == "chan" {
		return Channel()
	}
	if _, ok := Functs[atom]; ok {
		return Primitve(atom)
	}
	return Symbol(atom)
}

func readList() *Object {
	evalList := make([]*Object, 0)
	for {
		obj := readExpr()
		if obj == closeParenObj {
			return List(evalList)
		}
		evalList = append(evalList, obj)
	}
}

func readAtom(r rune) string {
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
