package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("testMap.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	var warehouse Warehouse
	var cycles int

	if scanner.Scan() {
		warehouse, cycles, err = parseWarehouse(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(scanner.Err())
	}

	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")

		if len(words) != 4 {
			break
		}
		pack, pos, err := parsePackages(words)

		if err != nil {
			log.Fatal(err)
		}
		if warehouse.packages.Exists(pos) {
			log.Fatal("Two entities can't be at the same position.")
		}
		warehouse.packages[pos] = pack
	}

	fmt.Println(warehouse, cycles)
	return
}

func parseWarehouse(line string) (warehouse Warehouse, cycles int, err error) {
	const minimumCycles = 10
	const maximumCycles = 100_000
	_, err = fmt.Sscanf(line, "%d %d %d", &warehouse.length, &warehouse.height, &cycles)

	if err != nil {
		return
	}
	if cycles < minimumCycles || cycles > maximumCycles {
		err = errors.New("cycle should be between 10 and 100 000")
		return
	}
	warehouse.packages = make(EntityMap[Package])
	warehouse.palletJacks = make(EntityMap[PalletJack])
	warehouse.trucks = make(EntityMap[Truck])
	return
}

func parsePackages(words []string) (pack Package, position Position, err error) {
	colorToWeight := map[string]Weight{
		"yellow": 100,
		"green":  200,
		"blue":   500,
	}
	pack.name = words[0]
	x, err1 := strconv.Atoi(words[1])
	y, err2 := strconv.Atoi(words[2])
	weight, ok := colorToWeight[strings.ToLower(words[3])]

	if err1 != nil || err2 != nil || !ok {
		err = errors.New("invalid package formatting")
		return
	}
	position.x = x
	position.y = y
	pack.weight = weight
	return
}

type Position struct {
	x, y int
}
type Weight int

type PalletJack struct {
	name string
	pack *Package
}

type Package struct {
	weight Weight
	name   string
}

type Truck struct {
	maxWeight             Weight
	elapseDischargingTime int
}

type EntityMap[T Package | Truck | PalletJack] map[Position]T

func (ettMap EntityMap[T]) Exists(pos Position) bool {
	_, exists := ettMap[pos]
	return exists
}

type Warehouse struct {
	length, height int
	packages       EntityMap[Package]
	palletJacks    EntityMap[PalletJack]
	trucks         EntityMap[Truck]
}
