package main

import (
	"github.com/clutterpot/clutterpot/cmd"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cmd.Execute()
}
