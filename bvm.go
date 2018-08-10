package main

import (
        "fmt"
        "strings"
        "os"
        "log"
        "strconv"
        //"github.com/kr/pretty"
        "github.com/llir/llvm/asm"
        "github.com/llir/llvm/ir"
        "github.com/llir/llvm/ir/constant"
        "github.com/llir/llvm/ir/types"
        "github.com/llir/llvm/ir/value"
)

func printUsage() {
        fmt.Printf("blessedvirginmary usage:\n")
        fmt.Printf("    blessedvirginmary <input>.ir [<input2>.ir [<input3>.ir ... ] ]\n")
}

func printFuncSig(f *ir.Function) {
        fmt.Printf("%s() {\n", f.Name)
}

func getConstant(c constant.Constant) string {
        switch c := c.(type) {
        case *constant.Int:
                return fmt.Sprintf("%d", c.X)
        default:
                return ""
        }
}

func getDstValue(v value.Value) string {
        switch val := v.(type) {
        case value.Named:
                return val.GetName()
        case constant.Constant:
                return getConstant(val)
        default:
                return ""
        }
}

func getSrcValue(v value.Value) string {
        switch val := v.(type) {
        case value.Named:
                return "${r" + val.GetName() + "}"
        case constant.Constant:
                return getConstant(val)
        default:
                return ""
        }
}

func printIcmp(inst ir.InstICmp) {
        var op string
        switch inst.Pred {
        case ir.IntNE:
                        op = " != "
        default:
                        op = ""
        }

        fmt.Printf("r%s=`if [ \"$r\"%s%s\"$r\"%s ]; then echo false; else echo true; fi`\n", inst.Name, getDstValue(inst.X), op, getDstValue(inst.Y))
        return
}

func instAllocaHelper(inst *ir.InstAlloca) string {
        switch t := inst.Typ.Elem.(type) {
        case *types.ArrayType:
                prefilled := strings.Repeat("0 ", int(t.Len))
                return fmt.Sprintf("s%s=(%s)\nr%s=s%s\n", inst.Name, prefilled, inst.Name, inst.Name)
        default:
                return fmt.Sprintf("declare s%s\nr%s=s%s\n", inst.Name, inst.Name, inst.Name)
        }
        return ""
}

func printInstruction(inst ir.Instruction) {
        switch inst := inst.(type) {
        /* Memory Instructions */

        case *ir.InstAlloca:
                fmt.Printf("%s", instAllocaHelper(inst))
                return
        case *ir.InstLoad:
                fmt.Printf("eval r%s=\\${%s}\n", inst.Name, getSrcValue(inst.Src))
                return
        case *ir.InstStore:
                fmt.Printf("eval %s=%s\n", getSrcValue(inst.Dst), getSrcValue(inst.Src))
                return
        case *ir.InstGetElementPtr:
                index, err := strconv.Atoi(getSrcValue(inst.Indices[1]))
                if err != nil {
                        panic("")
                }
                fmt.Printf("r%s=%s[%d]\n", inst.Name, getSrcValue(inst.Src), index)
                return

        /* Math Instructions */

        case *ir.InstAdd:
                fmt.Printf("r%s=$(expr %s + %s)\n", inst.Name, getSrcValue(inst.X), getSrcValue(inst.Y))
                return
        case *ir.InstSub:
                fmt.Printf("r%s=$(expr %s - %s)\n", inst.Name, getSrcValue(inst.X), getSrcValue(inst.Y))
                return
        case *ir.InstMul:
                fmt.Printf("r%s=$(expr %s \\* %s)\n", inst.Name, getSrcValue(inst.X), getSrcValue(inst.Y))
                return
        case *ir.InstSDiv:
                fmt.Printf("r%s=$(expr %s / %s)\n", inst.Name, getSrcValue(inst.X), getSrcValue(inst.Y))
                return
        case  *ir.InstICmp:
                printIcmp(*inst)
                return
        case *ir.InstSRem:
                fmt.Printf("r%s=$(expr %s %% %s)\n", inst.Name, getSrcValue(inst.X), getSrcValue(inst.Y))
                return
        // What about UDiv, URem?
        // Floating Point Instructions?

        default:
                panic(fmt.Sprintf("Unknown instruction %s", inst))
        }
}

func printFuncBlock(b *ir.BasicBlock) {
        for _, inst := range b.Insts {
                printInstruction(inst)
        }
        switch term := b.Term.(type) {
        case *ir.TermRet:
                fmt.Printf("return %s\n", getSrcValue(term.X))
		case *ir.TermCondBr:
				fun1 := "_br" + term.TargetTrue.GetName() + term.TargetTrue.Parent.GetName()
				fun2 := "_br" + term.TargetFalse.GetName() + term.TargetFalse.Parent.GetName()
				fmt.Printf("if [ $r%s ]; then %s; else %s; fi\n", getDstValue(term.Cond), fun1, fun2)
        }
}

func convertFuncToBash(f *ir.Function) {
        // Top level function
        fmt.Printf("%s() {\n", f.Name)
        fmt.Printf("%s\n", "_br" + f.Name + f.Blocks[0].GetName())
        fmt.Printf("}\n")

        // Blocks
        for _, block := range f.Blocks {
                fmt.Printf("%s() {\n", "_br" + f.GetName() + f.Blocks[0].GetName())
                printFuncBlock(block)
                fmt.Printf("}\n")
        }
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

        fmt.Println("main")
}

