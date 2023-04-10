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

	/*homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}*/

	homeDir := "."

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

	captures := make(map[string]string)

	for _, v := range provider["steps"].([]interface{}) {
		commandString := v.(map[interface{}]interface{})["command"].(string)

		regex, err := regexp.Compile(`\{([^}]+)\}`)

		if err != nil {
			log.Fatal(err)
		}

		match := regex.Find(
			[]byte(
				v.(map[interface{}]interface{})["command"].(string),
			),
		)

		if len(match) > 0 {
			for _, v := range provider["requiredFields"].([]interface{}) {
				fieldName := v.(map[interface{}]interface{})["fieldName"].(string)
				name := v.(map[interface{}]interface{})["name"]

				if config[fieldName] == "" {
					log.Fatalf(
						"Required field %s is not provided.\n",
						v.(map[interface{}]interface{})["fieldName"].(string),
					)
				}

				if regex.FindStringSubmatch(commandString)[1] == name.(string) {
					commandString = strings.ReplaceAll(
						commandString,
						string(match),
						config[fieldName].(string),
					)
				}
			}

			// re-check regex
			match := regex.FindStringSubmatch(commandString)

			if len(match) > 1 {
				if captures[match[1]] != "" {
					commandString = strings.ReplaceAll(
						commandString,
						match[0],
						captures[match[1]],
					)
				}
			}
		}

		command := strings.Split(
			commandString,
			" ",
		)

		cmd := exec.Command(command[0], command[1:]...)

		if cmd.Err != nil {
			log.Fatal(cmd.Err.Error())
		}

		if v.(map[interface{}]interface{})["capture"] != nil {
			capture := v.(map[interface{}]interface{})["capture"].(map[interface{}]interface{})

			pattern := capture["regex"]

			regex, err := regexp.Compile(pattern.(string))

			if err != nil {
				log.Fatal(err)
			}

			outputByte, err := cmd.Output()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(outputByte))

			output := strings.TrimSuffix(string(outputByte), "\n")

			captured := regex.FindStringSubmatch(output)

			if len(captured) > 1 {
				captures[capture["name"].(string)] = string(captured[1])
			}
		}
	}
}
