package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	//"github.com/kr/pretty"
	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func printUsage() {
	fmt.Printf("blessedvirginmary usage:\n")
	fmt.Printf("    blessedvirginmary <input>.ir [<input2>.ir [<input3>.ir ... ] ]\n")
}

func printFuncSig(f *ir.Function) {
	fmt.Printf("%s() {\n", name(f))
}

func getConstant(c constant.Constant) string {
	switch c := c.(type) {
	case *constant.Int:
		return fmt.Sprintf("%d", c.X)
	default:
		return ""
	}
}

func getLValue(v value.Value) string {
	switch val := v.(type) {
	case value.Named:
		return "local[r" + name(val) + "]"
	case constant.Constant:
		return getConstant(val)
	default:
		return ""
	}
}

func getBareName(v value.Value) string {
	switch val := v.(type) {
	case value.Named:
		return name(val)
	case constant.Constant:
		return getConstant(val)
	default:
		return ""
	}
}

func getRValue(v value.Value) string {
	switch val := v.(type) {
	case value.Named:
		return "${local[r" + name(val) + "]}"
	case constant.Constant:
		return getConstant(val)
	default:
		return ""
	}
}

func printIcmp(inst *ir.InstICmp) {
	var operand string
	switch inst.Pred {
	case enum.IPredNE:
		operand = "ne"
	case enum.IPredEQ:
		operand = "eq"
	case enum.IPredUGT:
	case enum.IPredSGT:
		operand = "gt"
	case enum.IPredUGE:
	case enum.IPredSGE:
		operand = "ge"
	case enum.IPredULT:
	case enum.IPredSLT:
		operand = "lt"
	case enum.IPredULE:
	case enum.IPredSLE:
		operand = "le"
	default:
		operand = "eq"
	}
	fmt.Printf("local[r%s]=`if [ \"%s\" -%s \"%s\" ]; then echo true; fi`\n", name(inst), getRValue(inst.X), operand, getRValue(inst.Y))
	return
}

func instAllocaHelper(inst *ir.InstAlloca) {
	switch t := inst.Typ.ElemType.(type) {
	case *types.ArrayType:
		for idx := 0; idx < int(t.Len); idx++ {
			fmt.Printf("local[s%s_%d]=0;", name(inst), idx)
		}
		fmt.Printf("\n%s=s%s\n", getLValue(inst), name(inst))
	default:
		fmt.Printf("local[s%s]=0\nlocal[r%s]=s%s\n", name(inst), name(inst), name(inst))
	}
}

func printInstruction(inst ir.Instruction) {
	switch inst := inst.(type) {
	/* Memory Instructions */

	case *ir.InstAlloca:
		instAllocaHelper(inst)
		return
	case *ir.InstLoad:
		fmt.Printf("local[r%s]=${local[%s]}\n", name(inst), getRValue(inst.Src))
		return
	case *ir.InstStore:
		fmt.Printf("local[%s]=%s\n", getRValue(inst.Dst), getRValue(inst.Src))
		return
	case *ir.InstGetElementPtr:
		index, err := strconv.Atoi(getRValue(inst.Indices[1]))
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s=%s_%d\n", getLValue(inst), getRValue(inst.Src), index)
		return

	case *ir.InstCall:
		// create callee function context and install arguments starting at index r0
		fmt.Printf("declare -A args\n")
		for idx, arg := range inst.Args {
			fmt.Printf("args[r%d]=%s\n", idx, getRValue(arg))
		}
		fmt.Printf("eval `%s \"$(declare -p args)\"`\n", getBareName(inst.Callee))
		fmt.Printf("%s=${ret}\n", getLValue(inst))

	/* Math Instructions */

	case *ir.InstAdd:
		fmt.Printf("%s=$(expr %s + %s)\n", getLValue(inst), getRValue(inst.X), getRValue(inst.Y))
		return
	case *ir.InstSub:
		fmt.Printf("%s=$(expr %s - %s)\n", getLValue(inst), getRValue(inst.X), getRValue(inst.Y))
		return
	case *ir.InstMul:
		fmt.Printf("%s=$(expr %s \\* %s)\n", getLValue(inst), getRValue(inst.X), getRValue(inst.Y))
		return
	case *ir.InstSDiv:
		fmt.Printf("%s=$(expr %s / %s)\n", getLValue(inst), getRValue(inst.X), getRValue(inst.Y))
		return
	case *ir.InstICmp:
		printIcmp(inst)
		return
	case *ir.InstSRem:
		fmt.Printf("%s=$(expr %s %% %s)\n", getLValue(inst), getRValue(inst.X), getRValue(inst.Y))
		return
	// What about UDiv, URem?
	// Floating Point Instructions?

	default:
		panic(fmt.Sprintf("Unknown instruction %s", inst))
	}
}

