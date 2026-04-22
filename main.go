package main

import (
	"fmt"
	"log"
	"os"

	"github.com/itshivams/studex-cli/cmd"
	"github.com/itshivams/studex-cli/internal/config"
)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-version" || arg == "--version" || arg == "-v" || arg == "version" {
			fmt.Println("studex-cli version 1.0.1")
			os.Exit(0)
		}
	}

	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	cmd.Execute()
}
