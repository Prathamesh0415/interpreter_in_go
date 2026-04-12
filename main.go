package main

import (
	"os"
	"fmt"
	"interpreter_go/repl"
)

func main(){
	fmt.Println("Lexer implemented :")
	repl.Start(os.Stdin, os.Stdout)
}