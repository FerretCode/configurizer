package commands

import (
	"fmt"
	"log"
	"os/exec"
)

func Test() {
	cmd := exec.Command("railway", "up")

	if cmd.Err != nil {
		log.Fatal(cmd.Err.Error())
	}

	byte, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(byte))
}
