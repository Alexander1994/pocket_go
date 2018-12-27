package main

type Method = func(args []Object) Object

var Functs = map[string]Method{
	"+":       add,
	"-":       minus,
	"/":       divide,
	"*":       multi,
	"println": printLn,
}

func isStartOfFunc(r rune) bool {
	for name := range Functs {
		if rune(name[0]) == r {
			return true
		}
	}
	return false
}

func add(args []Object) Object {
	sum := float32(0)
	for _, arg := range args {
		sum = sum + arg.Num()
	}
	return Num(sum)
}
func minus(args []Object) Object {
	diff := args[0].Num()

	for _, arg := range args[1:] {
		diff = diff - arg.Num()
	}
	return Num(diff)
}
func divide(args []Object) Object {
	diff := args[0].Num()
	for _, arg := range args[1:] {
		diff = diff / arg.Num()
	}
	return Num(diff)
}
func multi(args []Object) Object {
	sum := float32(0)
	for _, arg := range args {
		sum = sum * arg.Num()
	}
	return Num(sum)
}
func printLn(args []Object) Object {
	for _, arg := range args {
		arg.print()
	}
	println()
	return nilObj
}
