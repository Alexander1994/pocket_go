package main

import (
	"fmt"
	"strconv"
)

type varT int

const (
	NilT varT = iota
	ClosePT
	NumT
	CellT
	SymbolT
	ChanT
	PrimitveT
	FuncT
	MacroT
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
	expr         []*Object   // $expr...
	templateargs [][]*Object // argindex => array of references to arg in code
}

var nilObj = &Object{objT: NilT}
var closeParenObj = &Object{objT: ClosePT}

func Prints(os []*Object) {
	for _, op := range os {
		op.Print()
		print(" ")
	}
}

func (o *Object) Print() {
	switch o.objT {
	case PrimitveT, SymbolT:
		fmt.Printf("%s", o.Symbol())
	case FuncT:
		print("(")
		o.value.(*Func).args.Print()
		print(")")
		Prints(o.value.(*Func).expr)
	case ClosePT:
		fmt.Print(")")
	case NilT:
		fmt.Printf("nil")
	case NumT:
		fmt.Printf("%.2f", o.Num())
	case ChanT:
		fmt.Printf("chan")
	case MacroT:
		fmt.Printf("macrofn")

	case CellT:
		print("(")
		for i, no := range o.List() {
			no.Print()
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
	if o.objT == PrimitveT {
		return o.value.(*Prim).name
	} else if o.objT == SymbolT {
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
	return &Object{objT: NumT, value: n}
}

func List(os []*Object) *Object {
	return &Object{objT: CellT, value: os}
}

func Primitve(name string) *Object {
	return &Object{objT: PrimitveT, value: &Prim{method: Functs[name], name: name}}
}

func Symbol(name string) *Object {
	return &Object{objT: SymbolT, value: name}
}

func Function(args *Object, closure *Env, expr []*Object) *Object {
	return &Object{objT: FuncT, value: &Func{args: args, closure: closure, expr: expr}}
}

func Macro(tempateargs [][]*Object, expr []*Object) *Object {
	return &Object{objT: MacroT, value: &MacroFn{expr, tempateargs}}
}

func Channel() *Object {
	ch := make(chan *Object)
	return &Object{objT: ChanT, value: &ch}
}
