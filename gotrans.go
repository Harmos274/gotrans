// Responsible for parsing the input file and displaying the output
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Harmos274/gotrans/warehouse"
)

const HelpText = "Gotrans\n" +
	"=========\n" +
	"Giving a file that describes a warehouse with packages, forklifts and trucks inside it,\n" +
	"the program will have to optimise the distribution of packages to trucks using the forklifts.\n\n" +
	"Commands:\n" +
	"-h --help\tShow help\n" +
	"<file>\t\tlaunch the program\n"

func main() {
	arguments := os.Args

	if len(arguments) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, HelpText)
		os.Exit(1)
	} else if arguments[1] == "-h" || arguments[1] == "--help" {
		fmt.Printf("%s\n", HelpText)
		return
	}

	file, err := os.Open(arguments[1])
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

	ch := make(chan warehouse.CycleState)

	go warehouse.CleanWarehouse(initWr, ch, cycles)

	currentCycle := 1
	for state := range ch {
		fmt.Printf("tour %d/%d\n", currentCycle, cycles)
		fmt.Println(ShowableWarehouse(state))
		currentCycle++
	}
}
