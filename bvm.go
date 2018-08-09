package main

import (
        "fmt"
        "os"
        "log"
        "github.com/llir/llvm/asm"
        "github.com/llir/llvm/ir"
        "github.com/llir/llvm/ir/constant"
        //"github.com/llir/llvm/ir/types"
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

func getValue(v value.Value) string {
        switch val := v.(type) {
        case value.Named:
                return val.GetName()
        case constant.Constant:
                return getConstant(val)
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
                return "$" + val.GetName()
        case constant.Constant:
                return getConstant(val)
        default:
                return ""
        }
}

func printInstruction(inst ir.Instruction) {
        switch inst := inst.(type) {
        case *ir.InstAlloca:
                return
        case *ir.InstLoad:
                fmt.Printf("r%s=\n", inst.Name)
        case *ir.InstStore:
                fmt.Printf("r%s=%s\n", getDstValue(inst.Dst), getSrcValue(inst.Src))
        case *ir.InstAdd:
                fmt.Println("add")
                return
        case *ir.InstSub:
                fmt.Println("sub")
                return
        default:
                return
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
                for _, f := range parsedAsm.Funcs {
                        convertFuncToBash(f)
                }
        }

        fmt.Println("main")
}

