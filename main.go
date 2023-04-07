package configurizer

import "github.com/akamensky/argparse"

func main() {
	parser := argparse.NewParser("cfgr", "The configurizer command line tool")

	newProvider := parser.String("new", "new-provider", &argparse.Options{Required: false})
	path := parser.String("path", "config-path", &argparse.Options{Required: true})
}
