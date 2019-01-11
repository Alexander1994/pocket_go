package main

import "time"

type PrimFunc = func(args []*Object, env *Env) *Object

var Functs = map[string]PrimFunc{
	// arithmetic funcs
	"+": add,
	"-": minus,
	"/": divide,
	"*": multi,

	// variable funcs / mutates environments
	"def":  def,
	"defn": defn,
	"set":  set,

	// logic funcs
	"for": forloop,
	"if":  ifcond,
	"=":   equalval,
	"eq":  equalref,
	">":   cmp,

	// goruotine funcs
	"go": goroutine,
	"<-": channelop,

	// go flavored lisps
	"quote": quote,
	"[]":    subscript,
	"[:]":   sublist,

	// misc. funcs
	"sleep":   sleep,
	"println": printn,

	// macro funcs
	"macro": macro,
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

func printn(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	for i, arg := range args {
		arg.print()
		if i != len(args)-1 {
			print(" ")
		}
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

// (macro $symbol ($args...) $expr...)
func macro(args []*Object, env *Env) *Object {
	arglist := args[1].List()
	tempateargs := make([][]*Object, len(arglist))
	arrayindex := make(map[string]int)
	exprs := args[2:]
	for i, arg := range arglist {
		arrayindex[arg.Symbol()] = i
		tempateargs[i] = make([]*Object, 0)
	}
	for _, expr := range exprs {
		cacheMacro(expr, arrayindex, tempateargs)
	}
	env.Add(Car(args).Symbol(), Macro(tempateargs, exprs))
	return nilObj
}

func cacheMacro(obj *Object, arrayindex map[string]int, tempateargs [][]*Object) {
	if obj.Type() == symbolT {
		if ind, ok := arrayindex[obj.Symbol()]; ok {
			tempateargs[ind] = append(tempateargs[ind], obj)
		}
	} else if obj.Type() == cellT {
		for _, objIt := range obj.List() {
			cacheMacro(objIt, arrayindex, tempateargs)
		}
	}
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

/*
(macro plus1 (x) (set x (+ x 1)) )
(def y 0)
(plus1 y)
*/
// ($symbol ?$expr...)
func (o *Object) RunMacro(args []*Object, env *Env) (result *Object) {
	// setup
	macro := o.Macro()
	for i, arg := range args {
		for _, templ := range macro.templateargs[i] {
			(*templ) = *arg
		}
	}
	// run
	return EvalList(macro.expr, env)[len(macro.expr)-1]
}

// func expand(obj *Object, env *Env) {
// 	if obj.Type() == symbolT {
// 		evalobj, _ := env.find(obj.Symbol())
// 		if evalobj != nilObj {
// 			*obj = *Eval(evalobj, env)
// 		}
// 	} else if obj.Type() == cellT {
// 		for _, lobj := range obj.List() {
// 			expand(lobj, env)
// 		}
// 	}
// }

// (go $symbol ?$expr...)
func goroutine(args []*Object, env *Env) *Object {
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
func channelop(args []*Object, env *Env) *Object {
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

// '$expr
func quote(args []*Object, env *Env) *Object {
	if len(args) != 1 {
		panic("invalid arg count passed to quote")
	}
	return Car(args)
}

// ([] $expr $expr) $1 evals to num, $2 evals to list
func subscript(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid arg count passed to subscript")
	}
	numObj := Eval(args[0], env)
	listObj := Eval(args[1], env)
	if numObj.Type() != numT || listObj.Type() != cellT {
		panic("invalid types passed to subscript op")
	}
	return listObj.List()[uint(numObj.Num())]
}

// ([:] $expr $expr $expr) $1 evals to num, $2 evals to num, $3 evals to list
func sublist(args []*Object, env *Env) *Object {
	if len(args) != 3 {
		panic("invalid arg count passed to sublist")
	}
	upperindex := Eval(args[0], env)
	lowerindex := Eval(args[1], env)
	listObj := Eval(args[2], env)
	return List(listObj.List()[uint(lowerindex.Num()):uint(upperindex.Num())])
}
