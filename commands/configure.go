package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/.configurizer", homeDir)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s/.configurizer", homeDir), os.ModeDir)
	}

	// https://raw.githubusercontent.com/FerretCode/configurizer/main/providers/railway.yml

	client := http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://raw.githubusercontent.com/FerretCode/configurizer/main/providers/%s.yml",
			config["provider"],
		),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	providerByte, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	newFile, err := os.Create(fmt.Sprintf("%s/.configurizer/%s.yml", homeDir, config["provider"]))

	if err != nil {
		log.Fatal(err)
	}

	defer newFile.Close()

	_, err = newFile.Write(providerByte)

	if err != nil {
		log.Fatal(err)
	}

	provider := make(map[string]interface{})

	err = yaml.Unmarshal(providerByte, &provider)

	if err != nil {
		log.Fatal(err)
	}

	for k := range provider["requiredFields"].(map[string]string) {
		if config[k] == "" {
			log.Fatalf("Required field %s was not provided.\n", k)
		}
	}

	fmt.Println(config)
}
