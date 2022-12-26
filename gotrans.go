package main

import (
	"fmt"
	"github.com/Harmos274/gotrans/warehouse"
	"log"
	"os"
)

func main() {
	file, err := os.Open("testMap.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	initWr, cycles, err := parseInputFile(file)

	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan warehouse.Warehouse)

	go warehouse.CleanWarehouse(initWr, ch, cycles)

	for wr := range ch {
		fmt.Print(ShowableWarehouse(wr))
	}
	return
}
