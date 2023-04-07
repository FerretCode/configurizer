package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

func Configure(path string) {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("There was an error reading the supplied config file. Please make sure the path is correct.")
	}

	config := make(map[string]interface{})

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}