func getCondTermName(term ir.Terminator) string {
	switch term := term.(type) {
	case *ir.TermBr:
		return name(term.Target)
	default:
		return ""
	}
}

func printFuncBlock(b *ir.BasicBlock, funcname string) {
	for _, inst := range b.Insts {
		printInstruction(inst)
	}
	switch term := b.Term.(type) {
	case *ir.TermRet:
		fmt.Printf("local[ret]=%s\n", getRValue(term.X))
	case *ir.TermCondBr:
		fun1 := "_br" + funcname + name(term.TargetTrue)
		fun2 := "_br" + funcname + name(term.TargetFalse)
		fmt.Printf("if [ %s ]; then\n", getRValue(term.Cond))
		fmt.Printf("  eval `%s \"$(declare -p local)\"`\n", fun1)
		switch targetTerm := term.TargetTrue.Term.(type) {
		case *ir.TermBr:
			fmt.Printf("  eval `%s \"$(declare -p local)\"`\n", "_br"+funcname+name(targetTerm.Target))
		}
		fmt.Printf("else\n")
		fmt.Printf("  eval `%s \"$(declare -p local)\"`\n", fun2)
		switch targetTerm := term.TargetFalse.Term.(type) {
		case *ir.TermBr:
			fmt.Printf("  eval `%s \"$(declare -p local)\"`\n", "_br"+funcname+name(targetTerm.Target))
		}
		fmt.Printf("fi\n")
	}
}

func convertFuncToBash(f *ir.Function) {
	if len(f.Blocks) == 0 {
		return
	}
	// Assign IDs to unnamed local variables.
	if err := f.AssignIDs(); err != nil {
		panic(err)
	}
	// Top level function
	fmt.Printf("%s() {\n", name(f))
	fmt.Printf("declare -A local=${1#*=}\n")
	fmt.Printf("eval `%s\n", "_br"+name(f)+name(f.Blocks[0])+" \"$(declare -p local)\"`")
	fmt.Printf("ret=${local[ret]}\n")
	fmt.Printf("declare -p ret\n")
	fmt.Printf("}\n")

	// Blocks
	for _, block := range f.Blocks {
		fmt.Printf("%s() {\n", "_br"+name(f)+name(block))
		fmt.Printf("declare -A local=${1#*=}\n")
		printFuncBlock(block, name(f))
		fmt.Printf("declare -p local\n")
		fmt.Printf("}\n")
	}
}

// name returns the name or ID of the given value.
func name(v value.Named) string {
	const prefix = "%"
	return v.Ident()[len(prefix):]
}

func main() {
	args := os.Args[1:] // Slice the args after the program name
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	for i := 0; i < len(args); i++ {
		parsedAsm, err := asm.ParseFile(args[i])
		if err != nil {
			log.Fatal(err)
		}
		//pretty.Println(parsedAsm)
		for _, f := range parsedAsm.Funcs {
			convertFuncToBash(f)
		}
	}

	fmt.Println("declare -A args")
	// insert arguments to main here
	fmt.Println("eval `main \"$(declare -p args)\"`")
	fmt.Println("exit ${ret}")
}
