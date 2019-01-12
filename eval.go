package main

func Eval(o *Object, env *Env) *Object {
	switch o.objT {
	case NumT, PrimitveT, NilT, FuncT, ChanT:
		return o
	case SymbolT:
		obj, _ := env.Find(o.Symbol())
		if obj == nilObj {
			panic("undefined symbol " + o.Symbol())
		}
		return obj
	case CellT:
		function := Eval(o.Car(), env)
		args := o.Cdr()
		return Call(function, args, env)
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

func Call(function *Object, args []*Object, env *Env) *Object {
	switch function.Type() {
	case PrimitveT:
		return function.CallPrim(args, env)
	case FuncT:
		return function.CallFunc(args, env)
	case MacroT:
		return function.RunMacro(args, env)
	}
	panic("Head of cell/list is not a function: ")
}
