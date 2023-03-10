package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "  "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func main() {
	prettyPtr := flag.Bool("pretty", true, "display pretty output; otherwise do: -pretty=false")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Usage: scrub-pii <input json file> <sensitive fields file>")
		os.Exit(1)
	}

	inputPath := flag.Args()[0]
	sensitiveFieldsPath := flag.Args()[1]

	// scrub the input file for the given sensitive fields
	output, err := ScrubPersonalInfo(inputPath, sensitiveFieldsPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// pretty output if -pretty is specified
	if *prettyPtr {
		pretty, err := PrettyString(output)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println(pretty)
		return
	}

	// just output a json string
	fmt.Println(output)
}
