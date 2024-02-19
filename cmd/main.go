package main

import (
	"fmt"
	"go/parser"
	"go/token"
	// "os"
)

func main() {
	
	fst := token.NewFileSet()
	// file := os.Getenv("GOFILE")
	// fmt.Println(file)
	f, err := parser.ParseFile(fst, "/home/ethanperry/work/gobar/cmd/main.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	
	fmt.Print(f)
}