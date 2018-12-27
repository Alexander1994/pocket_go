package main

type Method = func(args []Object, env *Env) Object

var Functs = map[string]Method{
	"+":       add,
	"-":       minus,
	"/":       divide,
	"*":       multi,
	"println": printLn,
	"def":     def,
	"defn":    defn,
}

func isStartOfFunc(r rune) bool {
	for name := range Functs {
		if rune(name[0]) == r {
			return true
		}
	}
	return false
}

func add(args []Object, env *Env) Object {
	args = Eval_List(args, env)
	sum := float32(0)
	for _, arg := range args {
		sum = sum + arg.Num()
	}
	return Num(sum)
}
func minus(args []Object, env *Env) Object {
	args = Eval_List(args, env)
	diff := args[0].Num()

	for _, arg := range args[1:] {
		diff = diff - arg.Num()
	}
	return Num(diff)
}
func divide(args []Object, env *Env) Object {
	args = Eval_List(args, env)
	num := args[0].Num()
	for _, arg := range args[1:] {
		num = num / arg.Num()
	}
	return Num(num)
}
func multi(args []Object, env *Env) Object {
	args = Eval_List(args, env)
	sum := args[0].Num()
	for _, arg := range args[1:] {
		sum = sum * arg.Num()
	}
	return Num(sum)
}
func printLn(args []Object, env *Env) Object {
	args = Eval_List(args, env)
	for _, arg := range args {
		arg.print()
	}
	println()
	return nilObj
}

// (def $symbol $expr)
func def(args []Object, env *Env) Object {
	if len(args) != 2 {
		panic("invalid args length passed to def")
	}
	expr := Eval(args[1], env)
	env.Add(args[0].Symbol(), &expr)
	return nilObj
}

// (defn $symbol ($symbol...) $expr...)
func defn(args []Object, env *Env) Object {
	return nilObj
}
