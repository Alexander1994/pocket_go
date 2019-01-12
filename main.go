package main

import "os"

func main() {
	if len(os.Args) < 2 {
		println("enter file name")
		return
	}
	fname := os.Args[1]
	Run(fname)
}

func Run(fname string) {
	reader.Load(fname)
	var obj *Object
	env := CreateEnv()
	for {
		obj = ParseExpr()
		if obj == nilObj {
			break
		}
		if obj == closeParenObj {
			panic("extra paren hanging out")
		}
		obj = Eval(obj, env)
		obj.Print()
		println()
	}
}
