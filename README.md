# LLVM IR to Bash 4.x Transpiler

## Why?

Because

## Requirements

### System Utilities

- Bash 4.x (version 4 required for associative array support)
- GNU Parallel
- Make
- Clang
- Go

### Go Packages

- `go get github.com/llir/llvm/ir`
- `go get github.com/llir/llvm/asm`

## Getting Started

1. `mkdir -p ~/go/src && cd ~/go/src`
2. `git clone https://github.com/NateGraff/blessedvirginmary.git`
3. `cd blessedvirginmary`
4. `go build`
5. `bash test.bash`

