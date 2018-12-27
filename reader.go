package main

import (
	"io/ioutil"
)

var i = 0 // file index

var buffer []rune
var size int

func NextRune() (rune, bool) {
	curr := i
	i++
	if curr < size {
		return buffer[curr], false
	}
	return ' ', true
}

func Peek() rune {
	if i < size { // i is always pointing at the next rune
		return buffer[i]
	}
	panic("sudden end of file???? while peeking :0")
}

func SubBuffer(start, end int) string {
	return string(buffer[start : start+end])
}

func Load(fname string) {
	srcBytes, err := ioutil.ReadFile(fname)
	check(err)
	buffer = []rune(string(srcBytes))
	size = len(buffer)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
