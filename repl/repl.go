package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/smiksha1701/buggy/lexer"
	"github.com/smiksha1701/buggy/token"
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		l := lexer.New(scanner.Text())
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("TokenType = %q, TokenLiteral = %q\n", tok.Type, tok.Literal)
		}
	}
}
