package main

func Eval(o Object, env *Env) Object {
	switch o.objT {
	case numT, primitveT, nilT:
		return o
	case symbolT:
		obj := env.find(o.Symbol())
		if obj == &nilObj {
			panic("undefined symbol " + o.Symbol())
		}
		return *obj
	case cellT:
		function := Eval(o.Car(), env)
		args := o.Cdr()
		if function.Type() != primitveT {
			panic("Head of cell/list is not a function")
		}
		return call(function, args, env)
	}
	o.print()
	return nilObj
}

func Eval_List(list []Object, env *Env) []Object {
	evalList := make([]Object, len(list))
	for i, item := range list {
		evalList[i] = Eval(item, env)
	}
	return evalList
}

func call(function Object, args []Object, env *Env) Object {
	if function.Type() == primitveT {
		return function.Call(args, env)
	}
	panic("invalid call to func")
}
