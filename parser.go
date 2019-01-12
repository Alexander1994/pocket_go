package main

import (
	"strconv"
	"unicode"
)

func ParseExpr() (obj *Object) {
	for {
		c, isEOF := reader.NextRune()
		if isEOF {
			return nilObj
		}
		if unicode.IsSpace(c) {
			continue
		}
		if c == '\'' {
			return ParseQuote(c)
		}
		if c == ';' && reader.Peek() == ';' {
			for ; c != '\n' && !isEOF; c, isEOF = reader.NextRune() {
			}
		}
		if unicode.IsLetter(c) || IsStartOfFunc(c) {
			return ParsePhrase(c)
		}
		if unicode.IsDigit(c) || (c == '-' && unicode.IsDigit(reader.Peek())) {
			return ParseNum(c)
		}
		if c == '(' {
			return ParseList()
		}
		if c == ')' {
			return closeParenObj
		}
	}
}

func ParseQuote(r rune) *Object {
	list := []*Object{Primitve("quote"), ParseExpr()}
	return List(list)
}

func ParseNum(r rune) *Object {
	atom := ParseAtom(r)
	f, err := strconv.ParseFloat(atom, 32)
	if err != nil {
		panic("invalid num in parsing")
	}
	return Num(float32(f))
}

func ParsePhrase(r rune) *Object {
	atom := ParseAtom(r)
	if atom == "chan" {
		return Channel()
	}
	if _, ok := Functs[atom]; ok {
		return Primitve(atom)
	}
	return Symbol(atom)
}

func ParseList() *Object {
	evalList := make([]*Object, 0)
	for {
		obj := ParseExpr()
		if obj == closeParenObj {
			return List(evalList)
		}
		evalList = append(evalList, obj)
	}
}

func ParseAtom(r rune) string {
	start := reader.index
	end := 0
	for {
		c := reader.Peek()
		if unicode.IsSpace(c) || c == '(' || c == ')' {
			return string(r) + reader.SubBuffer(start, end)
		}
		end++
		reader.NextRune()
	}
}
