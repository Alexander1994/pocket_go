package main

type Env struct {
	vars     map[string]*Object
	upperEnv *Env
}

func (e *Env) dump() {
	it := e
	for {
		for name, v := range it.vars {
			print(name + " = ")
			v.Print()
			print(" ")
		}
		println()
		if it.upperEnv == nil {
			break
		}
		it = it.upperEnv
	}
}

func (env *Env) Find(symbol string) (obj *Object, objEnv *Env) {
	var found bool
	it := env
	for {
		obj, found = it.vars[symbol]
		if found {
			return obj, it
		}
		if it.upperEnv == nil {
			return nilObj, nil
		}
		it = it.upperEnv
	}
}

func (e *Env) Add(symbol string, obj *Object) {
	e.vars[symbol] = obj
}

func (e *Env) Set(symbol string, obj *Object) {
	currObj, env := e.Find(symbol)
	if currObj != nilObj {
		env.Add(symbol, obj)
		return
	}
	panic("call to set of non existant var")
}

func AddAndGetNewEnv(e *Env) (eNew *Env) {
	eNew = CreateEnv()
	eNew.upperEnv = e
	return eNew
}

func CreateEnv() (e *Env) {
	return &Env{vars: make(map[string]*Object), upperEnv: nil}
}

func (e *Env) PopFuncEnv() {
	e = e.upperEnv
}

func (e *Env) IsTempEnv() bool {
	return e.upperEnv != nil
}

func (o *Object) PushFuncEnv(args []*Object, env *Env) (newEnv *Env) {
	defArgs := o.Function().args.List()
	if len(defArgs) != len(args) {
		panic("args in call to function != function args")
	}
	evalArgs := EvalList(args, env)
	newEnv = AddAndGetNewEnv(env)

	for i, arg := range evalArgs {
		name := defArgs[i].Symbol()
		newEnv.Add(name, arg)
	}
	return newEnv
}
