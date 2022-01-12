package main

import (
	"github.com/ubiq/ubiq-burn-stats/daemon/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
