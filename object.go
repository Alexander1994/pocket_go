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

	macroT
	envT
)

type Prim struct {
	name   string
	method Method
}

type Object struct {
	objT  varT
	value interface{}
}

var nilObj = Object{objT: nilT}
var closeParenObj = Object{objT: closePT}

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
		for _, no := range o.List() {
			no.print()
		}
	default:
		panic("invalid type found: " + string(o.objT))
	}
}

func (o *Object) Type() varT {
	return o.objT
}

func (o *Object) Append(newO Object) {
	o.value = append(o.value.([]Object), newO)
}

func (o *Object) List() []Object {
	return o.value.([]Object)
}

func (o *Object) Car() Object {
	return o.List()[0]
}

func (o *Object) Cdr() []Object {
	return o.value.([]Object)[1:]
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

func (o *Object) Call(args []Object, env *Env) Object {
	return o.value.(Prim).method(args, env)
}

// Create Objects
func Num(n float32) Object {
	return Object{objT: numT, value: n}
}

func List(os []Object) Object {
	return Object{objT: cellT, value: os}
}

func Primitve(name string) Object {
	return Object{objT: primitveT, value: Prim{method: Functs[name], name: name}}
}

func Symbol(name string) Object {
	return Object{objT: symbolT, value: name}
}
