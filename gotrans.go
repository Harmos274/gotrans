package main

import (
	"fmt"
	"log"
	"os"
)

type Warehouse struct {
	length, height int
	packages       EntityMap[Package]
	palletJacks    EntityMap[PalletJack]
	trucks         EntityMap[Truck]
}

func (wh Warehouse) somethingExistsAtThisPosition(pos Position) bool {
	return wh.packages.Exists(pos) || wh.palletJacks.Exists(pos) || wh.trucks.Exists(pos)
}

type EntityMap[T Package | PalletJack | Truck] map[Position]T

type Position struct {
	x, y int
}

type Package struct {
	weight Weight
	name   string
}
type Weight int

type PalletJack struct {
	name string
	pack *Package
}

type Truck struct {
	name                  string
	maxWeight             Weight
	elapseDischargingTime int
}

func (ettMap EntityMap[T]) Exists(pos Position) bool {
	_, exists := ettMap[pos]
	return exists
}

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
	return
}
