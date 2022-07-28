package cmd

import (
	"fmt"

	"github.com/clutterpot/clutterpot/app"
)

func Execute() {
	fmt.Println("Clutterpot")

	app.New().Init()
}
