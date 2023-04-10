package main

import (
	"log"
	"os"

	"github.com/ferretcode/configurizer/commands"
)

func main() {
	switch command := os.Args[1]; command {
	case "configure":
		if len(os.Args) < 3 {
			log.Fatal("You need to provide a path to your config file!")
		}

		commands.Configure(os.Args[2])
	case "--help":
		log.Println(
			"Valid commands are `configure` and `new-provider`.\nconfigure\tdeploys your project\nnew-provider\tcreates a new provider",
		)
	case "test":
		commands.Test()
	default:
		log.Fatal("Please provide a valid command!")
	}
}
