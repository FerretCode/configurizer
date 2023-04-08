package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
		os.Mkdir(fmt.Sprintf("%s/.configurizer", homeDir), 0777)
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

	var captures []string

	for _, v := range provider["steps"].([]interface{}) {
		command := strings.Split(
			v.(map[interface{}]interface{})["command"].(string),
			" ",
		)

		regex := regexp.Regexp{}

		match := regex.Find(
			[]byte(
				v.(map[interface{}]interface{})["command"].(string),
			),
		)

		if len(match) > 0 {
			for _, v := range provider["requiredFields"].([]interface{}) {
				for key := range v.(map[interface{}]interface{}) {
					if config[key.(string)] == "" {
						log.Fatalf("Required field %s is not provided.\n", key)
					}
				}
			}
		}

		err := exec.Command(command[0], command[1:]...)

		if err != nil {
			log.Fatal(err)
		}

		if v.(map[interface{}]interface{})["capture"] != nil {
			capture := v.(map[interface{}]interface{})["capture"].(map[string]string)

			pattern := capture["regex"]

			captured := regex.Find([]byte(pattern))

			if len(captured) > 0 {
				captures = append(captures, string(captured))
			}
		}
	}
}
