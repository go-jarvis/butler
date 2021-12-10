package main

import "github.com/go-jarvis/jarvis/cmd/jarvis/cmd"

func main() {

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
