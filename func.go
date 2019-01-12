package main

import "time"

type PrimFunc = func(args []*Object, env *Env) *Object

var Functs = map[string]PrimFunc{
	// arithmetic funcs
	"+": Add,
	"-": Minus,
	"/": Divide,
	"*": Multi,

	// variable funcs / mutates environments
	"def":  Def,
	"defn": Defn,
	"set":  Set,

	// logic funcs
	"for": ForLoop,
	"if":  IfCond,
	"=":   EqualVal,
	"eq":  EqualRef,
	">":   Cmp,

	// goruotine funcs
	"go": GoRoutine,
	"<-": ChannelOp,

	// go flavored lisps
	"quote": Quote,
	"[]":    Subscript,
	"[:]":   Sublist,

	// misc. funcs
	"sleep":   Sleep,
	"println": Printn,

	// macro funcs
	"macro": DefMacro,
}

func IsStartOfFunc(r rune) bool {
	for name := range Functs {
		if rune(name[0]) == r {
			return true
		}
	}
	return false
}

func Add(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	sum := float32(0)
	for _, arg := range args {
		sum = sum + arg.Num()
	}
	return Num(sum)
}
func Minus(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	diff := Car(args).Num()

	for _, arg := range Cdr(args) {
		diff = diff - arg.Num()
	}
	return Num(diff)
}
func Divide(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	num := Car(args).Num()
	for _, arg := range Cdr(args) {
		num = num / arg.Num()
	}
	return Num(num)
}
func Multi(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	sum := Car(args).Num()
	for _, arg := range Cdr(args) {
		sum = sum * arg.Num()
	}
	return Num(sum)
}

func Printn(args []*Object, env *Env) *Object {
	args = EvalList(args, env)
	for i, arg := range args {
		arg.Print()
		if i != len(args)-1 {
			print(" ")
		}
	}
	println()
	return nilObj
}

// (sleep $num)
func Sleep(args []*Object, env *Env) *Object {
	if len(args) != 1 {
		panic("sleep gets 1 arg which is a num")
	}
	num := Eval(Car(args), env)
	if num.Type() != NumT {
		panic("sleep gets 1 arg which is a num")
	}
	length := time.Duration(num.Num()) * time.Millisecond
	time.Sleep(length)
	return nilObj
}

// (def $symbol $expr)
func Def(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid args length passed to def")
	}
	expr := args[1]
	expr = Eval(expr, env)
	env.Add(Car(args).Symbol(), expr)
	return nilObj
}

// (set $symbol $expr)
func Set(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid arg count in pass to set")
	}
	expr := Eval(args[1], env)
	env.Set(Car(args).Symbol(), expr)
	return nilObj
}

// (defn ?$symbol ($symbol...) $expr...)
func Defn(args []*Object, env *Env) *Object {
	if len(args) == 1 {
		panic("defn must have atleast 2 or more args: (defn ?$symbol ($symbol...) $expr...)")
	}

	var closure *Env
	if env.IsTempEnv() {
		closure = env
	}
	if len(args) == 2 {
		return Function(Car(args), closure, Cdr(args))
	}
	env.Add(Car(args).Symbol(), Function(args[1], closure, args[2:]))
	return nilObj
}

// (macro $symbol ($args...) $expr...)
func DefMacro(args []*Object, env *Env) *Object {
	arglist := args[1].List()
	tempateargs := make([][]*Object, len(arglist))
	arrayindex := make(map[string]int)
	exprs := args[2:]
	for i, arg := range arglist {
		arrayindex[arg.Symbol()] = i
		tempateargs[i] = make([]*Object, 0)
	}
	for _, expr := range exprs {
		CacheMacro(expr, arrayindex, tempateargs)
	}
	env.Add(Car(args).Symbol(), Macro(tempateargs, exprs))
	return nilObj
}

func CacheMacro(obj *Object, arrayindex map[string]int, tempateargs [][]*Object) {
	if obj.Type() == SymbolT {
		if ind, ok := arrayindex[obj.Symbol()]; ok {
			tempateargs[ind] = append(tempateargs[ind], obj)
		}
	} else if obj.Type() == CellT {
		for _, objIt := range obj.List() {
			CacheMacro(objIt, arrayindex, tempateargs)
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
	newEnv := o.PushFuncEnv(args, currEnv)
	resultList := EvalList(function.expr, newEnv)
	currEnv.PopFuncEnv()
	return resultList[len(resultList)-1]
}

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

// (go $symbol ?$expr...)
func GoRoutine(args []*Object, env *Env) *Object {
	if len(args) < 1 {
		panic("go primitive requires a function and its args")
	}
	function := Eval(args[0], env)
	if function.Type() != FuncT {
		panic("go primitive requires a function and its args")
	}
	go function.CallFunc(Cdr(args), env)
	return nilObj
}

// send: (<- $channel $expr) OR recv: (<- $channel)
func ChannelOp(args []*Object, env *Env) *Object {
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
func ForLoop(args []*Object, env *Env) *Object {
	num := Eval(Car(args), env)
	if len(args) <= 1 || num.Type() != NumT {
		panic("for loop must have a num in the first args")
	}
	for ; num.Type() == NumT && num.Num() != 0; num = Eval(Car(args), env) {
		EvalList(Cdr(args), env)
	}
	return nilObj
}

// (if $expr ?$expr...)
func IfCond(args []*Object, env *Env) *Object {
	num := Eval(Car(args), env)
	if len(args) <= 1 || num.Type() != NumT {
		panic("for loop must have a num in the first args")
	}
	if num.Num() != 0 {
		EvalList(Cdr(args), env)
	}
	return nilObj
}

// (= expr...)
func EqualVal(args []*Object, env *Env) *Object {
	if len(args) == 0 {
		panic("must have values/exprs in call to '=' function")
	}
	evalargs := EvalList(args, env)
	car := Car(evalargs)
	if car.Type() != NumT {
		return Num(0)
	}
	num := car.Num()
	for i := 1; i < len(evalargs); i++ {
		if evalargs[i].Type() != NumT || evalargs[i].Num() != num {
			return Num(0)
		}
	}
	return Num(1)
}

// (eq $expr...)
func EqualRef(args []*Object, env *Env) *Object {
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
func Cmp(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid args count passed to cmp")
	}
	evalargs := EvalList(args, env)
	if evalargs[0].Type() != NumT || evalargs[1].Type() != NumT {
		return Num(0)
	}
	if evalargs[0].Num() > evalargs[1].Num() {
		return Num(1)
	}
	return Num(0)
}

// '$expr
func Quote(args []*Object, env *Env) *Object {
	if len(args) != 1 {
		panic("invalid arg count passed to quote")
	}
	return Car(args)
}

// ([] $expr $expr) $1 evals to num, $2 evals to list
func Subscript(args []*Object, env *Env) *Object {
	if len(args) != 2 {
		panic("invalid arg count passed to subscript")
	}
	numObj := Eval(args[0], env)
	listObj := Eval(args[1], env)
	if numObj.Type() != NumT || listObj.Type() != CellT {
		panic("invalid types passed to subscript op")
	}
	return listObj.List()[uint(numObj.Num())]
}

// ([:] $expr $expr $expr) $1 evals to num, $2 evals to num, $3 evals to list
func Sublist(args []*Object, env *Env) *Object {
	if len(args) != 3 {
		panic("invalid arg count passed to sublist")
	}
	upperindex := Eval(args[0], env)
	lowerindex := Eval(args[1], env)
	listObj := Eval(args[2], env)
	return List(listObj.List()[uint(lowerindex.Num()):uint(upperindex.Num())])
}
