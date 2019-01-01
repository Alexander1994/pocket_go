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
	Load(fname)
	var obj *Object
	env := createEnv()
	for {
		obj = parseExpr()
		if obj == nilObj {
			break
		}
		if obj == closeParenObj {
			panic("extra paren hanging out")
		}
		obj = Eval(obj, env)
		obj.print()
		println()
	}
}
