package main

import (
	"fmt"

	config "github.com/eckeriaue/golang-url-shortener/internal"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
