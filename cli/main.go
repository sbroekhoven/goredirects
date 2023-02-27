package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sbroekhoven/goredirects"
)

func main() {
	url := os.Args[1]

	rd := goredirects.Get(url, "1.1.1.1")

	json, err := json.MarshalIndent(rd, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", json)
}
