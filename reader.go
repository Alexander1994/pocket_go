package main

import (
	"io/ioutil"
)

type Reader struct {
	index  int
	buffer []rune
	size   int
}

var reader Reader

func (r *Reader) NextRune() (rune, bool) {
	curr := r.index
	r.index++
	if curr < r.size {
		return r.buffer[curr], false
	}
	return ' ', true
}

func (r *Reader) Peek() rune {
	if r.index < r.size { // i is always pointing at the next rune
		return r.buffer[r.index]
	}
	panic("sudden end of file. while peeking")
}

func (r *Reader) SubBuffer(start, end int) string {
	return string(r.buffer[start : start+end])
}

func (r *Reader) Load(fname string) {
	srcBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	r.buffer = []rune(string(srcBytes))
	r.size = len(r.buffer)
}
