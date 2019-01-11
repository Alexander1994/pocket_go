package main

func Eval(o *Object, env *Env) *Object {
	switch o.objT {
	case numT, primitveT, nilT, funcT, chanT:
		return o
	case symbolT:
		obj, _ := env.find(o.Symbol())
		if obj == nilObj {
			panic("undefined symbol " + o.Symbol())
		}
		return obj
	case cellT:
		function := Eval(o.Car(), env)
		args := o.Cdr()
		return call(function, args, env)
	}
	return nilObj
}

func EvalList(list []*Object, env *Env) []*Object {
	evalList := make([]*Object, len(list))
	for i, item := range list {
		obj := Eval(item, env)
		evalList[i] = obj
	}
	return evalList
}

func call(function *Object, args []*Object, env *Env) *Object {
	switch function.Type() {
	case primitveT:
		return function.CallPrim(args, env)
	case funcT:
		return function.CallFunc(args, env)
	case macroT:
		return function.RunMacro(args, env)
	}
	panic("Head of cell/list is not a function: ")
}
