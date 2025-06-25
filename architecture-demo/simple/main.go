package main

import (
	"fmt"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := GetConfiguration()
	for i := 1; i <= cfg.Limit; i++ {
		fizzbuzz := GetFizzBuzz(i)
		fmt.Printf("%s ", fizzbuzz)
	}
	fmt.Println()
	return nil
}
