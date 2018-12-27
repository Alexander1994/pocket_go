package main

func main() {
	// if len(os.Args) != 1 {
	// 	print("enter file name")
	// }
	fname := /*os.Args[1]*/ "main" + ".pgo"
	Run(fname)
}

func Run(fname string) {
	Load(fname)

	var obj Object
	for {
		obj = readExpr()
		if obj == nilObj {
			break
		}
		if obj == closeParenObj {
			panic("extra paren hanging out")
		}
		obj = Eval(obj)
		obj.print()
		println()
	}
}
