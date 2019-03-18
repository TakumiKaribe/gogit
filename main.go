package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()
	argv := flag.Args()
	if len(argv) < 1 {
		log.Fatal("insufficient arguments")
	}

	cmd := flag.Arg(0)

	switch cmd {
	case "add",
		"cat-file",
		"checkout",
		"commit",
		"hash-object",
		"init",
		"log",
		"ls-tree",
		"merge",
		"rebase",
		"rev-parse",
		"rm",
		"show-ref",
		"tag":
		command(cmd)
	}
}

func command(cmd string) {
	fmt.Printf("your input command is %s\n", cmd)
}
