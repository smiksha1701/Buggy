package repl

import (
	"bufio"
	"fmt"
	"os"

	"io"

	"github.com/smiksha1701/buggy/evaluator"
	"github.com/smiksha1701/buggy/lexer"
	"github.com/smiksha1701/buggy/object"
	"github.com/smiksha1701/buggy/parser"
)

const PROMPT = `>>`
const Buggy = `     
   _|__|__|__|__|_ 
  /               \ 
  \X              / 
   \_____________/   
   /\      
   \_\_                              
`

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	}
}
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, Buggy)
	io.WriteString(out, "Oh, my developer won't be happy to see that...\n")
	io.WriteString(out, "  Here are some parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
