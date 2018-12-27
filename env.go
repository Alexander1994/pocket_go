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
			return &nilObj
		}
		it = it.upperEnv
	}
}

func (e *Env) Add(symbol string, obj *Object) {
	(*e).vars[symbol] = obj
}

func AddAndGetNewEnv(e *Env) (eNew *Env) {
	eNew = new(Env)
	(*eNew).upperEnv = e
	return eNew
}
