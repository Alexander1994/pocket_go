package main

import "fmt"

type varT int

const (
	nilT varT = iota
	trueT
	closePT
	numT
	cellT
	symbolT
	primitveT
	funcT
	// not in use
	macroT
)

type Prim struct {
	name   string
	method PrimFunc
}

type Object struct {
	objT  varT
	value interface{}
}

type Func struct {
	name string    // $symbol
	args *Object   // ($symbol...)
	expr *[]Object // $expr...
}

var nilObj = &Object{objT: nilT}
var closeParenObj = &Object{objT: closePT}
var trueObj = &Object{objT: trueT}

func prints(os []Object) {
	for _, op := range os {
		op.print()
	}
}

func (o *Object) print() {
	switch o.objT {
	case primitveT, symbolT:
		fmt.Printf("%s ", o.Symbol())
	case closePT:
		fmt.Print(") ")
	case nilT:
		fmt.Printf("nil ")
	case numT:
		fmt.Printf("%.2f ", o.Num())
	case cellT:
		fmt.Print("\n")
		for _, no := range *o.List() {
			no.print()
		}
	default:
		panic("invalid type found: " + string(o.objT))
	}
}

func (o *Object) Type() varT {
	return o.objT
}

func (o *Object) List() *[]Object {
	return o.value.(*[]Object)
}

func (o *Object) Car() *Object {
	return &(*o.List())[0]
}

func (o *Object) Cdr() *[]Object {
	cdr := ((*o.List())[1:])
	return &cdr
}

func Cdr(os *[]Object) *[]Object {
	cdr := (*os)[1:]
	return &cdr
}

func Car(os *[]Object) *Object {
	return &(*os)[0]
}

func (o *Object) Num() float32 {
	return o.value.(float32)
}

func (o *Object) Symbol() string {
	if o.objT == primitveT {
		return o.value.(Prim).name
	} else if o.objT == symbolT {
		return o.value.(string)
	}
	panic("no symbol found")
}

func (o *Object) Function() *Func {
	return o.value.(*Func)
}

func (o *Object) CallPrim(args *[]Object, env *Env) *Object {
	return o.value.(*Prim).method(args, env)
}

// Create Objects
func Num(n float32) *Object {
	return &Object{objT: numT, value: n}
}

func List(os *[]Object) *Object {
	return &Object{objT: cellT, value: os}
}

func Primitve(name string) *Object {
	return &Object{objT: primitveT, value: &Prim{method: Functs[name], name: name}}
}

func Symbol(name string) *Object {
	return &Object{objT: symbolT, value: name}
}

func Function(name string, args *Object, expr *[]Object) *Object {
	return &Object{objT: funcT, value: &Func{name: name, args: args, expr: expr}}
}
