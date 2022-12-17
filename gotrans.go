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

	warehouse, cycles, err := parseInputFile(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(warehouse, cycles)
	return
}

func parseInputFile(file *os.File) (warehouse Warehouse, cycles int, err error) {
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		warehouse, cycles, err = parseWarehouse(scanner.Text())
		if err != nil {
			return
		}
	} else {
		err = errors.New("invalid file format")
		return
	}

	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")

		if len(words) != 4 {
			break
		}
		pack, pos, packErr := parsePackages(words)

		if packErr != nil {
			err = packErr
			return
		}
		if warehouse.packages.Exists(pos) || warehouse.palletJacks.Exists(pos) || warehouse.trucks.Exists(pos) {
			err = errors.New("two entities can't be at the same position")
			return
		}
		warehouse.packages[pos] = pack
	}

	for {
		words := strings.Split(scanner.Text(), " ")

		if len(words) != 3 {
			break
		}
		pj, pos, pjErr := parsePalletJacks(words)
		if pjErr != nil {
			err = pjErr
			return
		}
		if warehouse.packages.Exists(pos) || warehouse.palletJacks.Exists(pos) || warehouse.trucks.Exists(pos) {
			err = errors.New("two entities can't be at the same position")
			return
		}
		warehouse.palletJacks[pos] = pj
		if !scanner.Scan() {
			break
		}
	}

	for {
		words := strings.Split(scanner.Text(), " ")

		if len(words) != 5 {
			err = errors.New("invalid formatting for truck and loading place")
			return
		}
		truck, pos, truckErr := parseTrucks(words)
		if truckErr != nil {
			err = truckErr
			return
		}
		if warehouse.packages.Exists(pos) || warehouse.palletJacks.Exists(pos) || warehouse.trucks.Exists(pos) {
			err = errors.New("two entities can't be at the same position")
			return
		}
		warehouse.trucks[pos] = truck
		if !scanner.Scan() {
			break
		}
	}
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

func parsePalletJacks(words []string) (pj PalletJack, position Position, err error) {
	pj.name = words[0]
	x, err1 := strconv.Atoi(words[1])
	y, err2 := strconv.Atoi(words[2])

	if err1 != nil || err2 != nil {
		err = errors.New("invalid pallet jack formatting")
		return
	}
	position.x = x
	position.y = y
	return
}

func parseTrucks(words []string) (truck Truck, position Position, err error) {
	truck.name = words[0]
	x, err1 := strconv.Atoi(words[1])
	y, err2 := strconv.Atoi(words[2])
	maxWeight, err3 := strconv.Atoi(words[3])
	elapseDischargingTime, err4 := strconv.Atoi(words[4])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		err = errors.New("invalid truck formatting")
		return
	}
	truck.elapseDischargingTime = elapseDischargingTime
	truck.maxWeight = Weight(maxWeight)
	position.x = x
	position.y = y
	return
}

type Position struct {
	x, y int
}

type Warehouse struct {
	length, height int
	packages       EntityMap[Package]
	palletJacks    EntityMap[PalletJack]
	trucks         EntityMap[Truck]
}

type Weight int

type Package struct {
	weight Weight
	name   string
}

type PalletJack struct {
	name string
	pack *Package
}

type Truck struct {
	name                  string
	maxWeight             Weight
	elapseDischargingTime int
}

type EntityMap[T Package | Truck | PalletJack] map[Position]T

func (ettMap EntityMap[T]) Exists(pos Position) bool {
	_, exists := ettMap[pos]
	return exists
}
