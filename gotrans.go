package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Harmos274/gotrans/warehouse"
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

	actualCycle := 1
	for wr := range ch {
		fmt.Println("tour", actualCycle)
		fmt.Println(ShowableWarehouse(wr))
		actualCycle += 1
	}
	return
}
