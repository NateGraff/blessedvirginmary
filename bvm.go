package main

import (
        "fmt"
        "strings"
        "os"
        "log"
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
                return "$r" + val.GetName()
        case constant.Constant:
                return getConstant(val)
        default:
                return ""
        }
}

func instAllocaHelper(inst *ir.InstAlloca) string {
        switch t := inst.Typ.Elem.(type) {
        case *types.ArrayType:
                prefilled := strings.Repeat("0 ", int(t.Len))
                return fmt.Sprintf("r%s=(%s)\n", inst.Name, prefilled)
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
                fmt.Printf("eval r%s=\\$%s\n", inst.Name, getSrcValue(inst.Src))
                return
        case *ir.InstStore:
                fmt.Printf("eval $r%s=%s\n", getDstValue(inst.Dst), getSrcValue(inst.Src))
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
        }
}

func printFuncEnd() {
        fmt.Println("}")
}


func convertFuncToBash(f *ir.Function) {
        printFuncSig(f)
        for _, block := range f.Blocks {
                printFuncBlock(block)
        }
        printFuncEnd()
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

