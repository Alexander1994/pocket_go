package main

import "time"

type PrimFunc = func(args []*Object, env *Env) *Object

var Functs = map[string]PrimFunc{
	// arithmetic funcs
	"+": add,
	"-": minus,
	"/": divide,
	"*": multi,

	// variable funcs
	"def":  def,
	"defn": defn,
	"set":  set,

	// goruotine funcs
	"go": goRoutine,
	"<-": channelOp,

	// logic funcs
	"for": forloop,
	"if":  ifcond,
	"=":   equalval,
	"eq":  equalref,
	">":   cmp,

	// misc. funcs
	"sleep":   sleep,
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

func add(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	sum := float32(0)
	for _, arg := range args {
		sum = sum + arg.Num()
	}
	return Num(sum)
}
func minus(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	diff := Car(args).Num()

	for _, arg := range Cdr(args) {
		diff = diff - arg.Num()
	}
	return Num(diff)
}
func divide(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	num := Car(args).Num()
	for _, arg := range Cdr(args) {
		num = num / arg.Num()
	}
	return Num(num)
}
func multi(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	sum := Car(args).Num()
	for _, arg := range Cdr(args) {
		sum = sum * arg.Num()
	}
	return Num(sum)
}
func printLn(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	for _, arg := range args {
		arg.print()
		print(" ")
	}
	println()
	return nilObj
}

// (sleep $num)
func sleep(args []*Object, env *Env) *Object {
	if len(args) != 1 {
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
func def(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid args length passed to def")
	}
	expr := args[1]
	expr = Eval(expr, env)
	env.Add(Car(args).Symbol(), expr)
	return nilObj
}

// (set $symbol $expr)
func set(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid arg count in pass to set")
	}
	expr := Eval(args[1], env)
	env.Set(Car(args).Symbol(), expr)
	return nilObj
}

// (defn ?$symbol ($symbol...) $expr...)
func defn(args []*Object, env *Env) *Object {
	if len(args) == 1 {
		panic("defn must have atleast 2 or more args: (defn ?$symbol ($symbol...) $expr...)")
	}
	var closure *Env
	if env.isTempEnv() {
		closure = env
	}
	if len(args) == 2 {
		return Function(Car(args), closure, Cdr(args))
	}
	env.Add(Car(args).Symbol(), Function(args[1], closure, args[2:]))
	return nilObj
}

// ($symbol ?$expr...)
func (o *Object) CallFunc(args []*Object, env *Env) (returnVal *Object) {
	function := o.Function()
	currEnv := env
	if function.closure != nil {
		currEnv = function.closure
	}
	newEnv := o.pushFuncEnv(args, currEnv)
	resultList := EvalList(function.expr, newEnv)
	currEnv.popFuncEnv()
	return resultList[len(resultList)-1]
}

// (go $symbol ?$expr...)
func goRoutine(args []*Object, env *Env) *Object {
	if len(args) < 1 {
		panic("go primitive requires a function and its args")
	}
	function := Eval(args[0], env)
	if function.Type() != funcT {
		panic("go primitive requires a function and its args")
	}
	go function.CallFunc(Cdr(args), env)
	return nilObj
}

// send: (<- $channel $expr) OR recv: (<- $channel)
func channelOp(args []*Object, env *Env) *Object {
	if len(args) == 2 { // send
		Eval(Car(args), env).Send(Eval(args[1], env))
		return nilObj
	} else if len(args) == 1 { // recv
		return Eval(Car(args), env).Recv()
	} else {
		panic("invalid call to channel op")
	}
}

// (for $expr ?$expr...)
func forloop(args []*Object, env *Env) *Object {
	num := Eval(Car(args), env)
	if len(args) <= 1 || num.Type() != numT {
		panic("for loop must have a num in the first args")
	}
	for ; num.Type() == numT && num.Num() != 0; num = Eval(Car(args), env) {
		EvalList(Cdr(args), env)
	}
	return nilObj
}

// (if $expr ?$expr...)
func ifcond(args []*Object, env *Env) *Object {
	num := Eval(Car(args), env)
	if len(args) <= 1 || num.Type() != numT {
		panic("for loop must have a num in the first args")
	}
	if num.Num() != 0 {
		EvalList(Cdr(args), env)
	}
	return nilObj
}

// (= expr...)
func eq(args []*Object, env *Env) *Object {
	if len(args) == 0 {
		panic("must have values/exprs in call to '=' function")
	}
	evalargs := EvalList(args, env)
	car := Car(evalargs)
	if car.Type() != numT {
		return Num(0)
	}
	num := car.Num()
	for i := 1; i < len(evalargs); i++ {
		if evalargs[i].Type() != numT || evalargs[i].Num() != num {
			return Num(0)
		}
	}
	return Num(1)
}

// (= expr...)
func equalval(args []*Object, env *Env) *Object {
	if len(args) == 0 {
		panic("must have values/exprs in call to '=' function")
	}
	evalargs := EvalList(args, env)
	car := Car(evalargs)
	if car.Type() != numT {
		return Num(0)
	}
	num := car.Num()
	for i := 1; i < len(evalargs); i++ {
		if evalargs[i].Type() != numT || evalargs[i].Num() != num {
			return Num(0)
		}
	}
	return Num(1)
}

// (eq $expr...)
func equalref(args []*Object, env *Env) *Object {
	if len(args) == 0 {
		panic("must have values/exprs in call to '=' function")
	}
	evalargs := EvalList(args, env)
	car := Car(evalargs)
	for i := 1; i < len(evalargs); i++ {
		if evalargs[i] != car {
			return Num(0)
		}
	}
	return Num(1)
}

// (> $expr $expr)
func cmp(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid args count passed to cmp")
	}
	evalargs := EvalList(args, env)
	if evalargs[0].Type() != numT || evalargs[1].Type() != numT {
		return Num(0)
	}
	if evalargs[0].Num() > evalargs[1].Num() {
		return Num(1)
	}
	return Num(0)
}
