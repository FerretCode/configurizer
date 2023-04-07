package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

func NewProvider(path string) {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("There was an error reading the supplied provider file. Please make sure the path is correct.")
	}

	provider := make(map[string]interface{})

	err = yaml.Unmarshal(file, &provider)

	if err != nil {
		log.Fatal(err)
	}

	for k, v := range provider {
		fmt.Println(k, v)
	}
}
