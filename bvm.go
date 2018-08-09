package main

import (
        "fmt"
        "os"
        "log"
        "github.com/llir/llvm/asm"
        "github.com/llir/llvm/ir"
        //"github.com/llir/llvm/ir/constant"
        //"github.com/llir/llvm/ir/types"
        //"github.com/llir/llvm/ir/value"
)

func printUsage() {
        fmt.Printf("blessedvirginmary usage:\n")
        fmt.Printf("    blessedvirginmary <input>.ir [<input2>.ir [<input3>.ir ... ] ]\n")
}

func printFuncSig(f *ir.Function) {
        // Print the function signature
        fmt.Println(f.Name)
}

func printFuncBlock(b *ir.BasicBlock) {
        // Loop over functions and print them out in bash
}

func printFuncEnd() {
        fmt.Println("}")
}


func convertFuncToBash(f *ir.Function) {
        printFuncSig(f)
        // for block in f.blocks
        //      printFuncBlock(f.block)
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
}

