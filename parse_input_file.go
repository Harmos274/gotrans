package main

import (
	"bufio"
	"errors"
	"fmt"
	. "github.com/Harmos274/gotrans/warehouse"
	"os"
	"strconv"
	"strings"
)

func parseInputFile(file *os.File) (warehouse Warehouse, cycles uint, err error) {
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
		pack, pos, packErr := parsePackage(words)

		if packErr != nil {
			err = packErr
			return
		}
		if warehouse.SomethingExistsAtThisPosition(pos) {
			err = errors.New("two entities can't be at the same position")
			return
		}
		warehouse.Packages[pos] = pack
	}

	for {
		words := strings.Split(scanner.Text(), " ")

		if len(words) != 3 {
			break
		}
		pj, pos, pjErr := parsePalletJack(words)
		if pjErr != nil {
			err = pjErr
			return
		}
		if warehouse.SomethingExistsAtThisPosition(pos) {
			err = errors.New("two entities can't be at the same position")
			return
		}
		warehouse.PalletJacks[pos] = pj
		if !scanner.Scan() {
			return
		}
	}

	for {
		words := strings.Split(scanner.Text(), " ")

		if len(words) != 5 {
			err = errors.New("invalid formatting for truck and loading place")
			return
		}
		truck, pos, truckErr := parseTruck(words)
		if truckErr != nil {
			err = truckErr
			return
		}
		if warehouse.SomethingExistsAtThisPosition(pos) {
			err = errors.New("two entities can't be at the same position")
			return
		}
		warehouse.Trucks[pos] = truck
		if !scanner.Scan() {
			return
		}
	}
}

func parseWarehouse(line string) (warehouse Warehouse, cycles uint, err error) {
	const minimumCycles = 10
	const maximumCycles = 100_000
	_, err = fmt.Sscanf(line, "%d %d %d", &warehouse.Length, &warehouse.Height, &cycles)

	if err != nil {
		return
	}
	if cycles < minimumCycles || cycles > maximumCycles {
		err = errors.New("cycle should be between 10 and 100 000")
		return
	}
	warehouse.Packages = make(EntityMap[Package])
	warehouse.PalletJacks = make(EntityMap[PalletJack])
	warehouse.Trucks = make(EntityMap[Truck])
	return
}

func parsePackage(words []string) (pack Package, position Position, err error) {
	colorToWeight := map[string]Weight{
		"yellow": 100,
		"green":  200,
		"blue":   500,
	}
	pack.Name = words[0]
	x, err1 := strconv.Atoi(words[1])
	y, err2 := strconv.Atoi(words[2])
	weight, ok := colorToWeight[strings.ToLower(words[3])]

	if err1 != nil || err2 != nil || !ok {
		err = errors.New("invalid package formatting")
		return
	}
	position.X = x
	position.Y = y
	pack.Weight = weight
	return
}

func parsePalletJack(words []string) (pj PalletJack, position Position, err error) {
	pj.Name = words[0]
	x, err1 := strconv.Atoi(words[1])
	y, err2 := strconv.Atoi(words[2])

	if err1 != nil || err2 != nil {
		err = errors.New("invalid pallet jack formatting")
		return
	}
	position.X = x
	position.Y = y
	return
}

func parseTruck(words []string) (truck Truck, position Position, err error) {
	truck.Name = words[0]
	x, err1 := strconv.Atoi(words[1])
	y, err2 := strconv.Atoi(words[2])
	maxWeight, err3 := strconv.Atoi(words[3])
	elapseDischargingTime, err4 := strconv.Atoi(words[4])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		err = errors.New("invalid truck formatting")
		return
	}
	truck.ElapseDischargingTime = elapseDischargingTime
	truck.MaxWeight = Weight(maxWeight)
	position.X = x
	position.Y = y
	return
}
