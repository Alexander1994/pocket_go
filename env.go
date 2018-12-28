package main

type Env struct {
	vars     map[string]*Object
	upperEnv *Env
}

func createEnv() (e *Env) {
	e = new(Env)
	e.vars = make(map[string]*Object)
	return e
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
	(*e).vars[symbol] = obj
}

func AddAndGetNewEnv(e *Env) (eNew *Env) {
	eNew = createEnv()
	(*eNew).upperEnv = e
	return eNew
}

func (e *Env) popFuncEnv() {
	e = e.upperEnv
}

func (o *Object) pushFuncEnv(args *[]Object, env *Env) (newEnv *Env) {
	funcDef := o.Function()
	defArgs := *funcDef.args.List()
	if len(defArgs) != len(*args) {
		panic("args in call to function != function args")
	}
	evalArgs := EvalList(args, env)
	newEnv = AddAndGetNewEnv(env)

	for i, arg := range *evalArgs {
		newEnv.Add(defArgs[i].Symbol(), &arg)
	}
	return newEnv
}
