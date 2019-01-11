package main

import (
	"fmt"
	"strconv"
)

type varT int

const (
	nilT varT = iota
	closePT
	numT
	cellT
	symbolT
	chanT
	primitveT
	funcT
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
	args    *Object // ($symbol...)
	closure *Env
	expr    []*Object // $expr...
}

type MacroFn struct {
	expr         []*Object
	templateargs [][]*Object // argindex => array of references to arg in code
}

var nilObj = &Object{objT: nilT}
var closeParenObj = &Object{objT: closePT}

func prints(os []*Object) {
	for _, op := range os {
		op.print()
		print(" ")
	}
}

func (o *Object) print() {
	switch o.objT {
	case primitveT, symbolT:
		fmt.Printf("%s", o.Symbol())
	case funcT:
		print("(")
		o.value.(*Func).args.print()
		print(")")
		prints(o.value.(*Func).expr)
	case closePT:
		fmt.Print(")")
	case nilT:
		fmt.Printf("nil")
	case numT:
		fmt.Printf("%.2f", o.Num())
	case chanT:
		fmt.Printf("chan")
	case macroT:
		fmt.Printf("macrofn")

	case cellT:
		print("(")
		for i, no := range o.List() {
			no.print()
			if i != len(o.List())-1 {
				print(" ")
			}
		}
		print(")")
	default:
		panic("invalid type found: " + o.TypeStr())
	}
}

func (o *Object) Type() varT {
	return o.objT
}

func (o *Object) TypeStr() string {
	return strconv.Itoa(int(o.Type()))
}

func (o *Object) List() []*Object {
	return o.value.([]*Object)
}

func (o *Object) Set(i int, newO *Object) {
	o.value.([]*Object)[i] = newO
}

func (o *Object) Car() *Object {
	return o.List()[0]
}

func (o *Object) Cdr() []*Object {
	return o.List()[1:]
}

func Cdr(os []*Object) []*Object {
	return os[1:]
}

func Car(os []*Object) *Object {
	return os[0]
}

func (o *Object) Num() float32 {
	return o.value.(float32)
}

func (o *Object) Symbol() string {
	if o.objT == primitveT {
		return o.value.(*Prim).name
	} else if o.objT == symbolT {
		return o.value.(string)
	}
	panic("no symbol found")
}

func (o *Object) Macro() *MacroFn {
	return o.value.(*MacroFn)
}

func (o *Object) Function() *Func {
	return o.value.(*Func)
}

func (o *Object) CallPrim(args []*Object, env *Env) *Object {
	return o.value.(*Prim).method(args, env)
}

func (o *Object) Send(s *Object) {
	(*o.value.(*chan *Object)) <- s
}

func (o *Object) Recv() (recv *Object) {
	recv = <-(*o.value.(*chan *Object))
	return recv
}

// Create Objects
func Num(n float32) *Object {
	return &Object{objT: numT, value: n}
}

func List(os []*Object) *Object {
	return &Object{objT: cellT, value: os}
}

func Primitve(name string) *Object {
	return &Object{objT: primitveT, value: &Prim{method: Functs[name], name: name}}
}

func Symbol(name string) *Object {
	return &Object{objT: symbolT, value: name}
}

func Function(args *Object, closure *Env, expr []*Object) *Object {
	return &Object{objT: funcT, value: &Func{args: args, closure: closure, expr: expr}}
}

func Macro(tempateargs [][]*Object, expr []*Object) *Object {
	return &Object{objT: macroT, value: &MacroFn{expr, tempateargs}}
}

func Channel() *Object {
	ch := make(chan *Object)
	return &Object{objT: chanT, value: &ch}
}
