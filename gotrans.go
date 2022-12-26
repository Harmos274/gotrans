package main

import (
	"fmt"
	"github.com/Harmos274/gotrans/clean_warehouse"
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

	warehouse, cycles, err := parseInputFile(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(warehouse, cycles)
	ch := make(chan clean_warehouse.Warehouse)

	go clean_warehouse.CleanWarehouse(warehouse, ch, cycles)

	for wr := range ch {
		fmt.Print(ShowWarehouse(wr))
	}
	return
}
