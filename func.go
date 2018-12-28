package main

import "time"

type PrimFunc = func(args *[]Object, env *Env) *Object

var Functs = map[string]PrimFunc{
	"+":       add,
	"-":       minus,
	"/":       divide,
	"*":       multi,
	"println": printLn,
	"def":     def,
	"defn":    defn,
	"go":      goRoutine,
	"sleep":   sleep,
}

func isStartOfFunc(r rune) bool {
	for name := range Functs {
		if rune(name[0]) == r {
			return true
		}
	}
	return false
}

func add(args *[]Object, env *Env) *Object {
	args = EvalList(args, env)
	sum := float32(0)
	for _, arg := range *args {
		sum = sum + arg.Num()
	}
	return Num(sum)
}
func minus(args *[]Object, env *Env) *Object {
	args = EvalList(args, env)
	diff := Car(args).Num()

	for _, arg := range *Cdr(args) {
		diff = diff - arg.Num()
	}
	return Num(diff)
}
func divide(args *[]Object, env *Env) *Object {
	args = EvalList(args, env)
	num := Car(args).Num()
	for _, arg := range *Cdr(args) {
		num = num / arg.Num()
	}
	return Num(num)
}
func multi(args *[]Object, env *Env) *Object {
	args = EvalList(args, env)
	sum := Car(args).Num()
	for _, arg := range *Cdr(args) {
		sum = sum * arg.Num()
	}
	return Num(sum)
}
func printLn(args *[]Object, env *Env) *Object {
	args = EvalList(args, env)
	for _, arg := range *args {
		arg.print()
		print(" ")
	}
	println()
	return nilObj
}

// (sleep $num)
func sleep(args *[]Object, env *Env) *Object {
	if len(*args) != 1 {
		panic("sleep gets 1 arg which is a num")
	}
	num := Eval(Car(args), env)
	if num.Type() != numT {
		panic("sleep gets 1 arg which is a num")
	}
	length := time.Duration(num.Num()) * time.Millisecond
	time.Sleep(length)
	return nilObj
}

// (def $symbol $expr)
func def(args *[]Object, env *Env) *Object {
	if len((*args)) != 2 {
		panic("invalid args length passed to def")
	}
	expr := Eval(&(*args)[1], env)
	env.Add((*args)[0].Symbol(), expr)
	return nilObj
}

// (defn $symbol ($symbol...) $expr...)
func defn(args *[]Object, env *Env) *Object {
	if len((*args)) < 3 {
		panic("defn must have atleast 3 args: defn $symbol ($symbol...) $expr...")
	}
	expr := (*args)[2:]
	symbol := (*args)[0].Symbol()
	env.Add(symbol, Function(symbol, &(*args)[1], &expr))
	return nilObj
}

// (go $symbol $expr...)
func goRoutine(args *[]Object, env *Env) *Object {
	if len((*args)) < 2 {
		panic("go primitive requires a function and its args")
	}
	function := Eval(&(*args)[0], env)
	if function.Type() != funcT {
		panic("go primitive requires a function and its args")
	}
	go function.CallFunc(Cdr(args), env)
	return nilObj
}

// ($symbol $expr...)
func (o *Object) CallFunc(args *[]Object, env *Env) (returnVal *Object) {
	newEnv := o.pushFuncEnv(args, env)
	resultList := *EvalList(o.Function().expr, newEnv)
	env.popFuncEnv()
	return &resultList[len(resultList)-1]
}
