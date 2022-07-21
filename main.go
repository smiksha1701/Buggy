package main

import (
	"fmt"

	"github.com/smiksha1701/buggy/repl"
)

const Welcome = `How do you do, fellow kids?
Let me introduce you Buggy language
Fill free to type in your commands
`

func main() {
	fmt.Print(Welcome)

	repl.Start()
}
