package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	infile := flag.String("input", "", "input json file for loading secrets")
	flag.Parse()

	if infile == nil || len(*infile) < 1 {
		fmt.Println("flag required: -input <filename>")
		flag.Usage()
		os.Exit(2)
	}

	var secretStore SecretStore
	if err := readSecretStore(*infile, &secretStore); err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	if err := secretStore.UpdateSecrets(ctx); err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}
