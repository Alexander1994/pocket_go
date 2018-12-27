package main

func Eval(o Object) Object {
	switch o.objT {
	case numT, primitveT, nilT:
		return o
	case cellT:
		function := Eval(o.Car())
		args := Eval_List(o.Cdr())
		if function.Type() != primitveT {
			panic("Head of cell/list is not a function")
		}
		return call(function, args)
	}
	o.print()
	return nilObj
}

func Eval_List(list []Object) []Object {
	evalList := make([]Object, len(list))
	for i, item := range list {
		evalList[i] = Eval(item)
	}
	return evalList
}

func call(function Object, args []Object) Object {
	if function.Type() == primitveT {
		return function.Call(args)
	}
	panic("invalid call to func")
}
