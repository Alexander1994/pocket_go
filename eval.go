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
		if function.Type() != primitveT && function.Type() != funcT {
			panic("Head of cell/list is not a function: ")
		}
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
	if function.Type() == primitveT {
		return function.CallPrim(args, env)
	}
	if function.Type() == funcT {
		return function.CallFunc(args, env)
	}
	panic("invalid call to func")
}
