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

        if config[fieldName] == "" {
          log.Fatalf(
            "Required field %s is not provided.\n", 
            v.(map[interface{}]interface{})["fieldName"].(string),
          )
        }

        commandString = strings.ReplaceAll(
          commandString, 
          string(match), 
          config[fieldName].(string),
        )    
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
      fmt.Println(v.(map[interface{}]interface{})["capture"])

			capture := v.(map[interface{}]interface{})["capture"].([]interface{})

			pattern := capture[0].(map[interface{}]interface{})["regex"]

      fmt.Println(pattern)

      regex, err := regexp.Compile(pattern.(string)) 

      if err != nil {
        log.Fatal(err)
      }

      output, err := cmd.Output()

      if err != nil {
        log.Fatal(err)
      }

      fmt.Println(string(output))

			captured := regex.Find(output)

			if len(captured) > 0 {
				captures = append(captures, string(captured))
			}
		}
	}
}
