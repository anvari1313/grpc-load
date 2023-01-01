package main

import (
	"fmt"
	"os"

	"github.com/anvari1313/grpc-load/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}
