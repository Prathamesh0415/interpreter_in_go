package repl

import (
	"bufio"
	"io"
	"fmt"
	"interpreter_go/token"
	"interpreter_go/lexer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		out.Write([]byte(PROMPT))
		scanner.Scan()
		text := scanner.Text()
		l := lexer.New(text)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("{ Type: %s :: Literal: %s }\n", string(tok.Type), tok.Literal)
		}  
	}

}