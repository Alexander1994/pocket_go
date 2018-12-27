package main

type PrimFunc = func(args []Object, env *Env) Object

var Functs = map[string]PrimFunc{
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
	args = EvalList(args, env)
	sum := float32(0)
	for _, arg := range args {
		sum = sum + arg.Num()
	}
	return Num(sum)
}
func minus(args []Object, env *Env) Object {
	args = EvalList(args, env)
	diff := args[0].Num()

	for _, arg := range args[1:] {
		diff = diff - arg.Num()
	}
	return Num(diff)
}
func divide(args []Object, env *Env) Object {
	args = EvalList(args, env)
	num := args[0].Num()
	for _, arg := range args[1:] {
		num = num / arg.Num()
	}
	return Num(num)
}
func multi(args []Object, env *Env) Object {
	args = EvalList(args, env)
	sum := args[0].Num()
	for _, arg := range args[1:] {
		sum = sum * arg.Num()
	}
	return Num(sum)
}
func printLn(args []Object, env *Env) Object {
	args = EvalList(args, env)
	for _, arg := range args {
		arg.print()
		print(" ")
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
	if len(args) < 3 {
		panic("defn must have atleast 3 args: defn $symbol ($symbol...) $expr...")
	}
	expr := args[2:]
	symbol := args[0].Symbol()
	fun := Function(symbol, &args[1], &expr)
	env.Add(symbol, &fun)
	return nilObj
}

// ($symbol $expr...)
func (o *Object) CallFunc(args []Object, env *Env) (returnVal Object) {
	funDef := o.Function()
	defArgs := funDef.args.List()
	if len(defArgs) != len(args) {
		panic("args in call to function != function args")
	}
	evalArgs := EvalList(args, env)
	newEnv := AddAndGetNewEnv(env)

	for i, arg := range evalArgs {
		newEnv.Add(defArgs[i].Symbol(), &arg)
	}
	exprList := *funDef.expr
	resultList := EvalList(exprList, newEnv)
	env.popEnvStack()
	return resultList[len(resultList)-1]
}
