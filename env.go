package main

type Env struct {
	vars     map[string]*Object
	upperEnv *Env
}

func (e *Env) dump() {
	it := e
	for {
		for name, v := range e.vars {
			print(name + " = ")
			v.print()
			println()
		}
		if it.upperEnv == nil {
			break
		}
		it = it.upperEnv
	}
}

func createEnv() (e *Env) {
	return &Env{vars: make(map[string]*Object), upperEnv: nil}
}

func (env *Env) find(symbol string) (obj *Object) {
	var found bool
	it := env
	for {
		obj, found = it.vars[symbol]
		if found {
			return obj
		}
		if it.upperEnv == nil {
			return nilObj
		}
		it = it.upperEnv
	}
}

func (e *Env) Add(symbol string, obj *Object) {
	e.vars[symbol] = obj
}

func AddAndGetNewEnv(e *Env) (eNew *Env) {
	eNew = createEnv()
	eNew.upperEnv = e
	return eNew
}

func (e *Env) popFuncEnv() {
	e = e.upperEnv
}

func (o *Object) pushFuncEnv(args []*Object, env *Env) (newEnv *Env) {
	funcDef := o.Function()
	defArgs := funcDef.args.List()
	if len(defArgs) != len(args) {
		panic("args in call to function != function args")
	}
	evalArgs := EvalList(args, env)
	newEnv = AddAndGetNewEnv(env)
	// (defn double (x ch)
	// (<- ch (* 2 x) )
	for i, arg := range evalArgs {
		name := defArgs[i].Symbol()
		newEnv.Add(name, arg)
	}
	return newEnv
}
