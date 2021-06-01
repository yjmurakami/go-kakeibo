package main

import (
	"fmt"
	"os"

	"github.com/yjmurakami/go-kakeibo/cmd/api/server"
)

func main() {
	if err := server.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
