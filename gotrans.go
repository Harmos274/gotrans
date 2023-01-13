// Responsible for parsing the input file and displaying the output
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Harmos274/gotrans/warehouse"
)

const helpText = "Gotrans\n" +
	"=========\n" +
	"Giving a file that describes a warehouse with packages, forklifts and trucks inside it,\n" +
	"the program will have to optimise the distribution of packages to trucks using the forklifts.\n\n" +
	"Commands:\n" +
	"-h --help\tShow help\n" +
	"-g --graphic\tActivate the graphic mode\n" +
	"<file>\t\tlaunch the program\n"

func main() {
	arguments := os.Args
	graphicMode := false

	if len(arguments) < 2 {
		_, _ = fmt.Fprint(os.Stderr, helpText)
		os.Exit(1)
	} else if arguments[1] == "-h" || arguments[1] == "--help" {
		fmt.Printf("%s\n", helpText)
		return
	} else if len(arguments) > 2 && (arguments[2] == "-g" || arguments[2] == "--graphic") {
		graphicMode = true
	}

	file, err := os.Open(arguments[1])
	if err != nil {
		fmt.Println("😱")
		_, _ = fmt.Fprint(os.Stderr, err)
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	initWr, cycles, err := parseInputFile(file)
	if err != nil {
		fmt.Println("😱")
		log.Fatal(err)
	}

	ch := make(chan warehouse.CycleState)

	go warehouse.CleanWarehouse(initWr, ch, cycles)

	currentCycle := 1
	for state := range ch {
		fmt.Printf("tour %d/%d\n", currentCycle, cycles)
		fmt.Println(showableWarehouse(state))
		currentCycle++
	}

	if currentCycle < int(cycles) {
		fmt.Println("😎")
	} else {
		if graphicMode {
			fmt.Print("!Warning: to active the graphic mode, the size of the map must not exceed 6x8\n\n")
		}
		// TUI
		ch := make(chan warehouse.CycleState)

		go warehouse.CleanWarehouse(initWr, ch, cycles)

		currentCycle := 1
		for state := range ch {
			fmt.Printf("tour %d/%d\n", currentCycle, cycles)
			fmt.Println(showableWarehouse(state))
			currentCycle++
		}
		if currentCycle < int(cycles) {
			fmt.Println("😎")
		} else {
			fmt.Println("🙂")
		}
	}

}
